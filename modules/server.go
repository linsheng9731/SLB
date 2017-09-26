package modules

import (
	"errors"
	"fmt"
	"github.com/linsheng9731/SLB/config"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

//type Configuration config.Configuration

var (
	errNoFrontend  = errors.New("No frontend configuration detected")
	errNoBackend   = errors.New("No backend configuration detected")
	errTimeout     = errors.New("Timeout")
	errPortExists  = errors.New("Port already in use")
	errRouteExists = errors.New("Route already in use")
)

type ShutdownChan chan bool

type Server struct {
	config.Configuration
	FrontendList
	ShutdownChan
	*WorkerPool
	sync.Mutex
	*sync.WaitGroup
	httpServers []*http.Server
}

func NewServer(configuration config.Configuration) *Server {
	return &Server{
		Configuration: configuration,
		ShutdownChan:  make(ShutdownChan),
		WaitGroup:     &sync.WaitGroup{},
		WorkerPool:    NewWorkerPool(configuration),
		httpServers:   []*http.Server{},
	}
}

func (s *Server) SetConfiguration(configuration config.Configuration) {
	s.Configuration = configuration
}

func (s *Server) setup() {
	runtime.GOMAXPROCS(s.Configuration.GeneralConfig.MaxProcs)

	for _, frontend := range s.Configuration.FrontendsConfig {

		newFrontend := NewFrontend(frontend)
		for _, backend := range frontend.BackendsConfig {
			newFrontend.BackendList = append(newFrontend.BackendList, NewBackend(backend))
		}

		if err := s.preChecksBeforeAdd(newFrontend); err != nil {
			log.Fatal(err.Error())
		} else {
			s.FrontendList = append(s.FrontendList, newFrontend)
		}
	}
}

// Some previous checking before run
func (s *Server) preChecksBeforeAdd(newFrontend *Frontend) error {
	for _, frontend := range s.FrontendList {
		if frontend.Route == newFrontend.Route {
			return errRouteExists
		}

		if frontend.Port == newFrontend.Port {
			return errPortExists
		}

		if len(newFrontend.BackendList) == 0 {
			return errNoBackend
		}
	}

	return nil
}

// Lets run the frontend
func (s *Server) RunFrontendServer(frontend *Frontend) {
	if len(frontend.BackendList) == 0 {
		log.Fatal(errNoBackend.Error())
	}
	for _, backend := range frontend.BackendList {
		backend.HeartCheck()
	}
	log.Printf("Run frontend server [%s] at [%s]", frontend.Name, frontend.Address())
	httpHandle := http.NewServeMux()
	requestHandler := func(w http.ResponseWriter, r *http.Request) {
		s.Lock()
		s.Add(1)
		s.Unlock()

		defer func() {
			if rec := recover(); rec != nil {
				log.Println("Err", rec)
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		// Get a channel the already attached to a worker
		chanResponse := s.Get(r, frontend)
		defer close(chanResponse)
		r.Close = true
		ticker := time.NewTicker(frontend.Timeout)
		defer ticker.Stop()
		select {
		case result := <-chanResponse:
			w, s, r = handleResponse(result, w, s, r)

		case <-r.Cancel:
			s.Lock()
			s.Done()
			s.Unlock()

		case <-ticker.C:
			s.Lock()
			s.Done()
			s.Unlock()
			http.Error(w, errTimeout.Error(), http.StatusRequestTimeout)
		}
	}

	httpHandle.HandleFunc(frontend.Route, requestHandler)
	server := &http.Server{
		Addr:    frontend.Address(),
		Handler: httpHandle,
	}
	s.httpServers = append(s.httpServers, server)
	go func() {
		<-frontend.Close
		server.Shutdown(nil)
		log.Print(fmt.Sprintf("Frontend closing: %s", frontend.Name))
	}()
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func handleResponse(result SSLBRequest, w http.ResponseWriter, s *Server, r *http.Request) (http.ResponseWriter, *Server, *http.Request) {
	// todo  it's valid ?
	for k, vv := range result.Header {
		for _, v := range vv {
			w.Header().Set(k, v)
		}
	}
	s.Lock()
	s.Done()
	s.Unlock()
	if result.Upgraded {
		if s.Configuration.GeneralConfig.Websocket {
			result.HijackWebSocket(w, r)
		}
	} else {
		w.WriteHeader(result.Status)

		if r.Method != "HEAD" {
			w.Write(result.Body)
		}
	}
	return w, s, r
}

func (s *Server) Run() {
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

func (s *Server) Stop() {
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
