package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

// this is a general configuration instance which is not bound to an instance of a struct
var Cfg *viper.Viper

func loadCfg(fileName string) (cfg *viper.Viper, err error) {
	cfg = viper.New()
	// set configuration filepath to the provided parameter
	cfg.SetConfigFile(fileName)
	// add default values if the given config file does not contains specific values or do not exist
	// default values taken from ../../../configs/default-config.toml
	cfg.SetDefault("webserver_address", "localhost:10711")
	cfg.SetDefault("storage_engine", "MongoDB+file")
	cfg.SetDefault("storage_engine_config", "./mongo-storage-config.toml")
	cfg.SetDefault("reverse_proxy_header", "")
	// read config from filepath
	err = cfg.ReadInConfig()
	return
}

func LoadMainConfig(fileName string) (err error) {
	if Cfg, err = loadCfg(fileName); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Could not read configuration from file, %T: %v. Falling back to defaults.\n", err, err)
			err = nil
		}
	}
	return
}
