package modules

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Worker struct {
	Idle bool
	sync.RWMutex
}

type Workers []*Worker

func NewWorker() *Worker {
	return &Worker{Idle: true}
}

func processReturn(result *http.Response) SLBRequest {
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return NewWorkerRequest(http.StatusInternalServerError, result.Header, []byte(err.Error()))
	}

	return NewWorkerRequest(result.StatusCode, result.Header, []byte(body))
}

func checkForWebsocket(r *http.Request) bool {
	result := false
	connHDR := ""
	connHDRS := r.Header["Connection"]

	if len(connHDRS) > 0 {
		connHDR = connHDRS[0]
	}

	if connHDR == "upgrade" || connHDR == "Upgrade" {
		upgradeHDRS := r.Header["Upgrade"]
		if len(upgradeHDRS) > 0 {
			result = (strings.ToLower(upgradeHDRS[0]) == "websocket")
		}
	}

	return result
}

func execRequest(backend *Backend, r *http.Request) SLBRequest {
	var httpRequest *http.Request
	var err error

	if checkForWebsocket(r) {
		ret := NewWorkerRequestUpgraded()
		ret.Backend = backend
		return ret
	}

	requestAddress := backend.BackendConfig.Address + r.URL.String()

	client := &http.Client{}
	httpRequest, err = http.NewRequest(r.Method, requestAddress, r.Body)

	for k, vv := range r.Header {
		for _, v := range vv {
			httpRequest.Header.Set(k, v)
		}
	}

	response, err := client.Do(httpRequest)

	if err != nil {
		return NewWorkerRequestErr(http.StatusRequestTimeout, []byte("No backend available"))
	}

	if response == nil {
		return NewWorkerRequestErr(http.StatusBadGateway, []byte("Method Not Supported By SLB"))
	}

	ret := processReturn(response)
	ret.Backend = backend
	return ret
}

func searchBackend(frontend *Frontend) *Backend {
	frontend.Lock()
	defer frontend.Unlock()
	var backendWithMinScore *Backend
	for idx, backend := range frontend.BackendList {
		if backend.Failed {
			continue
		}
		backend.RLock()
		if idx == 0 {
			backendWithMinScore = backend
		} else {
			if backend.Score < backendWithMinScore.Score {
				backendWithMinScore = backend
			}
		}
		backend.RUnlock()
	}
	return backendWithMinScore
}

func (w *Worker) Run(r *http.Request, frontend *Frontend) SLBRequestChan {
	w.Lock()
	w.Idle = false
	w.Unlock()

	chanReceiver := make(SLBRequestChan)
	go func(w *Worker, c SLBRequestChan, f *Frontend) {
		defer func() {
			if rec := recover(); rec != nil {
				// Channel is closed can happen
			}
		}()

		backend := searchBackend(f)

		if backend != nil {
			backend.Lock()
			backend.Score++
			backend.Unlock()
			c <- execRequest(backend, r)
		} else {
			chanReceiver <- NewWorkerRequestErr(http.StatusServiceUnavailable, []byte("Service Unavailable"))
		}

		w.Lock()
		w.Idle = true
		w.Unlock()
	}(w, chanReceiver, frontend)

	return chanReceiver
}
