package config

import (
	"fmt"
	"time"
)

type GeneralConfig struct {
	MaxProcs       int    `json:"maxProcs"`
	WorkerPoolSize int    `json:"workerPoolSize"`
	Websocket      bool   `json:"websocket"`
	LogLevel       string `json:"logLevel"` // Need to define how it works
	APIHost        string `json:"apihost"`
	APIPort        int    `json:"apiport"`
}

type FrontendConfig struct {
	Name           string        `json:"name"`
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	Route          string        `json:"route"`
	Timeout        time.Duration `json:"timeout"`
	Strategy       string        `json:"strategy"`
	BackendsConfig `json:"backends"`
	Heartbeat      string        `json:"heartbeat"`
	HBMethod       string        `json:"hbmethod"`
	ActiveAfter    int           `json:"activeAfter"`
	InactiveAfter  int           `json:"inactiveAfter"` // Consider inactive after max inactiveAfter
	HeartbeatTime  time.Duration `json:"heartbeatTime"` // Heartbeat time if health
	RetryTime      time.Duration `json:"retryTime"`     // Retry to time after failed
}

func (f *FrontendConfig) Address() string {
	return fmt.Sprintf("%s:%d", f.Host, f.Port)
}

type FrontendConfigs []FrontendConfig

// BackendConfig it's the configuration loaded
type BackendConfig struct {
	Name    string  `json:"name"`
	Weight  float64 `json:"weigth"`
	Address string  `json:"address"`
}

type BackendsConfig []BackendConfig

type Configuration struct {
	GeneralConfig   `json:"general"`
	FrontendConfigs `json:"frontends"`
}

func (c GeneralConfig) APIAddres() string {
	address := fmt.Sprintf("%s:%d",
		c.APIHost,
		c.APIPort,
	)
	return address
}
