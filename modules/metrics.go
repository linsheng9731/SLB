package modules

import (
	"math"
	"sync/atomic"
	"time"
)

var Cnt int64 = 0
var DurResults = [10]float64{}

// Metrics use to record some useful metrics
type Metrics struct {

	// record input traffic of service per second
	InTotalTraffic int64

	InAverageTraffic int64

	// record output traffic of service per second
	OutTotalTraffic int64

	OutAverageTraffic int64

	// record average spent time (ms) of response
	AverageResponseTime float64

	// count of 2xx http status code
	OkCount int64

	// count of 5xx http status code
	FailedCount int64
}

type MetricsFacade struct {
	// record input traffic of service per second
	InTraffic int64

	// record output traffic of service per second
	OutTraffic int64

	// record average spent time (ms) of response
	AverageResponseTime float64

	// count of 2xx http status code
	OkCount int64

	// count of 5xx http status code
	FailedCount int64

	// count of total request
	RequestCount int64
}

func (m *Metrics) Facade() MetricsFacade {
	return MetricsFacade{
		m.InAverageTraffic,
		m.OutAverageTraffic,
		m.AverageResponseTime,
		m.OkCount,
		m.FailedCount,
		m.OkCount + m.FailedCount,
	}
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) AddOkCount() {
	atomic.AddInt64(&m.OkCount, 1)
}

func (m *Metrics) AddFailedCount() {
	atomic.AddInt64(&m.FailedCount, 1)
}

func (m *Metrics) AddInTraffic(t int64) {
	atomic.AddInt64(&m.InTotalTraffic, t)
}

func (m *Metrics) AddOutTraffic(t int) {
	atomic.AddInt64(&m.OutTotalTraffic, int64(t))
}

func (m *Metrics) CalAverageResponse() {
	sum := 0.0
	cnt := 0.0
	for _, r := range DurResults {
		if r == 0 {
			continue
		}
		sum += r
		cnt++
	}
	if cnt == 0 {
		return
	}
	m.AverageResponseTime = math.Trunc(sum*1000) / cnt
}

func (m *Metrics) Status(s int) {
	if s > 400 {
		m.AddFailedCount()
	} else {
		m.AddOkCount()
	}
}

func (m *Metrics) IntervalTask() {
	timer := time.NewTicker(time.Second * 3)
	go func() {
		for {
			select {
			case <-timer.C:
				Cnt = 0
				m.InAverageTraffic = m.InTotalTraffic / 3
				m.OutAverageTraffic = m.OutTotalTraffic / 3
				m.InTotalTraffic = 0
				m.OutTotalTraffic = 0
			}
		}
	}()
}

func (m *Metrics) Dur(d float64) {
	atomic.AddInt64(&Cnt, 1)
	DurResults[Cnt%10] = d
}
