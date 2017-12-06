package modules

type Route struct {
	Service string            `json:"service"`
	Src     string            `json:"src"`
	Dst     string            `json:"dst"`
	Weight  float64           `json:"weight"`
}
