package gubernator

import (
	"github.com/mailgun/gubernator/cache"
	"github.com/mailgun/gubernator/metrics"
)

type UpdateFunc func(*PeerConfig)

// Syncs configs and peer listings between peers
type PeerSyncer interface {
	RegisterOnUpdate(UpdateFunc)
	Start(string) error
	Stop()
}

// Config shared across peers via the PeerSyncer
type PeerConfig struct {
	Peers []string
}

// Local config for our service instance
type ServerConfig struct {
	// This is the interface and port number the service will listen to for GRPC connections
	// If unset, the server will pick a random port on localhost. You can retrieve the listening
	// address and port by calling `server.Address()`.
	ListenAddress string

	// This is the address and port number the service will advertise to other peers in the cluster
	// If unset, defaults to `ListenAddress`.
	AdvertiseAddress string

	// The cache implementation
	Cache cache.Cache

	// This is the implementation of peer syncer this server will use to keep all the peer
	// configurations in sync across the cluster.
	PeerSyncer PeerSyncer

	// This is the peer picker algorithm the server will use decide which peer in the cluster
	// will coordinate a rate limit
	Picker PeerPicker

	// Metrics collector
	Metrics metrics.Collector
}

// An implementation of PeerSyncer suitable for testing local clusters
type LocalPeerSyncer struct {
	callbacks []UpdateFunc
}

// Emits a new cluster config to all registered callbacks
func (sc *LocalPeerSyncer) Update(config PeerConfig) {
	for _, cb := range sc.callbacks {
		cb(&config)
	}
}

func (sc *LocalPeerSyncer) RegisterOnUpdate(cb UpdateFunc) {
	sc.callbacks = append(sc.callbacks, cb)
}

func (sc *LocalPeerSyncer) Start(addr string) error { return nil }
func (sc *LocalPeerSyncer) Stop()                   {}