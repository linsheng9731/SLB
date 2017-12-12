package modules

import (
	"sync/atomic"
	"time"
)

// picker selects a target from a list of targets.
type picker func([]Route) *Route

var total uint64

// Picker contains the available picker functions.
// Update config/load.go#load after updating.
var Picker = map[string]picker{
	"rnd": rndPicker,
	"rr":  rrPicker,
}

// rndPicker picks a random target from the list of targets.
func rndPicker(rs []Route) *Route {
	if len(rs) == 0 {
		return nil
	}
	return &rs[randIntn(len(rs))]
}

// rrPicker picks the next target from a list of targets using round-robin.
func rrPicker(rs []Route) *Route {
	if len(rs) == 0 {
		return nil
	}
	u := rs[total%uint64(len(rs))]
	atomic.AddUint64(&total, 1)
	if total > 9223372036854775807 {
		total = 0
	}
	return &u
}

// stubbed out for testing
// we implement the randIntN function using the nanosecond time counter
// since it is 15x faster than using the pseudo random number generator
// (12 ns vs 190 ns) Most HW does not seem to provide clocks with ns
// resolution but seem to be good enough for µs resolution. Since
// requests are usually handled within several ms we should have enough
// variation. Within 1 ms we have 1000 µs to distribute among a smaller
// set of entities (<< 100)
var randIntn = func(n int) int {
	if n == 0 {
		return 0
	}
	return int(time.Now().UnixNano()/int64(time.Microsecond)) % n
}
