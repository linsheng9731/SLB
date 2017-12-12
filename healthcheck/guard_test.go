package healthcheck

import (
	"github.com/linsheng9731/slb/modules"
	"testing"
)

func TestGuardFlatRoutes(t *testing.T) {
	r1 := modules.NewRoute("", "", "", "back1", 1)
	r2 := modules.NewRoute("", "", "", "back2", 1)
	r3 := modules.NewRoute("", "", "", "back3", 1)
	routesMap := make(map[string][]modules.Route)
	routesMap["host1"] = []modules.Route{r1}
	routesMap["host2"] = []modules.Route{r2, r3}
	flattenRoutes := flatRoutes(routesMap)
	if len(flattenRoutes) != 3 {
		t.Fatal("flatRoutes return wrong reuslt !")
	}
}
