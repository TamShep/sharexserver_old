package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

// Cfg is a general configuration instance which is not bound to an instance of a struct
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
	cfg.SetDefault("whitelisted_content_types", []string{
		"image/png", "image/jpeg", "image/jpg", "image/gif",
		"text/plain", "text/plain; charset=utf-8",
		"video/mp4", "video/mpeg", "video/mpg4", "video/mpeg4", "video/flv",
	})
	// read config from filepath
	err = cfg.ReadInConfig()
	return
}

// LoadMainConfig loads the main config and stores the data into the Cfg variable.
func LoadMainConfig(fileName string) (err error) {
	if Cfg, err = loadCfg(fileName); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Could not read configuration from file, %T: %v. Falling back to defaults.\n", err, err)
			err = nil
		}
	}
	return
}
