package modules

import (
	"github.com/linsheng9731/SLB/config"
	"log"
	"net/http"
	"sync"
	"time"
)

type WorkerPool struct {
	config.Configuration
	Workers
	sync.RWMutex
}

func NewWorkerPool(configuration config.Configuration) *WorkerPool {
	wp := &WorkerPool{
		Configuration: configuration,
	}

	wp.createPool()
	return wp
}

func (wp *WorkerPool) createPool() {
	log.Printf("Create worker pool with [%d]", wp.Configuration.WorkerPoolSize)
	for i := 0; i <= wp.Configuration.WorkerPoolSize; i++ {
		worker := NewWorker()
		wp.Workers = append(wp.Workers, worker)
	}
}

// todo experiment feature
func (wp *WorkerPool) Resize() {
	log.Printf("Resize worker pool with [%d", wp.Configuration.WorkerPoolSize)
	if len(wp.Workers) > wp.Configuration.WorkerPoolSize {
		wp.Lock()
		defer wp.Unlock()
		wp.Workers = []*Worker{}
		for i := 0; i <= wp.Configuration.WorkerPoolSize; i++ {
			worker := NewWorker()
			wp.Workers = append(wp.Workers, worker)
		}

	} else if len(wp.Workers) < wp.Configuration.WorkerPoolSize {
		for i := 0; i <= wp.Configuration.WorkerPoolSize-len(wp.Workers); i++ {
			worker := NewWorker()
			wp.Workers = append(wp.Workers, worker)
		}
	}
}

func (wp *WorkerPool) CountIdle() int {
	count := 0

	for _, worker := range wp.Workers {
		worker.RLock()
		if worker.Idle {
			count++
		}
		worker.RUnlock()
	}

	return count
}

func (wp *WorkerPool) Get(r *http.Request, frontend *Frontend) SLBRequestChan {
	wp.Lock()
	defer wp.Unlock()

	var idleWorker *Worker

	for {

		for _, worker := range wp.Workers {
			worker.Lock()
			if worker.Idle {
				worker.Idle = false
				idleWorker = worker
				worker.Unlock()
				break
			}
			worker.Unlock()
		}

		if idleWorker == nil {
			idleWorker = NewWorker()
			idleWorker.Lock()
			idleWorker.Idle = false
			idleWorker.Unlock()
		}

		if idleWorker != nil {
			break
		}

		time.Sleep(time.Millisecond)
	}

	c := idleWorker.Run(r, frontend)
	return c
}
