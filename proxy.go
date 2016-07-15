package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

func serveProxy(reg *Registry) error {
	log.Println("Initializing service orchestrator...")

	proxy := NewMultipleHostReverseProxy(reg)

	log.Println("Initializing service orchestrator, finished, ", config.ProxyIp, config.ProxyPort)

	return http.ListenAndServe(config.ProxyIp+":"+strconv.Itoa(config.ProxyPort), proxy)
}

// NewMultipleHostReverseProxy creates a reverse proxy that will randomly
// select a host from the passed `targets`
func NewMultipleHostReverseProxy(reg *Registry) *httputil.ReverseProxy {

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, piggyBack string) (net.Conn, error) {

			if piggyBack == "" {
				return nil, errors.New("Failed to match any service")
			}

			// get rid of :80 appended by stack
			piggyBack = strings.Split(piggyBack, ":")[0]

			ss := strings.Split(piggyBack, "|")

			serviceId := ss[0]
			senderIp := ss[1]

			service := reg.GetService(serviceId)

			if service == nil {
				return nil, errors.New(fmt.Sprintf("Service is missing and/or deleted in the middle of routing, id: %s", serviceId))
			}
			endpoint := reg.NextAvailableEndpoint(service, senderIp)

			if endpoint == nil {
				return nil, errors.New(fmt.Sprintf("No endpoint is available/present for %s", service))
			}

			var conn net.Conn
			var err error

			for {
				conn, err = net.Dial(network, endpoint.Address())

				stats.IncrementEndpoint(endpoint.UniqueId)

				if err != nil {
					endpoint.SetState(Failing)
					store.UpdateEndpoint(endpoint)
					log.Println("Endpoint failing, address:", endpoint.Address())

					endpoint = reg.NextAvailableEndpoint(service, senderIp)

					if endpoint == nil {
						log.Println("All endpoints are failing for service:", service)
						return nil, errors.New("Endpoint failing, address: " + endpoint.Address())
					}

				}

				break

			}

			if endpoint.State() == Trying {
				endpoint.SetState(Active)
				store.UpdateEndpoint(endpoint)
			}

			return conn, err
		},

		TLSHandshakeTimeout: 10 * time.Second,
	}

	director := func(req *http.Request) {
		service := reg.ServiceForRequest(req)

		if service == nil {
			log.Printf("No service has matched for %s\n", req.URL.Path)
			return
		}

		senderIp := strings.Split(req.RemoteAddr, ":")[0]

		req.URL.Host = service.UniqueId + "|" + senderIp
	}

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
}
