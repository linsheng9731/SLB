package healthcheck

import (
	"github.com/linsheng9731/slb/modules"
	"log"
	"net/http"
	"time"
)

type BackendHealthCheck struct {
	backend *modules.Backend
}

func NewBackendHealthCheck(b *modules.Backend) BackendHealthCheck {
	return BackendHealthCheck{b}
}

func (check BackendHealthCheck) Check() {
	var b = check.backend
	go func() {
		for {
			select {
			case <-b.Close:
				goto end
			default:
				err, resp := doRequest(b)
				if err != nil || resp.StatusCode >= 400 {
					inActiveMark(b)
				} else {
					activeMark(b)
				}
				if b.Failed {
					time.Sleep(b.RetryTime)
				} else {
					time.Sleep(b.HeartbeatTime)
				}

			}
		}
	end:
		log.Print("Backend health check over.")
	}()
}

func activeMark(b *modules.Backend) {
	b.RWMutex.Lock()
	if b.ActiveTries >= b.ActiveAfter {
		if b.Failed {
			b.Failed = false
			log.Printf("Backend active  [%s]", b.Name)
		}
		b.Active = true
		b.InactiveTries = 0
	} else {
		b.ActiveTries++
	}
	b.RWMutex.Unlock()
}

func inActiveMark(b *modules.Backend) {
	b.RWMutex.Lock()
	if b.InactiveTries >= b.InactiveAfter {
		log.Printf("Backend inactive [%s]", b.Name)
		b.Active = false
		b.ActiveTries = 0
	} else {
		b.Failed = true
		b.InactiveTries++
		log.Printf("Error to check address [%s] name [%s] tries [%d]", b.Heartbeat, b.Name, b.InactiveTries)
	}
	b.RWMutex.Unlock()
}

func doRequest(b *modules.Backend) (error, *http.Response) {
	var request *http.Request
	var err error
	client := &http.Client{}
	request, err = http.NewRequest(b.HBMethod, b.Heartbeat, nil)
	request.Header.Set("User-Agent", "SLB-Heartbeat")
	resp, err := client.Do(request)
	return err, resp
}
