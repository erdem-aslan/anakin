package main

import (
	"bufio"
	"errors"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/satori/go.uuid"
	"os"
	"strings"
)

const SEPARATOR string = string(os.PathSeparator)

func initConfig() {

	var cfgPath string

	flag.StringVar(&cfgPath,
		"cfg",
		DefaultCfgPath,
		"Path of the configuration file, ex: /path/to/anakin/cfg/anakin.toml")
	flag.Parse()

	log.Println("Initializing configuration...")

	_, err := os.Stat(cfgPath)

	if err != nil {

		currentPath, err := os.Getwd()

		if err != nil {
			flag.PrintDefaults()
			log.Fatal("Attempting to use default values has failed, error:", err)
		}

		console := bufio.NewReader(os.Stdin)

		var response string = ""

		log.Print("Continue with default values? (N/y) : ")

		for {

			response, err = console.ReadString('\n')
			response = strings.TrimSuffix(response, "\n")

			if response != "" &&
				response != "N" && response != "n" && response != "Y" && response != "y" {
				log.Println("Only Y/y/N/n or empty for default value(s) are valid.")
				continue
			}

			break

		}

		if response == "" || response == "n" || response == "N" {
			log.Println("Exiting...")
			os.Exit(0)
		}

		log.Println("Generating anakin.toml with necessary directory structure...")

		cfgPath := currentPath + SEPARATOR + DefaultCfgPath
		dbPath := currentPath + SEPARATOR + DefaultDbPath

		_, err = os.Stat(cfgPath)

		if err != nil {
			err = os.Mkdir(cfgPath, os.ModePerm)

			if err != nil {
				log.Fatal("Failed creating directory", cfgPath, err)
			}
		}

		_, err = os.Stat(dbPath)

		if err != nil {
			err = os.Mkdir(dbPath, os.ModePerm)

			if err != nil {
				log.Fatal("Failed creating directory", dbPath, err)
			}
		}

		os.Mkdir(dbPath, os.ModePerm)

		cfgPath = cfgPath + SEPARATOR + "anakin.toml"

		tomlFile, err := os.Create(cfgPath)

		if err != nil {
			log.Fatal("Configuration file creation has failed, error:", err)
		}

		defer tomlFile.Close()

		config = &Configuration{
			AdminIp:       DefaultAdminIp,
			AdminPort:     DefaultAdminPort,
			AdminToken:    uuid.NewV4().String(),
			DbPath:        dbPath,
			DbFileName:    DefaultDbFileName,
			ProxyIp:       DefaultProxyIp,
			ProxyPort:     DefaultProxyPort,
			ProxyRootPath: DefaultProxyRootPath,
			LogDir:        DefaultLogDir,
		}

		toml.NewEncoder(tomlFile).Encode(config)

	}

	if config == nil {

		var c Configuration

		_, err = toml.DecodeFile(cfgPath+SEPARATOR+DefaultCfgFile, &c)

		if err != nil {
			log.Fatal("Failed parsing current configuration, error: ", err)
		}

		config = &c
	}

	err = validateConfig(config)

	if err != nil {
		log.Fatal("Configuration error: ", err)
	}

	log.Println("Initializing configuration, finished")

}

func validateConfig(c *Configuration) error {

	if c == nil {
		return errors.New("Invalid configuration file")
	}

	return c.IsMandatoryFieldsValid()
}

type Configuration struct {
	AdminIp    string
	AdminPort  int
	AdminToken string

	ProxyIp       string
	ProxyPort     int
	ProxyRootPath string

	ClusterIp      string
	ClusterPort    int
	ClusterMembers []string

	MongoServers []string

	DbPath     string
	DbFileName string
	LogDir     string
	CertFile   string
	PemFile    string
}

func (c *Configuration) IsMandatoryFieldsValid() error {

	if c.AdminPort != 0 &&
		c.AdminToken != "" &&
		c.ProxyPort != 0 &&
		c.ProxyRootPath != "" &&
		c.DbPath != "" {
		return nil
	} else {
		return errors.New("Missing mandatory fields")
	}
}
