package modules

import (
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

func (t Table )AddRoute(f *Frontend) *Table {
	var tmp Route
	for _, b := range f.BackendList {
		tmp = Route{
			Service :  b.Name,
			Src     :  b.Name,
			Dst     :  b.Address,
			Weight  :  b.Weight,
		}
		t[f.Port] = append(t[f.Port], tmp)
	}
	return &t
}


func (t Table) Lookup(req *http.Request, port int, pick picker) (route *Route) {
	routes := GetTable()[port]
	r := pick(routes)
	return r
}




