package api

import "github.com/linsheng9731/SLB/modules"
import "github.com/linsheng9731/SLB/server"

type Statistic struct {
	FrontendList []FrontendStat `json:"frontends"`
}

type FrontendStat struct {
	Name        string        `json:"name"`
	BackendList []BackendStat `json:"backends"`
}

type BackendStat struct {
	Name           string                 `json:"name"`
	BackendControl modules.BackendControl `json:"backend"`
}

func NewStat(s *server.LbServer) Statistic {
	frontStatList := []FrontendStat{}
	for _, front := range s.FrontendList {
		backendList := []BackendStat{}
		for _, b := range front.BackendList {
			backendList = append(backendList, BackendStat{b.Name, b.BackendControl})
		}
		frontStatList = append(frontStatList, FrontendStat{front.Name, backendList})
	}
	return Statistic{frontStatList}
}
