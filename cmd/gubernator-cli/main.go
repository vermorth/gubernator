/*
Copyright 2018-2019 Mailgun Technologies Inc

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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	guber "github.com/mailgun/gubernator"
	"github.com/mailgun/holster/v3/clock"
	"github.com/mailgun/holster/v3/setter"
	"github.com/mailgun/holster/v3/syncutil"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func checkErr(err error) {
	if err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
}

func randInt(min, max int) int64 {
	return int64(rand.Intn(max-min) + min)
}

func main() {
	var apiType, GRPCAddress, HTTPAddress string
	var configFile string
	var totalReq int64
	var routine int
	var err error

	log = logrus.StandardLogger()
	flags := flag.NewFlagSet("gubernator", flag.ContinueOnError)
	flags.StringVar(&apiType, "api_type", "grpc", "grpc or http")
	flags.Int64Var(&totalReq, "req_count", 100000, "total count of req")
	flags.IntVar(&routine, "routine", 16, "concurrency")
	flags.StringVar(&configFile, "config", "", "environment config file")
	flags.StringVar(&GRPCAddress, "grpc_address", "", "the gubernator GRPC endpoint address")
	flags.StringVar(&HTTPAddress, "http_address", "", "the gubernator HTTP endpoint address")
	checkErr(flags.Parse(os.Args[1:]))

	if apiType != "grpc" && apiType != "http" {
		checkErr(errors.New("please provide a api type via -api_type"))
	}
	if routine <= 0 {
		routine = 16
	}
	if totalReq <= 0 {
		totalReq = 100000
	}

	conf, err := guber.SetupDaemonConfig(log, configFile)
	checkErr(err)
	setter.SetOverride(&conf.GRPCListenAddress, GRPCAddress)
	setter.SetOverride(&conf.HTTPListenAddress, HTTPAddress)

	if apiType == "grpc" && configFile == "" && GRPCAddress == "" && os.Getenv("GUBER_GRPC_ADDRESS") == "" {
		checkErr(errors.New("please provide a GRPC endpoint via -grpc_address or from a config " +
			"file via -config or set the env GUBER_GRPC_ADDRESS"))
	}

	if apiType == "http" && configFile == "" && HTTPAddress == "" && os.Getenv("GUBER_HTTP_ADDRESS") == "" {
		checkErr(errors.New("please provide a HTTP endpoint via -http_address or from a config " +
			"file via -config or set the env GUBER_HTTP_ADDRESS"))
	}

	// 统计相关初始化
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	qps := 0
	go func() {
		for {
			select {
			case t := <-ticker.C:
				log.Infof("time:%v qps:%v", t, qps)
				qps = 0
			}
		}
	}()
	minCost, _ := time.ParseDuration("1h")
	maxCost := time.Duration(0)
	avgCost := time.Duration(0)

	// 请求初始化
	// Generate a selection of rate limits with random limits
	var rateLimits []*guber.RateLimitReq

	for i := 0; i < 10000; i++ {
		rateLimits = append(rateLimits, &guber.RateLimitReq{
			Name:      "put_object_per_seconds",
			UniqueKey: fmt.Sprintf("appid=%v|bucket=%v", 1255000001, guber.RandomString(10)),
			Hits:      1,
			Limit:     randInt(1000, 50000),
			Duration:  1000,
			Algorithm: guber.Algorithm_TOKEN_BUCKET,
		})
	}

	// grpc bench
	if apiType == "grpc" {
		err = guber.SetupTLS(conf.TLS)
		checkErr(err)

		log.Infof("Connecting to '%s'...", conf.GRPCListenAddress)
		client, err := guber.DialV1Server(conf.GRPCListenAddress, conf.ClientTLS())
		checkErr(err)

		fan := syncutil.NewFanOut(routine)
		for i := int64(1); i <= totalReq; {
			for _, rateLimit := range rateLimits {
				fan.Run(func(obj interface{}) error {
					r := obj.(*guber.RateLimitReq)
					ctx, cancel := context.WithTimeout(context.Background(), clock.Millisecond*5000)
					// Now hit our cluster with the rate limits
					startT := time.Now()
					_, err := client.GetRateLimits(ctx, &guber.GetRateLimitsReq{
						Requests: []*guber.RateLimitReq{r},
					})
					tc := time.Since(startT)
					checkErr(err)
					cancel()

					if tc > maxCost {
						maxCost = tc
					}
					if tc < minCost {
						minCost = tc
					}
					avgCost = time.Duration(int64(avgCost)*(i-1)/i + int64(tc)/i)

					// if resp.Responses[0].Status == guber.Status_OVER_LIMIT {
					// 	spew.Dump(resp)
					// }
					qps++
					return nil
				}, rateLimit)
				i++
				if i > totalReq {
					break
				}
			}
		}
	}

	// http bench
	if apiType == "http" {
		tr := &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     30 * time.Second,
		}
		client := &http.Client{
			Transport: tr,
			Timeout:   5 * time.Second,
		}
		URL := "http://" + conf.HTTPListenAddress + "/v1/GetRateLimits"
		log.Infof("Posting to '%s'...", URL)

		fan := syncutil.NewFanOut(routine)
		for i := int64(1); i <= totalReq; {
			for _, rateLimit := range rateLimits {
				fan.Run(func(obj interface{}) error {
					r := obj.(*guber.RateLimitReq)
					body, _ := json.Marshal(r)
					startT := time.Now()
					// send req
					resp, err := client.Post(URL, "application/json", bytes.NewBuffer(body))
					tc := time.Since(startT)
					checkErr(err)
					_, err = ioutil.ReadAll(resp.Body)
					checkErr(err)
					//log.Info(string(buf))
					resp.Body.Close()
					if tc > maxCost {
						maxCost = tc
					}
					if tc < minCost {
						minCost = tc
					}
					avgCost = time.Duration(int64(avgCost)*(i-1)/i + int64(tc)/i)

					qps++
					return nil
				}, rateLimit)
				i++
				if i > totalReq {
					break
				}
			}
		}
	}

	log.Infof("apiType:%v, totalReq:%v, routineNum:%v, avgCost:%v, maxCost:%v, minCost:%v", apiType, totalReq, routine, avgCost, maxCost, minCost)
}
