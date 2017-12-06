package config

import (
	"fmt"
	"time"
)

type GeneralConfig struct {
	MaxProcs         int    `json:"maxProcs"`
	WorkerPoolSize   int    `json:"workerPoolSize"`
	GracefulShutdown bool   `json:"gracefulShutdown"`
	Websocket        bool   `json:"websocket"`
	LogLevel         string `json:"logLevel"` // Need to define how it works
	RPCHost          string `json:"rpchost"`
	RPCPort          int    `json:"rpcport"`
	APIHost          string `json:"apihost"`
	APIPort          int    `json:"apiport"`
}

type FrontendConfig struct {
	Name           string        `json:"name"`
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	Route          string        `json:"route"`
	Timeout        time.Duration `json:"timeout"`
	Strategy       string        `json:"strategy"`
	BackendsConfig `json:"backends"`
}
type FrontendsConfig []FrontendConfig

// BackendConfig it's the configuration loaded
type BackendConfig struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Heartbeat string `json:"heartbeat"`
	HBMethod  string `json:"hbmethod"`

	ActiveAfter   int `json:"activeAfter"`
	InactiveAfter int `json:"inactiveAfter"` // Consider inactive after max inactiveAfter
	Weight        float64 `json:"weigth"`

	HeartbeatTime time.Duration `json:"heartbeatTime"` // Heartbeat time if health
	RetryTime     time.Duration `json:"retryTime"`     // Retry to time after failed
}
type BackendsConfig []BackendConfig

type Configuration struct {
	GeneralConfig   `json:"general"`
	FrontendsConfig `json:"frontends"`
}

func (c GeneralConfig) RPCAddres() string {
	address := fmt.Sprintf("%s:%d",
		c.RPCHost,
		c.RPCPort,
	)
	return address
}

func (c GeneralConfig) APIAddres() string {
	address := fmt.Sprintf("%s:%d",
		c.APIHost,
		c.APIPort,
	)
	return address
}
