package modules

import (
	"fmt"
	"github.com/linsheng9731/SLB/config"
	"sync"
	"time"
)

type Frontend struct {
	config.FrontendConfig
	BackendList
	sync.RWMutex
	Close chan bool
}

type FrontendList []*Frontend

func NewFrontend(frontendConfig config.FrontendConfig) *Frontend {
	frontendConfig.Timeout = frontendConfig.Timeout * time.Millisecond
	return &Frontend{
		FrontendConfig: frontendConfig,
		Close:          make(chan bool),
	}
}

func (f *Frontend) Address() string {
	return fmt.Sprintf("%s:%d", f.Host, f.Port)
}
