package server

import (
	"fmt"
	"github.com/linsheng9731/SLB/common"
	"github.com/linsheng9731/SLB/config"
	"github.com/linsheng9731/SLB/healthcheck"
	"github.com/linsheng9731/SLB/modules"
	"log"
	"net/http"
	"runtime"
	"sync"
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
	*modules.WorkerPool
	sync.Mutex
	*sync.WaitGroup
}

func NewServer(configuration config.Configuration) *LbServer {
	return &LbServer{
		Configuration: configuration,
		ShutdownChan:  make(ShutdownChan),
		WaitGroup:     &sync.WaitGroup{},
		WorkerPool:    modules.NewWorkerPool(configuration),
	}
}

func (s *LbServer) Start() {
	log.Println("Setup and check configuration")
	s.setup()

	if len(s.FrontendList) == 0 {
		log.Fatal(errNoFrontend.Error())
	}
	log.Println("Setup ok ...")
	for _, frontend := range s.FrontendList {
		go s.RunFrontendServer(frontend)
	}
}

func (s *LbServer) Stop() {
	log.Println("Showtdun frontends")
	for _, front := range s.FrontendList {
		front.Close <- true
		for _, b := range front.BackendList {
			b.Close <- true
		}
	}
	if s.Configuration.GeneralConfig.GracefulShutdown {
		log.Println("Wait for graceful shutdown")
		s.Wait()
		log.Println("Bye")
	}
}

func (s *LbServer) setup() {
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

func (s *LbServer) RunFrontendServer(frontend *modules.Frontend) {
	var httpHandler = NewHttpHandler(frontend, s)
	if len(frontend.BackendList) == 0 {
		log.Fatal(errNoBackend.Error())
	}
	for _, backend := range frontend.BackendList {
		healthcheck.NewBackendHealthCheck(backend).Check()
	}
	log.Printf("Start frontend http server [%s] at [%s]", frontend.Name, frontend.Address())
	httpHandle := http.NewServeMux()

	httpHandle.HandleFunc(frontend.Route, httpHandler.HandleRequest)
	httpServer := &http.Server{
		Addr:    frontend.Address(),
		Handler: httpHandle,
	}

	go func() {
		<-frontend.Close
		err := httpServer.Shutdown(nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(fmt.Sprintf("Frontend http server %s closed.", frontend.Name))
	}()

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Println(err)
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
