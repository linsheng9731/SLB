package modules

import (
	"github.com/linsheng9731/slb/config"
	"net/http"
	"strings"
	"sync/atomic"
)

// Table contains a set of RoutesMap grouped by host.
// The host RoutesMap are sorted from most to least specific
// by sorting the RoutesMap in reverse order by path.
type Table struct {
	RoutesMap       map[int]map[string][]Route
	ActiveRoutesMap map[int]map[string][]Route
}

var table atomic.Value

func init() {
	table.Store(Table{make(map[int]map[string][]Route), make(map[int]map[string][]Route)})
}

func GetTable() Table {
	return table.Load().(Table)
}

func (t Table) GetActiveRoutes(port int) (bool, map[string][]Route) {
	routes, ok := t.ActiveRoutesMap[port]
	if ok {
		return true, routes
	} else {
		return false, make(map[string][]Route)
	}
}

func (t Table) AddRoute(f *config.FrontendConfig) *Table {
	var route Route
	routesMap := make(map[string][]Route)
	for _, b := range f.BackendsConfig {
		route = NewRoute(b.Name, b.Name, b.Hostname, b.Address, b.Weight)
		routes, ok := routesMap[b.Hostname]
		if ok {
			routesMap[b.Hostname] = append(routes, route)
		} else {
			// first one initial a slice
			routesMap[b.Hostname] = []Route{route}
		}
	}
	t.RoutesMap[f.Port] = routesMap
	t.ActiveRoutesMap = t.RoutesMap
	return &t
}

// Lookup function return a suitable route for request .
// If one route has hostname attribute, lookup check if the
// host of request and route's hostname are equal, finally
// the picker you specified will pick one suitable route.
func (t Table) Lookup(req *http.Request, port int, pick picker) (route *Route) {
	hostname := host(req.Host)
	var r *Route
	activeRoutes, ok := t.ActiveRoutesMap[port][hostname]
	if ok {
		r = pick(activeRoutes)
	} else {
		// default host
		routesMap, ok := t.ActiveRoutesMap[port][""]
		if ok {
			r = pick(routesMap)
		} else {
			r = nil
		}
	}
	return r
}

func host(hostname string) string {
	return strings.Split(hostname, ":")[0]
}
