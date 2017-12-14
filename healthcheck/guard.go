package healthcheck

import (
	"github.com/linsheng9731/slb/config"
	"github.com/linsheng9731/slb/logger"
	"github.com/linsheng9731/slb/modules"
	"net/http"
	"time"
)

var lg = logger.Server

type Guard struct {
	table          *modules.Table
	frontendConfig config.FrontendConfig
	active         bool
}

func NewGuard(t *modules.Table, f config.FrontendConfig) *Guard {
	return &Guard{
		table:          t,
		frontendConfig: f,
		active:         true,
	}
}

func (g Guard) Stop() {
	g.active = false
}

func (g Guard) Check() {
	interval := g.frontendConfig.HeartbeatTime
	t := time.NewTicker(time.Second * interval)
	port := g.frontendConfig.Port
	routes := g.table.RoutesMap[port]
	flattenRoutes := flatRoutes(routes)
	go func() {
		for {
			select {
			case <-t.C:
				if !g.active {
					lg.Info("receive a kill signal, try to go end !")
					goto END
				}
				activeRoutes := g.detect(flattenRoutes)
				g.table.ActiveRoutesMap[port] = activeRoutes
			}
		}
	END:
	}()
}

func flatRoutes(routes map[string][]modules.Route) []modules.Route {
	var flattenRoutes []modules.Route
	for _, rr := range routes {
		for _, r := range rr {
			flattenRoutes = append(flattenRoutes, r)
		}
	}

	return flattenRoutes
}

func (g Guard) detect(flattenRoutes []modules.Route) map[string][]modules.Route {
	activeRoutes := make(map[string][]modules.Route)
	for i, r := range flattenRoutes {
		if r.IgnoreCheck {
			flattenRoutes[i].Active = true
			activeRoutes[r.Hostname] = append(activeRoutes[r.Hostname], r)
			continue
		}
		h := r.Dst + g.frontendConfig.Heartbeat
		err, rep := doRequest(h)
		if err != nil || rep.StatusCode >= 400 {
			flattenRoutes[i].Active = false
			lg.Error(r.Src + " " + r.Dst + " is inactive !")
		} else if !r.Active {
			lg.Info(r.Src + " " + r.Dst + " is active again !")
			flattenRoutes[i].Active = true
			activeRoutes[r.Hostname] = append(activeRoutes[r.Hostname], r)
		} else {
			flattenRoutes[i].Active = true
			activeRoutes[r.Hostname] = append(activeRoutes[r.Hostname], r)
		}
	}
	return activeRoutes
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
