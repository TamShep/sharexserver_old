package config

import (
	"github.com/mmichaelb/sharexserver/pkg/storage"
	"github.com/mmichaelb/sharexserver/pkg/storage/storages"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"time"
)

func loadMongoCfg(fileName string) (mongoCfg *viper.Viper, err error) {
	mongoCfg = viper.New()
	// set configuration filepath to the provided parameter
	mongoCfg.SetConfigFile(fileName)
	// add default values if the given config file does not contains specific values or do not exist
	// default values taken from ../../../configs/default-mongo-storage-config.toml
	mongoCfg.SetDefault("address", "localhost:27017")
	mongoCfg.SetDefault("connect_timeout", time.Second*4)
	mongoCfg.SetDefault("storage_folder", "./files/")
	mongoCfg.SetDefault("auth_db", "")
	mongoCfg.SetDefault("auth_user", "")
	mongoCfg.SetDefault("auth_passwd", "")
	mongoCfg.SetDefault("storage_db", "sharexserver")
	mongoCfg.SetDefault("storage_file_col", "uploads")
	// read config from filepath
	err = mongoCfg.ReadInConfig()
	return
}

func ParseMongoStorageFromConfig(fileName string) (storage storage.FileStorage, err error) {
	var mongoCfg *viper.Viper
	mongoCfg, err = loadMongoCfg(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Could not read Mongo storage configuration from file, %T: %v. Falling back to defaults.\n",
				err, err)
			err = nil
		}
	}
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoCfg.GetString("address")},
		Timeout:  mongoCfg.GetDuration("connect_timeout"),
		Source:   mongoCfg.GetString("auth_db"),
		Username: mongoCfg.GetString("auth_user"),
		Password: mongoCfg.GetString("auth_passwd"),
	}
	storage = &storages.MongoStorage{
		DialInfo:       dialInfo,
		DataFolder:     mongoCfg.GetString("storage_folder"),
		DatabaseName:   mongoCfg.GetString("storage_db"),
		CollectionName: mongoCfg.GetString("storage_file_col"),
	}
	return
}
