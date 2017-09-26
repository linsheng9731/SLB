package modules

import (
	"github.com/linsheng9731/SLB/config"
	"log"
	"net/http"
	"sync"
	"time"
)

// BackendControl keep the control data
type BackendControl struct {
	Failed bool // The last request failed
	Active bool

	InactiveTries int
	ActiveTries   int
	Score         int
}

// Backend structure
type Backend struct {
	config.BackendConfig
	BackendControl
	sync.RWMutex
	Close chan bool
}

type BackendList []*Backend

func NewBackend(backendConfig config.BackendConfig) *Backend {
	backendConfig.HeartbeatTime = backendConfig.HeartbeatTime * time.Millisecond
	backendConfig.RetryTime = backendConfig.RetryTime * time.Millisecond

	return &Backend{
		BackendConfig: backendConfig,
		BackendControl: BackendControl{
			true, false,
			0, 0, 0,
		},
		Close: make(chan bool, 2),
	}
}

// Monitoring the backend, can add or remove if heartbeat fail
func (b *Backend) HeartCheck() {
	go func() {
		for {
			select {
			case <-b.Close:
				goto end
			default:
				var request *http.Request
				var err error

				client := &http.Client{}
				request, err = http.NewRequest(b.HBMethod, b.Heartbeat, nil)
				request.Header.Set("User-Agent", "SSLB-Heartbeat")

				resp, err := client.Do(request)
				if err != nil || resp.StatusCode >= 400 {
					b.RWMutex.Lock()
					// Max tries before consider inactive
					if b.InactiveTries >= b.InactiveAfter {
						log.Printf("Backend inactive [%s]", b.Name)
						b.Active = false
						b.ActiveTries = 0
					} else {
						// Ok that guy it's out of the game
						b.Failed = true
						b.InactiveTries++
						log.Printf("Error to check address [%s] name [%s] tries [%d]", b.Heartbeat, b.Name, b.InactiveTries)
					}
					b.RWMutex.Unlock()
				} else {
					defer resp.Body.Close()

					// Ok, let's keep working boys
					b.RWMutex.Lock()
					if b.ActiveTries >= b.ActiveAfter {
						if b.Failed {
							log.Printf("Backend active [%s]", b.Name)
						}

						b.Failed = false
						b.Active = true
						b.InactiveTries = 0
					} else {
						b.ActiveTries++
					}
					b.RWMutex.Unlock()
				}

				if b.Failed {
					time.Sleep(b.RetryTime)
				} else {
					time.Sleep(b.HeartbeatTime)
				}
			}
		}
	end:
		log.Print("Backend heart check done.")
	}()
}
