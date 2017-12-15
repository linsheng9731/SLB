package api

import (
	"encoding/json"
	"fmt"
	"github.com/containous/mux"
	"github.com/linsheng9731/slb/common"
	"github.com/linsheng9731/slb/logger"
	"github.com/linsheng9731/slb/modules"
	"github.com/linsheng9731/slb/server"
	"github.com/urfave/negroni"
	"net/http"
)

var lg = logger.Server

type API struct {
	Serer *server.LbServer
	msg   chan int
}

func NewAPI(s *server.LbServer, msg chan int) *API {
	return &API{s, msg}
}

func (api *API) reload(w http.ResponseWriter, r *http.Request) {
	lg.Info("API server get reload configuration request.")
	go func() { api.msg <- common.RELOAD }()
}

func (api *API) check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func (api *API) configuration(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(api.Serer.Configuration, "", "  ")
	if err != nil {
		lg.Error(err)
	}
	fmt.Fprintln(w, string(b))
}

func (api *API) statistic(w http.ResponseWriter, r *http.Request) {
	api.Serer.CalAverageResponse()
	b, err := json.MarshalIndent(api.Serer.Metrics.Facade(), "", "  ")
	if err != nil {
		lg.Error(err)
	}

	fmt.Fprintln(w, string(b))
}

func (api *API) goroutine(w http.ResponseWriter, r *http.Request) {
	modules.ProcessInput("lookup goroutine", w)
}

func (api *API) heap(w http.ResponseWriter, r *http.Request) {
	modules.ProcessInput("lookup heap", w)
}

func (api *API) thread(w http.ResponseWriter, r *http.Request) {
	modules.ProcessInput("lookup threadcreate", w)
}

func (api *API) block(w http.ResponseWriter, r *http.Request) {
	modules.ProcessInput("lookup block", w)
}

func (api *API) gc(w http.ResponseWriter, r *http.Request) {
	modules.ProcessInput("gc summary", w)
}

func (api *API) Listen(address string) {
	var handlerInstance = negroni.New()
	router := mux.NewRouter()
	//router.Methods("GET").Path("/reload").HandlerFunc(api.reload)
	router.Methods("GET").Path("/health").HandlerFunc(api.check)
	router.Methods("GET").Path("/config").HandlerFunc(api.configuration)
	router.Methods("GET").Path("/status").HandlerFunc(api.statistic)
	router.Methods("GET").Path("/profile/goroutine").HandlerFunc(api.goroutine)
	router.Methods("GET").Path("/profile/heap").HandlerFunc(api.heap)
	router.Methods("GET").Path("/profile/thread").HandlerFunc(api.thread)
	router.Methods("GET").Path("/profile/block").HandlerFunc(api.block)
	router.Methods("GET").Path("/profile/gc").HandlerFunc(api.gc)
	handlerInstance.UseHandler(router)
	lg.Info(fmt.Sprintf("Api server listen on %s", address))
	go func() {
		err := http.ListenAndServe(address, handlerInstance)
		if err != nil {
			lg.Fatal(err)
		}
	}()
}
