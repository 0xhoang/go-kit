package common

import (
	"sync"
)

type Worker interface {
	Task()
	TrackHistory(status string, message string, responseData string)
}

type WorkerPool struct {
	name      string
	tasksChan chan Worker
	wg        sync.WaitGroup
}

func NewWorkerPool(maxGoroutines int, buffered uint, name string) *WorkerPool {
	workerChain := make(chan Worker)
	if buffered != 0 {
		workerChain = make(chan Worker, buffered)
	}

	p := WorkerPool{
		name:      name,
		tasksChan: workerChain,
	}

	for i := 0; i < maxGoroutines; i++ {
		p.wg.Add(1)

		go func() {
			for w := range p.tasksChan {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Run submits work to the pool.
func (p *WorkerPool) Run(w Worker) {
	p.tasksChan <- w
}

// Shutdown waits for all the goroutines to shutdown.
func (p *WorkerPool) Shutdown() {
	close(p.tasksChan)
	p.wg.Wait()
}
