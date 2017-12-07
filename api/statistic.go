package api

import "github.com/linsheng9731/slb/server"

type Statistic struct {
	FrontendList []FrontendStat `json:"frontends"`
}

type FrontendStat struct {
	Name        string        `json:"name"`
	BackendList []BackendStat `json:"backends"`
}

type BackendStat struct {
	Name string `json:"name"`
}

func NewStat(s *server.LbServer) Statistic {
	frontStatList := []FrontendStat{}
	for _, front := range s.FrontendConfigs {
		backendList := []BackendStat{}
		for _, b := range front.BackendsConfig {
			backendList = append(backendList, BackendStat{b.Name})
		}
		frontStatList = append(frontStatList, FrontendStat{front.Name, backendList})
	}
	return Statistic{frontStatList}
}
