package config

import (
	"testing"
	"time"
)

func TestFileparserGeneral(t *testing.T) {
	jsonConf := []byte(`
    {
        "general": {
            "maxProcs": 2,
            "gracefulShutdown": true,
            "logLevel": "info",
            "websocket": true
        },
        "frontends" : [
            {
                "name" : "Front1",
                "host" : "127.0.0.1",
                "port" : 9000,
                "route" : "/",
                "timeout" : 5000,
				"strategy": "rnd",
				"inactiveAfter" : 3,
				"activeAfter" : 1,
				"heartbeatTime" : 5000,
				"retryTime" : 1000,
				"hbmethod" : "HEAD",
				"heartbeat" : "http://127.0.0.1:9001"
            }
        ]
    }`)

	conf := ConfParser(jsonConf)

	if conf.GeneralConfig.LogLevel != "info" {
		t.Fatal("LogLevel is wrong", conf.GeneralConfig.LogLevel)
	}

}

func TestFileparserFrontend(t *testing.T) {
	jsonConf := []byte(`
    {
        "general": {
            "maxProcs": 4,
            "workerPoolSize": 1000,
            "gracefulShutdown": true,
            "logLevel": "info",
            "websocket": true
        },
        "frontends" : [
            {
                "name" : "Front1",
                "host" : "127.0.0.1",
                "port" : 9000,
                "route" : "/",
				"strategy": "rnd",
				"inactiveAfter" : 3,
				"activeAfter" : 1,
				"heartbeatTime" : 5000,
				"retryTime" : 1000,
				"hbmethod" : "HEAD",
				"heartbeat" : "http://127.0.0.1:9001"
            }
        ]
    }`)

	conf := ConfParser(jsonConf)
	if conf.FrontendConfigs[0].Name != "Front1" {
		t.Fatal("Name is wrong", conf.FrontendConfigs[0].Name)
	}

	if conf.FrontendConfigs[0].Host != "127.0.0.1" {
		t.Fatal("Host is wrong", conf.FrontendConfigs[0].Host)
	}

	if conf.FrontendConfigs[0].Port != 9000 {
		t.Fatal("Port is wrong", conf.FrontendConfigs[0].Port)
	}

	timeout := time.Millisecond * 30000
	if conf.FrontendConfigs[0].Timeout != timeout {
		t.Fatal("Timeout is wrong", conf.FrontendConfigs[0].Timeout)
	}
}

func TestFileparserBackend(t *testing.T) {
	jsonConf := []byte(`
    {
        "general": {
            "maxProcs": 4,
            "workerPoolSize": 1000,
            "gracefulShutdown": true,
            "logLevel": "info",
            "websocket": true,
            "host": "127.0.0.1",
            "port": 42555
        },

        "frontends" : [
            {
                "name" : "Front1",
                "host" : "127.0.0.1",
                "port" : 9000,
                "route" : "/",
                "timeout" : 5000,
				"strategy": "rnd",
				"inactiveAfter" : 3,
				"activeAfter" : 1,
				"heartbeatTime" : 5000,
				"retryTime" : 1000,
				"hbmethod" : "HEAD",
				"heartbeat" : "http://127.0.0.1:9001",
                "backends" : [
                    {
                        "name" : "Back1",
                        "address" : "http://127.0.0.1:9001",
                        "weigth": 1,
						"ignoreCheck": true
                    }
                ]
            }
        ]
    }`)

	conf := ConfParser(jsonConf)
	if conf.FrontendConfigs[0].BackendsConfig[0].Name != "Back1" {
		t.Fatal("Name is wrong", conf.FrontendConfigs[0].BackendsConfig[0].Name)
	}

	if conf.FrontendConfigs[0].BackendsConfig[0].IgnoreCheck != true {
		t.Fatal("Ignore check is wrong", conf.FrontendConfigs[0].BackendsConfig[0].IgnoreCheck)
	}
}
