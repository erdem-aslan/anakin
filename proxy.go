package main

import (
	"errors"
	"log"
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

	log.Println("Initializing service orchestrator, finished")

	return http.ListenAndServe(config.ProxyIp+":"+strconv.Itoa(config.ProxyPort), proxy)
}

// NewMultipleHostReverseProxy creates a reverse proxy that will randomly
// select a host from the passed `targets`
func NewMultipleHostReverseProxy(reg *Registry) *httputil.ReverseProxy {

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, endpointId string) (net.Conn, error) {

			// get rid of :80 appended by stack
			endpointId = strings.Split(endpointId, ":")[0]

			endpoint := reg.GetEndpoint(endpointId)

			conn, err := net.Dial(network, endpoint.Address())

			if err != nil {
				endpoint.SetState(Failing)
				store.UpdateEndpoint(endpoint)
				return nil, errors.New("Endpoint failing, address: " + endpoint.Address())
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
		req.URL.Scheme = "http"

		if service == nil {
			log.Printf("No service has matched for %s\n", req.URL.Path)
			return
		}

		endpoint := reg.Endpoint(service, req)

		if endpoint == nil {
			log.Printf("No endpoint is available/present for %s\n", req.URL.Path)
			return

		}

		req.URL.Host = endpoint.UniqueId
	}

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
}
