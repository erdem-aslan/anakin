package main

import (
	"hash/fnv"
	"log"
)

type State string

const (
	Active    State = "active"
	Passive   State = "passive"
	Suspended State = "suspended"
	Failing   State = "failing"
	Trying    State = "trying"

	Version           string = "1.0"
	DefaultCfgPath    string = "cfg"
	DefaultCfgFile    string = "anakin.toml"
	DefaultDbPath     string = "db"
	DefaultDbFileName string = "anakin.db"

	DefaultAdminIp   string = ""
	DefaultAdminPort int    = 16016

	DefaultProxyIp       string = ""
	DefaultProxyPort     int    = 16015
	DefaultProxyRootPath string = "/"
)

var (
	config *Configuration
	store  Store
)

func init() {
	log.SetPrefix("<anakin " + Version + "> ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {

	// Config init does not return error, just exits the runtime.
	initConfig()

	err := initStore()

	if err != nil {
		log.Fatal("Store access has failed with error: ", err)
	}

	r, err := initRegistry()

	if err != nil {
		log.Fatal("Registry initialization has failed with error: ", err)
	}

	store.AddListener(r)

	go serveAdminBackend()

	err = serveProxy(r)

	if err != nil {
		log.Fatal("Failed serving reverse proxy, error: ", err)
	}
}

func initStore() error {

	store = NewFsStore()
	return store.Initialize(config.DbPath + SEPARATOR + config.DbFileName)
}

func initRegistry() (r *Registry, err error) {
	r = NewRegistry()
	err = r.Init(store)

	return
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}
