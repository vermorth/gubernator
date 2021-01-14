/*
Copyright 2018-2020 Mailgun Technologies Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gubernator

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"io"
	"net"
	"runtime"
	"strconv"

	ml "github.com/hashicorp/memberlist"
	"github.com/mailgun/holster/v3/clock"
	"github.com/mailgun/holster/v3/retry"
	"github.com/mailgun/holster/v3/setter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type MemberListPool struct {
	log        logrus.FieldLogger
	memberList *ml.Memberlist
	conf       MemberListPoolConfig
	events     *memberListEventHandler
}

type MemberListPoolConfig struct {
	// (Required) Callback address. Defaults to GUBER_ADVERTISE_ADDRESS
	CallbackAddress string

	// (Required) This is the address:port the member list protocol listen for other members on
	MemberListAddress string

	// (Required) This is the address:port the member list will advertise to other members it finds
	AdvertiseAddress string

	// (Required) A list of nodes this member list instance can contact to find other members.
	KnownNodes []string

	// (Required) A callback function which is called when the member list changes
	OnUpdate UpdateFunc

	// (Optional) The name of the node this member list identifies itself as.
	NodeName string

	// (Optional) An interface through which logging will occur (Usually *logrus.Entry)
	Logger logrus.FieldLogger

	// (Optional) The datacenter this instance belongs too
	DataCenter string
}

func NewMemberListPool(ctx context.Context, conf MemberListPoolConfig) (*MemberListPool, error) {
	setter.SetDefault(conf.Logger, logrus.WithField("category", "gubernator"))
	m := &MemberListPool{
		log:  conf.Logger,
		conf: conf,
	}

	lisHost, lisPort, err := splitAddress(conf.MemberListAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "MemberListAddress=`%s` is invalid;", conf.MemberListAddress)
	}

	// Member list requires the address to be an ip address
	if ip := net.ParseIP(lisHost); ip == nil {
		addrs, err := net.LookupHost(lisHost)
		if err != nil {
			return nil, errors.Wrapf(err, "while preforming host lookup for '%s'", lisHost)
		}
		if len(addrs) == 0 {
			return nil, errors.Wrapf(err, "net.LookupHost() returned no addresses for '%s'", lisHost)
		}
		lisHost = addrs[0]
	}

	advHost, advPort, err := splitAddress(conf.AdvertiseAddress)
	if err != nil {
		return nil, errors.Wrap(err, "AdvertiseAddress=`%s` is invalid;")
	}

	// Configure member list event handler
	m.events = newMemberListEventHandler(m.log, conf)

	// Configure member list
	config := ml.DefaultWANConfig()
	config.Name = conf.AdvertiseAddress
	config.Events = m.events
	config.BindAddr = lisHost
	config.BindPort = lisPort
	config.AdvertiseAddr = advHost
	config.AdvertisePort = advPort

	if conf.NodeName != "" {
		config.Name = conf.NodeName
	}

	config.LogOutput = newLogWriter(m.log)

	// Create and set member list
	memberList, err := ml.Create(config)
	if err != nil {
		return nil, err
	}
	m.memberList = memberList

	// Prep metadata
	_, callPort, err := splitAddress(conf.CallbackAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "CallbackAddress=`%s` is invalid;", conf.CallbackAddress)
	}
	gob.Register(memberListMetadata{})
	metadata := memberListMetadata{
		DataCenter:       conf.DataCenter,
		AdvertiseAddress: conf.CallbackAddress,
		GubernatorPort:   callPort,
	}

	// Join member list pool
	err = m.joinPool(ctx, conf.KnownNodes, metadata)
	if err != nil {
		return nil, errors.Wrap(err, "while attempting to join the member-list pool")
	}

	return m, nil
}

func (m *MemberListPool) joinPool(ctx context.Context, knownNodes []string, metadata memberListMetadata) error {
	// Get local node and set metadata
	node := m.memberList.LocalNode()
	serializedMetadata, err := serializeMemberListMetadata(metadata)
	if err != nil {
		return err
	}
	node.Meta = serializedMetadata

	err = retry.Until(ctx, retry.Interval(clock.Millisecond*300), func(ctx context.Context, i int) error {
		// Join member list
		m.log.Infof("MemberList Join KnownNodes: %v", knownNodes)
		_, err = m.memberList.Join(knownNodes)
		if err != nil {
			return errors.Wrap(err, "while joining member-list")
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "timed out attempting to join member list")
	}

	// Add the local node to the event handler's peer list
	m.events.addPeer(node)

	return nil
}

func (m *MemberListPool) Close() {
	err := m.memberList.Leave(clock.Second)
	if err != nil {
		m.log.Warn(errors.Wrap(err, "while leaving member-list"))
	}
}

type memberListEventHandler struct {
	peers map[string]PeerInfo
	log   logrus.FieldLogger
	conf  MemberListPoolConfig
}

func newMemberListEventHandler(log logrus.FieldLogger, conf MemberListPoolConfig) *memberListEventHandler {
	handler := memberListEventHandler{
		conf: conf,
		log:  log,
	}
	handler.peers = make(map[string]PeerInfo)
	return &handler
}

func (e *memberListEventHandler) addPeer(node *ml.Node) {
	ip := getIP(node.Address())

	// Deserialize metadata
	metadata, err := deserializeMemberListMetadata(node.Meta)
	if err != nil {
		e.log.WithError(err).Warnf("while adding to peers")
	} else {
		// Handle deprecated GubernatorPort
		if metadata.AdvertiseAddress == "" {
			metadata.AdvertiseAddress = makeAddress(ip, metadata.GubernatorPort)
		}
		e.peers[ip] = PeerInfo{GRPCAddress: metadata.AdvertiseAddress, DataCenter: metadata.DataCenter}
		e.callOnUpdate()
	}
}

func (e *memberListEventHandler) NotifyJoin(node *ml.Node) {
	ip := getIP(node.Address())

	// Deserialize metadata
	metadata, err := deserializeMemberListMetadata(node.Meta)
	if err != nil {
		// This is called during memberlist initialization due to the fact that the local node
		// has no metadata yet
		e.log.WithError(err).Warn("while deserialize member-list metadata")
	} else {
		// Handle deprecated GubernatorPort
		if metadata.AdvertiseAddress == "" {
			metadata.AdvertiseAddress = makeAddress(ip, metadata.GubernatorPort)
		}
		e.peers[ip] = PeerInfo{GRPCAddress: metadata.AdvertiseAddress, DataCenter: metadata.DataCenter}
		e.callOnUpdate()
	}
}

func (e *memberListEventHandler) NotifyLeave(node *ml.Node) {
	ip := getIP(node.Address())

	// Remove PeerInfo
	delete(e.peers, ip)

	e.callOnUpdate()
}

func (e *memberListEventHandler) NotifyUpdate(node *ml.Node) {
	ip := getIP(node.Address())

	// Deserialize metadata
	metadata, err := deserializeMemberListMetadata(node.Meta)
	if err != nil {
		e.log.WithError(err).Warn("while updating member-list")
	} else {
		// Construct Gubernator address and create PeerInfo
		gubernatorAddress := makeAddress(ip, metadata.GubernatorPort)
		e.peers[ip] = PeerInfo{GRPCAddress: gubernatorAddress, DataCenter: metadata.DataCenter}
		e.callOnUpdate()
	}
}

func (e *memberListEventHandler) callOnUpdate() {
	var peers []PeerInfo

	for _, p := range e.peers {
		e.log.Debugf("Peer GRPCAddress:%s, Peer HTTPAddress:%s, My CallbackAddress:%s, IsOwner:%v",
			p.GRPCAddress,
			p.HTTPAddress,
			e.conf.CallbackAddress,
			p.GRPCAddress == e.conf.CallbackAddress)
		if p.GRPCAddress == e.conf.CallbackAddress {
			p.IsOwner = true
		}
		peers = append(peers, p)
	}
	e.conf.OnUpdate(peers)
}

func getIP(address string) string {
	addr, _, _ := net.SplitHostPort(address)
	return addr
}

func makeAddress(ip string, port int) string {
	return net.JoinHostPort(ip, strconv.Itoa(port))
}

type memberListMetadata struct {
	DataCenter       string
	AdvertiseAddress string
	// Deprecated
	GubernatorPort int
}

func serializeMemberListMetadata(metadata memberListMetadata) ([]byte, error) {
	b, err := json.Marshal(&metadata)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling metadata as JSON")
	}
	return b, nil
}

func deserializeMemberListMetadata(b []byte) (*memberListMetadata, error) {
	var metadata memberListMetadata
	if err := json.Unmarshal(b, &metadata); err != nil {
		decoder := gob.NewDecoder(bytes.NewBuffer(b))
		if err := decoder.Decode(&metadata); err != nil {
			return nil, errors.Wrap(err, "error decoding metadata")
		}
	}
	return &metadata, nil
}

func newLogWriter(log logrus.FieldLogger) *io.PipeWriter {
	reader, writer := io.Pipe()

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Errorf("Error while reading from Writer: %s", err)
		}
		reader.Close()
	}()
	runtime.SetFinalizer(writer, func(w *io.PipeWriter) {
		writer.Close()
	})

	return writer
}

func splitAddress(addr string) (string, int, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return host, 0, errors.New(" expected format is `address:port`")
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return host, intPort, errors.Wrap(err, "port must be a number")
	}
	return host, intPort, nil
}
