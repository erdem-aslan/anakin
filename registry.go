package main

import (
	"container/ring"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func NewRegistry() *Registry {

	return &Registry{
		apps:        make(map[string]*Application),
		urlAppIndex: make(map[string]*Application),
		services:    make(map[string]*Service),
		endpoints:   make(map[string]*Endpoint),
	}
}

type Registry struct {
	apps        map[string]*Application
	urlAppIndex map[string]*Application
	al          sync.RWMutex

	services map[string]*Service
	sl       sync.RWMutex

	endpoints map[string]*Endpoint
	el        sync.RWMutex

	s Store
}

func (r *Registry) Init(store Store) error {

	log.Println("Initializing registry...")

	r.s = store

	appsSlice, err := store.GetApplications()

	if err != nil {
		return err
	}

	r.al.Lock()
	for _, app := range appsSlice {
		r.apps[app.UniqueId] = app
		r.urlAppIndex[app.BaseUrl] = app
	}
	r.al.Unlock()

	endSlice, err := store.GetEndpoints()

	if err != nil {
		return err
	}

	r.el.Lock()
	for _, end := range endSlice {
		log.Println("Caching endpoint: ", end)
		r.endpoints[end.UniqueId] = end
	}
	r.el.Unlock()

	serviceSlice, err := store.GetServices()

	if err != nil {
		return err
	}

	r.sl.Lock()
	for _, service := range serviceSlice {

		r.services[service.UniqueId] = service

		set := service.EndpointsSet()

		// build the ring of endpoints for round robin
		ri := ring.New(len(set))

		// build the slice of endpoints for source hashing
		sl := make([]*Endpoint, len(set))

		for eId, _ := range set {
			e := r.endpoints[eId]
			ri.Value = e
			ri = ri.Next()

			sl = append(sl, e)

		}

		service.serviceEPRing = ri
		service.serviceEPList = sl

	}
	r.sl.Unlock()

	log.Println("Initializing registry, finished")

	return nil

}

func (r *Registry) GetApplication(id string) *Application {
	r.al.RLock()
	defer r.al.RUnlock()
	return r.apps[id]
}

func (r *Registry) GetService(id string) *Service {
	r.sl.RLock()
	defer r.sl.RUnlock()
	return r.services[id]
}

func (r *Registry) GetEndpoint(id string) *Endpoint {
	r.el.RLock()
	defer r.el.RUnlock()
	return r.endpoints[id]
}

func (r *Registry) ExtractBaseUrl(target *url.URL) (baseUrl string, err error) {

	path := target.Path
	// Trim the leading `/`
	if len(path) > 1 && path[0] == '/' {
		path = path[1:]
	}
	// Explode on `/` and make sure we have at least
	// 1 element (application baseUrl)
	tmp := strings.Split(path, "/")

	if len(tmp) < 1 {
		return "", fmt.Errorf("Invalid path")
	}
	baseUrl = tmp[0]
	// Rewrite the request's path without the prefix.
	target.Path = "/" + strings.Join(tmp[1:], "/")
	return baseUrl, nil
}

func (r *Registry) ServiceForRequest(request *http.Request) *Service {

	target := request.URL
	baseUrl, err := r.ExtractBaseUrl(target)

	if err != nil {
		log.Println("Missing baseUrl within ", target)
		return nil
	}

	apps, err := store.GetApplications()

	if err != nil {
		log.Println("Error: ", err)
		return nil
	}

	var matchedApp *Application = nil

	for _, app := range apps {

		if app.BaseUrl == baseUrl && app.State == Active {
			matchedApp = app
		}
	}

	if matchedApp == nil {
		log.Println("No application has matched to base url: ", baseUrl)
		return nil
	}

	services := matchedApp.ServicesSorted()

	if len(services) == 0 {
		log.Println("No services has defined for base url: ", baseUrl)
	}

	var matchedService *Service = nil

	for _, serviceId := range services {

		service := r.GetService(serviceId)

		if service == nil {
			log.Println("Registry cache mismatch, reconstructing...")
			r.Init(r.s)
			return nil
		}

		if service.Nested {
			if strings.HasPrefix(target.Path, service.ServiceUrl) {
				matchedService = service
				break
			}
		} else {
			if service.ServiceUrl == target.Path {
				matchedService = service
				break
			}
		}
	}

	if matchedService == nil {
		log.Println("No service for target path: ", target.Path)
	}

	return matchedService
}

func (r *Registry) Endpoint(s *Service, req *http.Request) *Endpoint {

	var e *Endpoint = nil

	switch s.BalanceStrategy {

	case Round_Robin:

		s.slb.Lock()
		defer s.slb.Unlock()

		ri := s.serviceEPRing

		if ri == nil {
			return nil
		}

		for i := 0; i < ri.Len(); i++ {

			e = ri.Value.(*Endpoint)

			state := e.State()

			s.serviceEPRing = ri.Next()

			if state == Active || state == Trying {
				break
			}
		}

	case Source_Hashing:

		s.slb.RLock()
		defer s.slb.RUnlock()

		li := s.serviceEPList

		lil := len(li)

		if li == nil || lil == 0 {
			return nil
		}

		idx := hash(strings.Split(req.RemoteAddr, ":")[0]) % lil

		e = s.serviceEPList[idx]
		state := e.State()

		if state != Active || state != Trying {
			for _, value := range s.serviceEPList {

				if value.State() == Active || value.State() == Trying {
					e = value
				}

			}
		}
	}

	return e

}

// -- store callback interface impl.
//
func (r *Registry) ApplicationAdded(a *Application) {
	log.Println("Application added: ", a)

	r.al.Lock()
	defer r.al.Unlock()

	old := r.apps[a.UniqueId]

	if old != nil {
		delete(r.urlAppIndex, old.BaseUrl)
		delete(r.apps, old.UniqueId)
	}

	r.apps[a.UniqueId] = a
	r.urlAppIndex[a.BaseUrl] = a
}

func (r *Registry) ApplicationUpdated(a *Application) {
	log.Println("Application updated: ", a)

	r.al.Lock()
	defer r.al.Unlock()

	old := r.apps[a.UniqueId]

	if old != nil {
		delete(r.urlAppIndex, old.BaseUrl)
		delete(r.apps, old.UniqueId)
	}

	r.apps[a.UniqueId] = a
	r.urlAppIndex[a.BaseUrl] = a
}

func (r *Registry) ApplicationRemoved(id string) {
	log.Println("Application removed: ", id)

	r.al.Lock()
	defer r.al.Unlock()

	a := r.apps[id]

	if a != nil {
		delete(r.urlAppIndex, a.BaseUrl)
		delete(r.apps, id)
	}
}

func (r *Registry) ServiceAdded(s *Service) {
	log.Println("Service added: ", s)

	r.sl.Lock()
	defer r.sl.Unlock()
	r.services[s.UniqueId] = s
}

func (r *Registry) ServiceUpdated(s *Service) {
	log.Println("Service updated: ", s)
	r.sl.Lock()
	defer r.sl.Unlock()
	r.services[s.UniqueId] = s
}

func (r *Registry) ServiceRemoved(id string) {
	log.Println("Service removed: ", id)

	r.sl.Lock()
	defer r.sl.Unlock()

	delete(r.services, id)
}

func (r *Registry) EndpointAdded(e *Endpoint) {

	log.Println("Endpoint added: ", e)

	r.el.Lock()
	defer r.el.Unlock()
	r.endpoints[e.UniqueId] = e
}

func (r *Registry) EndpointUpdated(e *Endpoint) {

	log.Println("Endpoint updated: ", e)

	r.el.Lock()
	defer r.el.Unlock()
	r.endpoints[e.UniqueId] = e
}

func (r *Registry) EndpointRemoved(id string) {

	log.Println("Endpoint removed: ", id)

	r.al.Lock()
	defer r.al.Unlock()
	delete(r.apps, id)
}

func (r *Registry) RemoteRegistryEvent(message AnakinEvent) {

	payload := message.Payload

	switch message.EventType {
	case AppCreated:
		app, err := store.GetApplication(payload)

		if err != nil {
			log.Println("Failed processing remote appCreated event, error: ", err)
		}

		if app != nil || err != nil {
			r.ApplicationAdded(app)
		}
	case AppDeleted:
		r.ApplicationRemoved(payload)
	case AppUpdated:
		app, err := store.GetApplication(payload)

		if err != nil {
			log.Println("Failed processing remote appUpdated event, error: ", err)
		}

		if app != nil || err != nil {
			r.ApplicationUpdated(app)
		}

	case SvcCreated:
		svc, err := store.GetService(payload)

		if err != nil {
			log.Println("Failed processing remote svcCreated event, error: ", err)
		}

		if svc != nil || err != nil {
			r.ServiceAdded(svc)
		}

	case SvcDeleted:
		r.ServiceRemoved(payload)

	case SvcUpdated:
		svc, err := store.GetService(payload)

		if err != nil {
			log.Println("Failed processing remote svcUpdated event, error: ", err)
		}

		if svc != nil || err != nil {
			r.ServiceUpdated(svc)
		}

	case EndpCreated:
		endp, err := store.GetEndpoint(payload)

		if err != nil {
			log.Println("Failed processing remote endpCreated event, error: ", err)
		}

		if endp != nil || err != nil {
			r.EndpointAdded(endp)
		}

	case EndpDeleted:
		r.EndpointRemoved(payload)

	case EndpUpdated:
		endp, err := store.GetEndpoint(payload)

		if err != nil {
			log.Println("Failed processing remote appUpdated event, error: ", err)
		}

		if endp != nil || err != nil {
			r.EndpointUpdated(endp)
		}

	default:
		log.Println("Unhandled remote event, type: ", message.EventType, ", payload: ", payload)
	}
}
