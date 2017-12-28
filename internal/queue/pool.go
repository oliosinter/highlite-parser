package queue

import (
	"sync"
)

// NewPool creates a new instance of pool.
func NewPool(workerCount uint) *Pool {
	jobs := make(chan IJob)

	workers := make([]*Worker, 0, workerCount)
	for ; workerCount > 0; workerCount-- {
		workers = append(workers, NewWorker(jobs))
	}

	return &Pool{
		jobs:    jobs,
		workers: workers,
		stopped: make(chan bool),
	}
}

// Pool creates workers and handles new job request.
type Pool struct {
	jobs     chan IJob
	workers  []*Worker
	stopOnce sync.Once
	stopped  chan bool
}

// AddJob adds a new job into worker queue.
func (p *Pool) AddJob(job IJob) {
	go func() {
		select {
		case p.jobs <- job:
		case <-p.stopped:
		}
	}()
}

// Stop gracefully stops all workers. It doesn't block.
// Returns a channel to notify when all workers have been stopped.
func (p *Pool) Stop() <-chan bool {
	p.stopOnce.Do(func() {
		go p.stopWorkersGracefully()
	})

	return p.stopped
}

// Stops all workers, waits until every worker has finished its work.
func (p *Pool) stopWorkersGracefully() {
	stops := make(chan bool)
	for _, w := range p.workers {
		go func(worker *Worker) {
			stops <- <-worker.Stop()
		}(w)
	}

	for i := 0; i < len(p.workers); i++ {
		<-stops
	}

	close(p.stopped)
}
