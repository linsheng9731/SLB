package modules

type Route struct {
	Service string  `json:"service"`
	Src     string  `json:"src"`
	Dst     string  `json:"dst"`
	Weight  float64 `json:"weight"`
	Active  bool
}

func NewRoute(s string, sr string, dst string, weight float64) Route {
	return Route{
		Service: s,
		Src:     sr,
		Dst:     dst,
		Weight:  weight,
		Active:  true,
	}
}
