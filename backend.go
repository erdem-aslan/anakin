package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func serveAdminBackend() {

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/anakin/v1",
		handlePing).
		Methods("GET")

	r.HandleFunc("/anakin/v1/cluster",
		handleCluster).
		Methods("GET")

	r.HandleFunc("/anakin/v1/local",
		handleLocal).
		Methods("GET")

	r.HandleFunc("/anakin/v1/apps",
		handleApplications).
		Methods("GET", "POST")

	r.HandleFunc("/anakin/v1/apps/{appId}",
		handleApplication).
		Methods("GET", "PUT", "DELETE")

	r.HandleFunc("/anakin/v1/apps/{appId}/services",
		handleServices).
		Methods("GET", "POST")

	r.HandleFunc("/anakin/v1/apps/{appId}/services/{serviceId}",
		handleService).
		Methods("GET", "PUT", "DELETE")

	r.HandleFunc("/anakin/v1/apps/{appId}/services/{serviceId}/endpoints",
		handleEndpoints).
		Methods("GET", "POST")

	r.HandleFunc("/anakin/v1/apps/{appId}/services/{serviceId}/endpoints/{endpointId}",
		handleEndpoint).
		Methods("GET", "PUT", "DELETE")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	var address string = ""

	if config.AdminIp != DefaultAdminIp {
		address = DefaultAdminIp
	}

	address = address + ":" + strconv.Itoa(config.AdminPort)

	log.Println("Serving administration backend on ", address)

	http.ListenAndServe(address, r)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handleCluster(w http.ResponseWriter, r *http.Request) {

	others, local := anakinCluster.Instances()

	instances := make([]*Instance, 0)
	instances = append(instances, local)

	if others != nil || len(others) != 0 {

		for _, instance := range others {

			req, _ := http.NewRequest("GET", "http://"+instance.AdminIp+
				":"+instance.AdminPort+"/anakin/v1/local", nil)

			req.Header.Set("Accept", "application/json")

			resp, err := http.DefaultClient.Do(req)

			if err != nil {
				log.Println("Instance cannot be reached: ", instance, err)
				instance.State = Failing
				continue
			}

			var remote *Instance = new(Instance)

			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(remote)

			if err != nil {
				internalError(w, err)
				return
			}

			instances = append(instances, remote)

		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	sort.Sort(SortInstanceById(instances))

	err := json.NewEncoder(w).Encode(instances)

	if err != nil {
		internalError(w, err)
		return
	}
}

func handleLocal(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(anakinCluster.LocalInstance())

	if err != nil {
		internalError(w, err)
		return
	}

}

func handleApplications(w http.ResponseWriter, r *http.Request) {

	applications, err := store.GetApplications()

	if err != nil {
		internalError(w, err)
		return
	}

	switch r.Method {

	case "GET":

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		sort.Sort(SortAppById(applications))

		err := json.NewEncoder(w).Encode(applications)

		if err != nil {
			internalError(w, err)
			return
		}

	case "POST":

		contentType := r.Header.Get("Content-Type")

		if contentType == "" || contentType != "application/json" {
			badRequest(w, errors.New("Invalid/missing content-type header"))
			return
		}

		var app *Application = new(Application)

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(app)

		if err != nil {
			log.Println("Error: ", err)
			badRequest(w, err)
			return
		}

		err = app.Init()

		if err != nil {
			badRequest(w, err)
			return
		}

		apps, err := store.GetApplications()

		if err != nil {
			internalError(w, err)
			return
		}

		for _, ap := range apps {
			if ap.Name == app.Name {
				badRequest(w, errors.New("Name already present: "+ap.Name))
				return
			}

			if ap.BaseUrl == app.BaseUrl {
				badRequest(w, errors.New("Base url already present: "+ap.BaseUrl))
				return
			}
		}

		err = store.CreateApplication(app)

		if err != nil {

			if err == AlreadyPresentError {
				badRequest(w, err)
			} else {
				internalError(w, err)
			}

			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(app)

		if err != nil {
			internalError(w, err)
		}

	}
}

func handleApplication(w http.ResponseWriter, r *http.Request) {

	appId := mux.Vars(r)["appId"]
	app, err := store.GetApplication(appId)

	switch r.Method {

	case "GET":

		if err != nil {
			internalError(w, err)
			return
		}

		if app == nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Add("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(app)

		if err != nil {
			internalError(w, err)
			return
		}

	case "PUT":

		contentType := r.Header.Get("Content-Type")

		if contentType == "" || contentType != "application/json" {
			badRequest(w, errors.New("Invalid/missing content-type header"))
			return
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(app)

		if err != nil {
			badRequest(w, err)
			return
		}

		apps, err := store.GetApplications()

		if err != nil {
			internalError(w, err)
			return
		}

		for _, a := range apps {

			if a.Id() == app.Id() {
				continue
			}

			if a.BaseUrl == app.BaseUrl {
				badRequest(w, errors.New("Base url already defined: "+app.BaseUrl))
				return
			}

			if a.Name == app.Name {
				badRequest(w, errors.New("Name already defined: "+app.BaseUrl))
				return
			}
		}

		err = store.UpdateApplication(app)

		if err != nil {

			if err == MissingEntryError {
				http.NotFound(w, r)
				return
			}

			internalError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)

	case "DELETE":

		err := store.DeleteApplication(appId)
		if err != nil {
			internalError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleServices(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	appId := vars["appId"]

	app, err := store.GetApplication(appId)

	if err != nil {
		internalError(w, err)
		return
	}

	if app == nil {
		http.NotFound(w, r)
		return
	}

	switch r.Method {

	case "GET":

		w.Header().Add("Content-Type", "application/json")

		services := make([]*Service, 0, 10)

		for serviceId, _ := range app.ServicesCopy() {
			s, err := store.GetService(serviceId)

			if err != nil {
				log.Println("Service fetch error:", err)
				continue
			}

			services = append(services, s)
		}

		sort.Sort(SortServiceById(services))

		err := json.NewEncoder(w).Encode(services)

		if err != nil {
			internalError(w, err)
			return
		}

	case "POST":

		contentType := r.Header.Get("Content-Type")

		if contentType == "" || contentType != "application/json" {
			badRequest(w, errors.New("Invalid/missing content-type header"))
			return
		}

		var s *Service = &Service{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(s)

		if err != nil {
			badRequest(w, err)
			return
		}

		err = s.Init()

		if err != nil {
			badRequest(w, err)
			return
		}

		if !strings.HasPrefix(s.ServiceUrl, "/") {
			s.ServiceUrl = "/" + s.ServiceUrl
		}

		for serviceId, _ := range app.ServicesCopy() {

			svc, err := store.GetService(serviceId)

			if err != nil {
				internalError(w, err)
				return
			}

			if svc.ServiceUrl == s.ServiceUrl {
				badRequest(w, errors.New("Service url is already present: "+s.ServiceUrl))
				return
			}
		}

		if s.ServiceUrl == "" {
			badRequest(w, errors.New("ServiceUrl is mandatory"))
			return
		}

		err = store.CreateService(s)

		if err != nil {
			internalError(w, err)
			return
		}

		app.AddService(s)

		err = store.UpdateApplication(app)

		if err != nil {

			if err == MissingEntryError {
				http.NotFound(w, r)
			} else {
				internalError(w, err)
			}

			store.DeleteService(s.Id())

			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(s)

		if err != nil {
			internalError(w, err)
		}
	}

}

func handleService(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	appId := vars["appId"]
	serviceId := vars["serviceId"]

	app, err := store.GetApplication(appId)

	if err != nil {
		internalError(w, err)
		return
	}

	if app == nil {
		http.NotFound(w, r)
		return
	}

	if !app.ContainsService(serviceId) {
		http.NotFound(w, r)
		return
	}

	s, err := store.GetService(serviceId)

	if err != nil {
		internalError(w, err)
		return
	}

	switch r.Method {

	case "GET":

		if s == nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(s)

		if err != nil {
			internalError(w, err)
			return
		}

	case "PUT":

		if s == nil {
			http.NotFound(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")

		if contentType == "" || contentType != "application/json" {
			badRequest(w, errors.New("Invalid/missing content-type header"))
			return
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(s)

		if err != nil {
			badRequest(w, err)
			return
		}

		if !strings.HasPrefix(s.ServiceUrl, "/") {
			s.ServiceUrl = "/" + s.ServiceUrl
		}

		for id, serviceUrl := range app.ServicesCopy() {

			if id == s.Id() {
				continue
			}

			if s.ServiceUrl == serviceUrl {
				badRequest(w, errors.New("Service url belongs to other service"))
				return
			}
		}

		log.Println("Update attempt on: ", s)

		err = store.UpdateService(s)

		if err != nil {

			if err == MissingEntryError {
				http.NotFound(w, r)
			} else {

				internalError(w, err)
			}

			return
		}

		w.WriteHeader(http.StatusOK)

	case "DELETE":
		app.RemoveServiceId(serviceId)
		store.UpdateApplication(app)
		store.DeleteService(serviceId)
		w.WriteHeader(http.StatusOK)
	}

}

func handleEndpoints(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	appId := vars["appId"]
	serviceId := vars["serviceId"]

	app, err := store.GetApplication(appId)

	if err != nil {
		internalError(w, err)
		return
	}

	if app == nil {
		http.NotFound(w, r)
		return
	}

	services := app.ServicesCopy()

	if len(services) == 0 || services[serviceId] == "" {
		http.NotFound(w, r)
		return
	}

	s, err := store.GetService(serviceId)

	if err != nil {
		internalError(w, err)
		return
	}

	switch r.Method {

	case "GET":

		w.Header().Add("Content-Type", "application/json")

		endpoints := make([]*Endpoint, 0, 10)

		for id, _ := range s.EndpointsSet() {
			endp, err := store.GetEndpoint(id)

			if err != nil {
				log.Println(err)
				continue
			}
			endpoints = append(endpoints, endp)

		}

		sort.Sort(SortEndpointById(endpoints))

		err := json.NewEncoder(w).Encode(endpoints)

		if err != nil {
			internalError(w, err)
			return
		}

	case "POST":

		contentType := r.Header.Get("Content-Type")

		if contentType == "" || contentType != "application/json" {
			badRequest(w, errors.New("Invalid/missing content-type header"))
			return
		}

		var e *Endpoint = &Endpoint{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(e)

		if err != nil {
			badRequest(w, err)
			return
		}

		if e.Host == "" || e.Port == "" {
			badRequest(w, errors.New("invalid host/port"))
			return
		}

		pValue, err := strconv.Atoi(e.Port)

		if err != nil {
			badRequest(w, errors.New("Port value should be numeric"))
			return
		}

		if pValue > 65535 {
			badRequest(w, errors.New("Port cannot be greater than 65535"))
			return
		}

		for id, _ := range s.EndpointsSet() {

			v, err := store.GetEndpoint(id)

			if err != nil {
				internalError(w, err)
				return
			}

			if e.Address() == v.Address() {
				badRequest(w, AlreadyPresentError)
				return
			}
		}

		e.Init()

		err = store.CreateEndpoint(e)

		if err != nil {
			internalError(w, err)
			return
		}

		s.AddEndpoint(e.Id())

		err = store.UpdateService(s)

		if err != nil {

			if err == MissingEntryError {
				http.NotFound(w, r)
			} else {
				internalError(w, err)
			}

			store.DeleteEndpoint(e.Id())

		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(e)

		if err != nil {
			internalError(w, err)
		}

	}

}

func handleEndpoint(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	appId := vars["appId"]
	serviceId := vars["serviceId"]
	endpointId := vars["endpointId"]

	app, err := store.GetApplication(appId)

	if err != nil {
		internalError(w, err)
		return
	}

	if app == nil {
		http.NotFound(w, r)
		return
	}

	services := app.ServicesCopy()

	if len(services) == 0 || services[serviceId] == "" {
		http.NotFound(w, r)
		return
	}

	s, err := store.GetService(serviceId)

	if err != nil {
		internalError(w, err)
		return
	}

	if s == nil {
		http.NotFound(w, r)
		return
	}

	e, err := store.GetEndpoint(endpointId)

	if err != nil {
		internalError(w, err)
		return
	}

	switch r.Method {

	case "GET":

		if e == nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(e)

		if err != nil {
			internalError(w, err)
			return
		}

	case "PUT":

		if e == nil {
			http.NotFound(w, r)
			return
		}

		contentType := r.Header.Get("Content-Type")

		if contentType == "" || contentType != "application/json" {
			badRequest(w, errors.New("Invalid/missing content-type header"))
			return
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(e)

		if err != nil {
			badRequest(w, err)
			return
		}

		if e.Host == "" || e.Port == "" {
			badRequest(w, errors.New("invalid host/port"))
			return
		}

		pValue, err := strconv.Atoi(e.Port)

		if err != nil {
			badRequest(w, errors.New("Port value should be numeric"))
			return
		}

		if pValue > 65535 {
			badRequest(w, errors.New("Port cannot be greater than 65535"))
			return
		}

		err = store.UpdateEndpoint(e)

		if err != nil {

			if err == MissingEntryError {
				http.NotFound(w, r)
				return
			}

			internalError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)

	case "DELETE":
		s.RemoveEndpoint(endpointId)
		store.UpdateService(s)
		store.DeleteEndpoint(endpointId)
		w.WriteHeader(http.StatusOK)
	}
}

func internalError(w http.ResponseWriter, err error) {
	log.Println("Internal error: ", err)
	http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, err error) {
	log.Println("Bad request error: ", err)
	http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
}
