package modules

import (
	"github.com/linsheng9731/SLB/config"
	"sync"
	"time"
)

type BackendControl struct {
	Failed        bool
	Active        bool
	InactiveTries int
	ActiveTries   int
	Score         int
}

type Backend struct {
	config.BackendConfig
	BackendControl
	sync.RWMutex
	Close chan bool
}

type BackendList []*Backend

func NewBackend(backendConfig config.BackendConfig) *Backend {
	backendConfig.HeartbeatTime = backendConfig.HeartbeatTime * time.Millisecond
	backendConfig.RetryTime = backendConfig.RetryTime * time.Millisecond

	return &Backend{
		BackendConfig: backendConfig,
		BackendControl: BackendControl{
			true, false,
			0, 0, 0,
		},
		Close: make(chan bool, 2),
	}
}
