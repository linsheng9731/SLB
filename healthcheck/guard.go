package healthcheck

import (
	"github.com/linsheng9731/slb/config"
	"github.com/linsheng9731/slb/modules"
	"github.com/lunny/log"
	"net/http"
	"time"
)

type Guard struct {
	table          modules.Table
	frontendConfig config.FrontendConfig
	active         bool
}

func NewGuard(t modules.Table, f config.FrontendConfig) *Guard {
	return &Guard{
		table:          t,
		frontendConfig: f,
		active:         true,
	}
}

func (g Guard)Stop()  {
	g.active = false
}

func (g Guard) Check() {
	interval := g.frontendConfig.HeartbeatTime
	t := time.NewTicker(time.Second * interval )
	routes := g.table[g.frontendConfig.Port]
	go func() {
		for {
			select {
			case <-t.C:
				if !g.active {
					log.Info("receive a kill signal, try to return !")
					goto END
				}
				for i, r := range routes {
					h := r.Dst + g.frontendConfig.Heartbeat
					err, rep := doRequest(h)
					if err != nil || rep.StatusCode >= 400 {
						routes[i].Active = false
						log.Error(r.Src + " " + r.Dst + " is inactive !")
					} else if !r.Active {
						log.Info(r.Src + " " + r.Dst + " is active again !")
						routes[i].Active = true
					} else {
						routes[i].Active = true
					}
				}
			}
		}
	END:
	}()
}

func doRequest(address string) (error, *http.Response) {
	var request *http.Request
	var err error
	client := &http.Client{}
	request, err = http.NewRequest("HEAD", address, nil)
	request.Header.Set("User-Agent", "SLB-Heartbeat")
	resp, err := client.Do(request)
	return err, resp
}
