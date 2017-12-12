package modules

import (
	"github.com/linsheng9731/slb/config"
	"log"
	"net/http"
	"testing"
	"time"
)

var (
	cf config.Configuration
	tt *Table
)

func init() {
	cf = config.Configuration{
		GeneralConfig: config.GeneralConfig{
			Websocket: true,
			LogLevel:  "info",
			APIHost:   "127.0.0.1",
			APIPort:   9292,
		},
		FrontendConfigs: []config.FrontendConfig{
			{
				Host:          "0.0.0.0",
				Port:          9000,
				Timeout:       time.Millisecond * 30000,
				HeartbeatTime: time.Millisecond * 30000,
				BackendsConfig: []config.BackendConfig{
					{Weight: 1, Hostname: "host0.com"},
					{Weight: 1, Hostname: "host1.com"},
					{Weight: 1, Hostname: "host1.com"},
					{Weight: 1, Hostname: "host2.com"},
					{Weight: 1, Hostname: "host2.com"},
					{Weight: 1, Hostname: "host2.com"},
					{Weight: 1, Hostname: ""},
				},
			},
		},
	}

	for _, f := range cf.FrontendConfigs {
		tt = GetTable().AddRoute(&f)
	}
}

func TestTable_AddRoute(t *testing.T) {
	ok, activeMaps := tt.GetActiveRoutes(9000)
	if !ok {
		t.Fatal("Get active routes failed!")
	}
	if len(activeMaps["host0.com"]) != 1 {
		log.Fatal("Get active routes for 'host.com0' failed! ")
	}
	if len(activeMaps["host1.com"]) != 2 {
		log.Fatal("Get active routes for 'host.com1' failed! ")
	}
	if len(activeMaps["host2.com"]) != 3 {
		log.Fatal("Get active routes for 'host.com2' failed! ")
	}
	if len(activeMaps[""]) != 1 {
		log.Fatal("Get default active routes failed! ")
	}
}

func TestTable_Lookup(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = ""
	r := tt.Lookup(req, 9000, rndPicker)
	if r.Hostname != "" {
		t.Fatal("Table lookup route failed!")
	}
	req.Host = "hostx.com"
	r = tt.Lookup(req, 9000, rndPicker)
	if r.Hostname != "" {
		t.Fatal("Table lookup route failed!")
	}
	req.Host = "host1.com"
	r = tt.Lookup(req, 9000, rndPicker)
	if r.Hostname != "host1.com" {
		t.Fatal("Table lookup route failed!")
	}
}
