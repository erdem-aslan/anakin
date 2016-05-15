package main

import (
	"github.com/hashicorp/logutils"
	"hash/fnv"
	"io"
	logger "log"
	"os"
	"os/signal"
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
	config        *Configuration
	anakinCluster *AnakinCluster
	store         Store
	registry      *Registry
	log           *logger.Logger
	filter        *logutils.LevelFilter
	stats         *StatsContainer
)

func init() {
	filter = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}

	log = logger.New(io.MultiWriter(os.Stdout), "<anakin>", logger.Ldate|logger.Ltime|logger.Lshortfile)

}

func main() {

	attachShutdownHook()

	// Config init does not return error, just exits the runtime.
	initConfig()

	err := initStore()

	if err != nil {
		log.Fatal("Store access has failed with error: ", err)
	}

	// init clustering
	err = initCluster()

	if err != nil {
		log.Fatal("Clustering has failed: ", err)
	}

	registry, err = initRegistry()

	if err != nil {
		log.Fatal("Registry initialization has failed with error: ", err)
	}

	store.AddListener(registry)

	initStats()

	go serveAdminBackend()

	err = serveProxy(registry)

	if err != nil {
		log.Fatal("Failed serving reverse proxy, error: ", err)
	}
}

func initStore() error {

	if config.MongoServers == nil {
		store = NewFsStore()
		return store.Initialize(config.DbPath + SEPARATOR + config.DbFileName)
	}

	store = NewMongoStore()
	return store.Initialize("")

}

func initRegistry() (r *Registry, err error) {
	r = NewRegistry()
	err = r.Init(store)

	return
}

func initStats() {
	stats = NewStatsContainer()
	stats.Start()
}

func attachShutdownHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		<-c
		if anakinCluster != nil {
			anakinCluster.Shutdown(true)
		}

		os.Exit(0)
	}()
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

// Helper wrapper type which implements sort.Interface
type SortByDESCServiceUrlLength []*Service

func (a SortByDESCServiceUrlLength) Len() int {
	return len(a)
}
func (a SortByDESCServiceUrlLength) Less(i, j int) bool {

	if a[i] != nil && a[j] != nil {
		return len(a[i].ServiceUrl) > len(a[j].ServiceUrl)
	}

	return true
}
func (a SortByDESCServiceUrlLength) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type SortInstanceById []*Instance

func (a SortInstanceById) Len() int {
	return len(a)
}
func (a SortInstanceById) Less(i, j int) bool {
	return hash(a[i].Id) < hash(a[j].Id)
}
func (a SortInstanceById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type SortAppById []*Application

func (a SortAppById) Len() int {
	return len(a)
}
func (a SortAppById) Less(i, j int) bool {
	return hash(a[i].Id()) < hash(a[j].Id())
}
func (a SortAppById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type SortServiceById []*Service

func (a SortServiceById) Len() int {
	return len(a)
}
func (a SortServiceById) Less(i, j int) bool {
	return hash(a[i].Id()) < hash(a[j].Id())
}
func (a SortServiceById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type SortEndpointById []*Endpoint

func (a SortEndpointById) Len() int {
	return len(a)
}
func (a SortEndpointById) Less(i, j int) bool {
	return hash(a[i].Id()) < hash(a[j].Id())
}
func (a SortEndpointById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
