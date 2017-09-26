package server

import (
	"github.com/linsheng9731/SLB/common"
	"github.com/linsheng9731/SLB/modules"
	"log"
	"net/http"
	"time"
)

type HttpHandler struct {
	server   *LbServer
	frontend *modules.Frontend
}

func NewHttpHandler(frontend *modules.Frontend, server *LbServer) HttpHandler {
	return HttpHandler{server, frontend}
}

func (h HttpHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var s = h.server
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
	chanResponse := s.Get(r, h.frontend)
	defer close(chanResponse)
	r.Close = true
	ticker := time.NewTicker(h.frontend.Timeout)
	defer ticker.Stop()
	select {
	case result := <-chanResponse:
		h.HandleResponse(result, w, s, r)

	case <-r.Cancel:
		s.Lock()
		s.Done()
		s.Unlock()

	case <-ticker.C:
		s.Lock()
		s.Done()
		s.Unlock()
		http.Error(w, common.ErrTimeout.Error(), http.StatusRequestTimeout)
	}
}

func (h HttpHandler) HandleResponse(result modules.SLBRequest, w http.ResponseWriter, s *LbServer, r *http.Request) {
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
}
