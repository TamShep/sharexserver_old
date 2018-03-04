package config

import (
	"strconv"
	"testing"
	"time"
)

func TestMongoConfig(t *testing.T) {
	cfg, err := loadMongoCfg("../../../test/test-mongo-storage-config.toml")
	if err != nil {
		t.Fatalf("Could not load Mongo storage test config file, %T: %v\n", err, err)
	}
	if address := cfg.GetString("address"); address != "0.0.0.0:1337" {
		t.Fatalf(`Invalid value for "address": %s`, strconv.Quote(address))
	}
	if connectTimeout := cfg.GetDuration("connect_timeout"); connectTimeout != time.Minute+time.Second*30 {
		t.Fatalf(`Invalid value for "connect_timeout": %s`, strconv.Quote(connectTimeout.String()))
	}
	if storageFolder := cfg.GetString("storage_folder"); storageFolder != "./sharex-files/" {
		t.Fatalf(`Invalid value for "storage_folder": %s`, strconv.Quote(storageFolder))
	}
	if authDb := cfg.GetString("auth_db"); authDb != "sharex-admin-db" {
		t.Fatalf(`Invalid value for "auth_db": %s`, strconv.Quote(authDb))
	}
	if authUser := cfg.GetString("auth_user"); authUser != "l_torvalds" {
		t.Fatalf(`Invalid value for "auth_user": %s`, strconv.Quote(authUser))
	}
	if authPasswd := cfg.GetString("auth_passwd"); authPasswd != "MySuperSecurePassword+!#" {
		t.Fatalf(`Invalid value for "auth_passwd": %s`, strconv.Quote(authPasswd))
	}
	if storageDb := cfg.GetString("storage_db"); storageDb != "sharex-upload-metadata" {
		t.Fatalf(`Invalid value for "storage_db": %s`, strconv.Quote(storageDb))
	}
	if storageFileCol := cfg.GetString("storage_file_col"); storageFileCol != "uploads" {
		t.Fatalf(`Invalid value for "storage_file_col": %s`, strconv.Quote(storageFileCol))
	}
}
