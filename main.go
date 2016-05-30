package main

import (
	"bufio"
	"github.com/hashicorp/logutils"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"hash/fnv"
	"io"
	logger "log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"time"
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
	DefaultLogDir        string = "log"
)

var (
	config        *Configuration
	anakinCluster *AnakinCluster
	store         Store
	registry      *Registry
	log           *logger.Logger
	filter        *logutils.LevelFilter
	fileLogger    *RotatingFileWriter
	stats         *StatsContainer
)

func init() {

	filter = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}

	log = logger.New(os.Stdout, "<anakin>", logger.Ldate|logger.Ltime|logger.Lshortfile)

}

func main() {

	attachShutdownHook()

	// Config init does not return error, just exits the runtime.
	initConfig()

	fileLogger = NewRotatingFileWriter(config.LogDir, "anakin", config.LogDir+"/backups", 1024*1024*1024*10, 10)
	err := fileLogger.Init()

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(fileLogger)

	err = initStore()

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

	go serveAdminBackend("")

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

		if fileLogger != nil {
			fileLogger.Shutdown()
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

func NewRotatingFileWriter(filePath string,
	filePrefix string,
	backupDir string,
	maxSizeInBytes int64,
	maxBackup int) *RotatingFileWriter {

	return &RotatingFileWriter{
		maxBackup:  maxBackup,
		interval:   time.Millisecond * 100,
		filePath:   filePath,
		filePrefix: filePrefix,
		backupDir:  backupDir,
		maxSize:    maxSizeInBytes,
		logCh:      make(chan string, 10000),
		shCh:       make(chan struct{}, 1),
	}
}

type RotatingFileWriter struct {
	interval   time.Duration
	maxSize    int64
	maxBackup  int
	filePrefix string
	filePath   string
	file       *os.File

	backupDir string

	ticker *time.Ticker
	logCh  chan string
	shCh   chan struct{}

	isInit bool
	initL  sync.RWMutex // only used at init/start
	fileL  sync.RWMutex // general usage
}

// io.Writer interface implementation
func (rlw *RotatingFileWriter) Write(b []byte) (n int, err error) {
	rlw.logCh <- string(b)
	return len(b), err
}

func (rlw *RotatingFileWriter) Init() error {

	rlw.initL.RLock()
	if rlw.isInit {
		rlw.initL.RUnlock()
		log.Println("RotatingLogWriter is already in init state, ignoring init request.")
		return nil
	}
	rlw.initL.RUnlock()

	// ----- Init
	rlw.initL.Lock()
	defer rlw.initL.Unlock()

	var err error

	if !filepath.IsAbs(rlw.filePath) {

		rlw.filePath, err = filepath.Abs(rlw.filePath)

		if err != nil {
			return err
		}
	}

	log.Println("Creating directory", rlw.filePath, " if missing...")

	err = os.MkdirAll(rlw.filePath, 0755)

	if err != nil {
		return err
	}

	f, err := os.OpenFile(rlw.filePath+string(os.PathSeparator)+rlw.filePrefix+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	rlw.file = f

	// create backupDirectory if missing

	if !filepath.IsAbs(rlw.backupDir) {
		rlw.backupDir = rlw.filePath + string(os.PathSeparator) + rlw.backupDir
	}

	err = os.MkdirAll(rlw.backupDir, 0755)

	if err != nil {
		return err
	}

	rlw.rotateIfNecessary()

	rlw.ticker = time.NewTicker(rlw.interval)

	go func() {

	forever:
		for {
			select {
			case <-rlw.ticker.C:
				rlw.rotateIfNecessary()
				if len(rlw.logCh) != 0 {
					rlw.persistBuffer()
				}

			case <-rlw.shCh:
				log.Println("Flushing and shutting down...")

				rlw.persistBuffer()
				rlw.file.Sync()
				rlw.file.Close()
				break forever
			}
		}
	}()

	return nil
}

func (rlw *RotatingFileWriter) Shutdown() {
	rlw.shCh <- struct{}{}
}

func (rlw *RotatingFileWriter) rotateIfNecessary() error {
	rlw.fileL.RLock()

	info, err := rlw.file.Stat()

	if err != nil {
		return err
	}

	size := info.Size()

	rlw.fileL.RUnlock()

	if size < rlw.maxSize {
		return nil
	}

	rlw.fileL.Lock()
	defer rlw.fileL.Unlock()

	rlw.file.Sync()

	backup, err := os.OpenFile(rlw.backupDir+
		strconv.Itoa(os.PathSeparator)+
		rlw.filePrefix+"-"+time.Now().Format(time.RFC3339)+
		".log",
		os.O_RDWR,
		0660)

	if err != nil {
		return err
	}

	defer backup.Close()

	r := bufio.NewReader(rlw.file)
	w := bufio.NewWriter(backup)

	buf := make([]byte, 4096)

	for {
		read, err := r.Read(buf)

		if err != nil && err != io.EOF {
			return err
		}

		if read == 0 {
			break
		}

		write, err := w.Write(buf)

		if err != nil {
			return err
		}

		if write != read {
			return errors.New("Failed rotating file, read" + strconv.Itoa(read) + " bytes but written only" + strconv.Itoa(write) + "bytes")
		}
	}

	return nil

}

func (rlw *RotatingFileWriter) persistBuffer() {

persistLoop:
	for {
		select {
		case s := <-rlw.logCh:
			rlw.file.WriteString(s)
		default:
			break persistLoop
		}
	}
}
