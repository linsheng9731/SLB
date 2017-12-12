package modules

import "testing"

func TestRndPicker(t *testing.T) {
	r1 := NewRoute("", "", "", "back1", 1)
	r2 := NewRoute("", "", "", "back2", 1)
	r3 := NewRoute("", "", "", "back3", 1)
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
	r1 := NewRoute("", "", "", "back1", 1)
	r2 := NewRoute("", "", "", "back2", 1)
	r3 := NewRoute("", "", "", "back3", 1)
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
