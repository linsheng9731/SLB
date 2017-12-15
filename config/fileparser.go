package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	lg "log"
	"time"
)

const DEFAULT_FILENAME = "config.json"

func openFile(filename string) []byte {
	var file []byte
	var err error

	if filename != "" {
		file, err = ioutil.ReadFile(filename)
		if err == nil {
			return file
		} else {
			lg.Fatal(err)
		}
	}

	file, err = ioutil.ReadFile("/etc/sslb/" + DEFAULT_FILENAME)
	if err == nil {
		return file
	}

	file, err = ioutil.ReadFile("~/./sslb/" + DEFAULT_FILENAME)
	if err == nil {
		return file
	}

	file, err = ioutil.ReadFile("./" + DEFAULT_FILENAME)
	if err != nil {
		lg.Fatal("No config file found, in /etc/sslb or ~/.sslb or in current dir")
	}

	return file
}

// ConfParser to Parse JSON FILE
func ConfParser(file []byte) *Configuration {
	if err := Validate(file); err != nil {
		lg.Fatal("Can't validate config.json ", err)
	}

	jsonConfig := Configuration{
		GeneralConfig: GeneralConfig{
			LogLevel: "info",
			APIHost:  "127.0.0.1",
			APIPort:  9292,
		},
		FrontendConfigs: []FrontendConfig{
			{
				Timeout:        time.Millisecond * 30000,
				HeartbeatTime:  time.Millisecond * 30000,
				BackendsConfig: []BackendConfig{{Weight: 1, Hostname: "", IgnoreCheck: false}},
			},
		},
	}

	err := json.Unmarshal(file, &jsonConfig)

	if err != nil {
		lg.Fatal("Error to parse json conf", err.Error())
	}

	return &jsonConfig
}

// Setup parse config file and return configuration
func Setup(filename string) *Configuration {
	if v := flag.Lookup("test.v"); v == nil {
		file := openFile(filename)
		return ConfParser(file)
	} else {
		file := openFile("../config.json")
		return ConfParser(file)
	}
}
