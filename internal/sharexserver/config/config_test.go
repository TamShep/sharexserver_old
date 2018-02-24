package config

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestMainConfig(t *testing.T) {
	cfg, err := loadCfg("../../../test/test-config.toml")
	if err != nil {
		t.Fatalf("Could not load test config file, %T: %v", err, err)
	}
	if webServerAddress := cfg.GetString("webserver_address"); webServerAddress != ":80" {
		t.Fatalf(`Invalid value for "webserver_address": %s`, strconv.Quote(webServerAddress))
	}
	if storageEngine := cfg.GetString("storage_engine"); storageEngine != "Magic+fairy_dust" {
		t.Fatalf(`Invalid value for "storage_engine": %s`, strconv.Quote(storageEngine))
	}
	if storageEngineConfig := cfg.GetString("storage_engine_config"); storageEngineConfig != "./mongo-storage-config.toml" {
		t.Fatalf(`Invalid value for "storage_engine_config": %s`, strconv.Quote(storageEngineConfig))
	}
	if reverseProxyHeader := cfg.GetString("reverse_proxy_header"); reverseProxyHeader != "This-Header-Contains-The-Real-IP" {
		t.Fatalf(`Invalid value for "reverse_proxy_header": %s`, strconv.Quote(reverseProxyHeader))
	}
	if whitelistedContentTypes := cfg.GetStringSlice("whitelisted_content_types"); !reflect.DeepEqual(whitelistedContentTypes, []string{"first-ct", "a-mime-type", "sp€ci4l"}) {
		t.Fatalf(`Invalid value for "whitelisted_content_types": %s`, strconv.Quote(fmt.Sprintf("%+v", whitelistedContentTypes)))
	}
}
