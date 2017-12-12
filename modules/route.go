package modules

type Route struct {
	Service  string  `json:"service"`
	Src      string  `json:"src"`
	Dst      string  `json:"dst"`
	Hostname string  `json:"hostname"`
	Weight   float64 `json:"weight"`
	Active   bool
}

func NewRoute(s string, sr string, hostname string, dst string, weight float64) Route {
	return Route{
		Service:  s,
		Src:      sr,
		Dst:      dst,
		Hostname: hostname,
		Weight:   weight,
		Active:   true,
	}
}
