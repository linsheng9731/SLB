package api

import (
	"encoding/json"
	"fmt"
	"github.com/containous/mux"
	"github.com/linsheng9731/SLB/common"
	"github.com/linsheng9731/SLB/server"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

type API struct {
	Serer *server.LbServer
	msg   chan int
}

func NewAPI(s *server.LbServer, msg chan int) *API {
	return &API{s, msg}
}

func (api *API) reload(w http.ResponseWriter, r *http.Request) {
	log.Println("API server get reload configuration request.")
	go func() { api.msg <- common.RELOAD }()
}

func (api *API) check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func (api *API) configuration(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(api.Serer.Configuration, "", "  ")
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintln(w, string(b))
}

func (api *API) statistic(w http.ResponseWriter, r *http.Request) {
	stat := NewStat(api.Serer)
	b, err := json.MarshalIndent(stat, "", "  ")
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintln(w, string(b))
}

func (api *API) Listen(address string) {
	var handlerInstance = negroni.New()
	router := mux.NewRouter()
	router.Methods("GET").Path("/reload").HandlerFunc(api.reload)
	router.Methods("GET").Path("/health-check").HandlerFunc(api.check)
	router.Methods("GET").Path("/config").HandlerFunc(api.configuration)
	router.Methods("GET").Path("/status").HandlerFunc(api.statistic)
	handlerInstance.UseHandler(router)
	log.Println(fmt.Sprintf("Api server listen on %s", address))
	go func() {
		err := http.ListenAndServe(address, handlerInstance)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
