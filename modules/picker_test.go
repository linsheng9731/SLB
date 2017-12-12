package modules

import "testing"

var (
	r1 Route
	r2 Route
	r3 Route
)

func init() {

	r1 = NewRoute("", "", "", "back1", false, 1)
	r2 = NewRoute("", "", "", "back2", false, 1)
	r3 = NewRoute("", "", "", "back3", false, 1)
}
func TestRndPicker(t *testing.T) {

	routes1 := []Route{}
	routes2 := []Route{r1}
	routes3 := []Route{r1, r2, r3}
	if rndPicker(routes1) != nil {
		t.Fatal("picker return wrong result !")
	}

	if rndPicker(routes2).Dst != "back1" {
		t.Fatal("picker return wrong result !")
	}

	if rndPicker(routes3) == nil {
		t.Fatal("picker return wrong result !")
	}
}

func TestRrdPicker(t *testing.T) {
	routes1 := []Route{}
	routes2 := []Route{r1}
	routes3 := []Route{r1, r2, r3}
	if rrPicker(routes1) != nil {
		t.Fatal("picker return wrong result !")
	}

	if rrPicker(routes2).Dst != "back1" {
		t.Fatal("picker return wrong result !")
	}

	if rrPicker(routes3) == nil {
		t.Fatal("picker return wrong result !")
	}
}
