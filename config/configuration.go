package config

import (
	"fmt"
	"github.com/linsheng9731/slb/common"
	"time"
)

var GlobalConfig = Setup(common.CONFIG_FILENAME)

type GeneralConfig struct {
	LogLevel string `json:"logLevel"`
	LogSize  int    `json:"logSize"`
	Silence  bool   `json:"silence"`
	APIHost  string `json:"apihost"`
	APIPort  int    `json:"apiport"`
}

type FrontendConfig struct {
	Name           string        `json:"name"`
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	Timeout        time.Duration `json:"timeout"`
	Strategy       string        `json:"strategy"`
	BackendsConfig `json:"backends"`
	Heartbeat      string        `json:"heartbeat"`
	HeartbeatTime  time.Duration `json:"heartbeatTime"`
}

func (f *FrontendConfig) Address() string {
	return fmt.Sprintf("%s:%d", f.Host, f.Port)
}

type FrontendConfigs []FrontendConfig

// BackendConfig it's the configuration loaded
type BackendConfig struct {
	Name        string  `json:"name"`
	Hostname    string  `json:"hostname"`
	Weight      float64 `json:"weigth"`
	Address     string  `json:"address"`
	IgnoreCheck bool    `json:"ignoreCheck"`
}

type BackendsConfig []BackendConfig

type Configuration struct {
	GeneralConfig   `json:"general"`
	FrontendConfigs `json:"frontends"`
}

func (c GeneralConfig) Addres() string {
	address := fmt.Sprintf("%s:%d",
		c.APIHost,
		c.APIPort,
	)
	return address
}
