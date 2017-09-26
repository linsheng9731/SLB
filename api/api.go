package api

import (
	"fmt"
	"github.com/containous/mux"
	"github.com/linsheng9731/SLB/common"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

const (
	CONFIG_FILENAME = "config.json"
)

type API struct {
	msg chan int
}

func NewAPI(msg chan int) *API {
	return &API{msg}
}

func (api *API) reload(response http.ResponseWriter, request *http.Request) {
	log.Println("API server get reload configuration request.")
	go func() { api.msg <- common.RELOAD }()
}

func (api *API) Listen(address string) {
	var handlerInstance = negroni.New()
	router := mux.NewRouter()
	router.Methods("GET").Path("/reload").HandlerFunc(api.reload)
	handlerInstance.UseHandler(router)
	log.Println(fmt.Sprintf("Api server listen on %s", address))
	go func() {
		err := http.ListenAndServe(address, handlerInstance)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
