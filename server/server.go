package server

import (
	"crypto/tls"
	"github.com/linsheng9731/slb/config"
	"github.com/linsheng9731/slb/healthcheck"
	"github.com/linsheng9731/slb/modules"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type ShutdownChan chan bool

var holders []*healthcheck.Guard

type LbServer struct {
	config.Configuration
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

	for _, f := range s.FrontendConfigs {
		modules.GetTable().AddRoute(&f)
	}

	for _, f := range s.FrontendConfigs {
		var f = f // catch variable
		log.Println("start to listen frontend " + f.Address())
		go func() {
			h := newHTTPProxy(&f)
			l := modules.Listen{
				Addr:         f.Address(),
				Proto:        "http",
				ReadTimeout:  f.Timeout * time.Millisecond,
				WriteTimeout: f.Timeout * time.Millisecond,
			}
			// listen port
			if err := modules.ListenAndServeHTTP(l, h); err != nil {
				log.Fatal("[FATAL]", err)
			}
		}()
	}

	t := modules.GetTable()
	for _, f := range s.FrontendConfigs {
		g := healthcheck.NewGuard(t, f)
		g.Check()
		holders = append(holders, g)
	}
}

func (s *LbServer) Stop() {
	// stop guards
	for _, h := range holders {
		h.Stop()
	}
}

func newHTTPProxy(f *config.FrontendConfig) http.Handler {

	pick := modules.Picker[f.Strategy] // random target strategy and next target strategy
	if pick == nil {
		log.Print(f.Strategy)
		log.Fatal("strategy is illegal !")
	}
	newTransport := func(tlscfg *tls.Config) *http.Transport {
		return &http.Transport{
			ResponseHeaderTimeout: f.Timeout * time.Millisecond,
			MaxIdleConnsPerHost:   10,
			Dial: (&net.Dialer{
				Timeout:   f.Timeout * time.Millisecond,
				KeepAlive: f.Timeout * time.Millisecond,
			}).Dial,
			TLSClientConfig: tlscfg,
		}
	}
	// http handler
	return &modules.HttpProxy{
		Transport: newTransport(nil),
		Lookup: func(r *http.Request) *modules.Route {
			t := modules.GetTable().Lookup(r, f.Port, pick)
			if t == nil {
				//notFound.Inc(1)
				log.Print("[WARN] No route for ", r.Host, r.URL)
			}
			return t
		},
	}
}
