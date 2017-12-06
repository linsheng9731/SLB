package server

import (
	"github.com/linsheng9731/slb/common"
	"github.com/linsheng9731/slb/config"
	"github.com/linsheng9731/slb/modules"
	"log"
	"net/http"
	"runtime"
	"sync"
	"crypto/tls"
	"net"
)

var (
	errNoFrontend = common.ErrNoFrontend
	errNoBackend  = common.ErrNoBackend
	errPortExists = common.ErrPortExists
)

type ShutdownChan chan bool

type LbServer struct {
	config.Configuration
	modules.FrontendList
	ShutdownChan
	sync.Mutex
	*sync.WaitGroup
}

func NewServer(configuration config.Configuration) *LbServer {
	return &LbServer{
		Configuration: configuration,
		ShutdownChan:  make(ShutdownChan),
		WaitGroup:     &sync.WaitGroup{},
	}
}

func (s *LbServer) Run() {
	for _, f := range s.FrontendList {
		modules.GetTable().AddRoute(f)
		h := newHTTPProxy(f)
		l :=  modules.Listen {
			f.Address(),
			"http",
			f.Timeout,
			f.Timeout,
		}
		// listen port
		if err := modules.ListenAndServeHTTP(l, h); err != nil {
			log.Fatal("[FATAL]", err)
		}
	}
}

func (s *LbServer) Stop() {

}

func (s *LbServer) Setup() {
	runtime.GOMAXPROCS(s.Configuration.GeneralConfig.MaxProcs)

	for _, frontend := range s.Configuration.FrontendsConfig {

		newFrontend := modules.NewFrontend(frontend)
		for _, backend := range frontend.BackendsConfig {
			newFrontend.BackendList = append(newFrontend.BackendList, modules.NewBackend(backend))
		}

		if err := s.preChecksBeforeAdd(newFrontend); err != nil {
			log.Fatal(err.Error())
		} else {
			s.FrontendList = append(s.FrontendList, newFrontend)
		}
	}
}

func newHTTPProxy(f *modules.Frontend) http.Handler {

	pick := modules.Picker[f.Strategy] // random target strategy and next target strategy
	if pick == nil {
		log.Print(f.Strategy)
		log.Fatal("strategy is illegal !")
	}
	newTransport := func(tlscfg *tls.Config) *http.Transport {
		return &http.Transport{
			ResponseHeaderTimeout: f.Timeout,
			MaxIdleConnsPerHost:   10,
			Dial: (&net.Dialer{
				Timeout:   f.Timeout,
				KeepAlive: f.Timeout,
			}).Dial,
			TLSClientConfig: tlscfg,
		}
	}
	// http handler
	return &modules.HttpProxy{
		Transport:         newTransport(nil),
		Lookup: func(r *http.Request) *modules.Route {
			t := modules.GetTable().Lookup(r,f.Port, pick)
			if t == nil {
				//notFound.Inc(1)
				log.Print("[WARN] No route for ", r.Host, r.URL)
			}
			return t
		},
	}
}

func (s *LbServer) preChecksBeforeAdd(newFrontend *modules.Frontend) error {
	for _, frontend := range s.FrontendList {
		if frontend.Port == newFrontend.Port {
			return errPortExists
		}
		if len(newFrontend.BackendList) == 0 {
			return errNoBackend
		}
	}
	return nil
}
