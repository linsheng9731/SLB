package modules

import (
	"github.com/linsheng9731/slb/config"
	"net/http"
	"sync/atomic"
)

// Table contains a set of routes grouped by host.
// The host routes are sorted from most to least specific
// by sorting the routes in reverse order by path.
type Table map[int][]Route

var table atomic.Value

func init() {
	table.Store(make(Table))
}

func GetTable() Table {
	return table.Load().(Table)
}

func (t Table) AddRoute(f *config.FrontendConfig) *Table {
	var tmp Route
	for _, b := range f.BackendsConfig {
		tmp = NewRoute(b.Name, b.Name, b.Address, b.Weight)
		t[f.Port] = append(t[f.Port], tmp)
	}
	return &t
}

func (t Table) Lookup(req *http.Request, port int, pick picker) (route *Route) {
	routes := GetTable()[port]
	var activeRoutes []Route
	for _, r := range routes {
		if r.Active {
			activeRoutes = append(activeRoutes, r)
		}
	}
	r := pick(activeRoutes)
	return r
}
