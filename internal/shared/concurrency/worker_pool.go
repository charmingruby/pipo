package concurrency

import (
	"errors"
	"sync"
)

var ErrWorkerPoolStillRunning = errors.New("worker pool is still running")

type Message struct {
	ID   string
	Data []byte
	Err  error
}

type Processor interface {
	Process(msg Message) Message
}

type WorkerPool struct {
	inCh        chan Message
	outCh       chan Message
	errCh       chan Message
	doneCh      chan struct{}
	proc        Processor
	wg          sync.WaitGroup
	concurrency int
}

func NewWorkerPool(proc Processor, concurrency int) *WorkerPool {
	return &WorkerPool{
		inCh:        make(chan Message),
		outCh:       make(chan Message),
		errCh:       make(chan Message),
		doneCh:      make(chan struct{}, 1),
		proc:        proc,
		concurrency: concurrency,
		wg:          sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Start() {
	for range wp.concurrency {
		wp.wg.Add(1)

		go func() {
			defer wp.wg.Done()

			for msg := range wp.inCh {
				processedMsg := wp.proc.Process(msg)

				if processedMsg.Err != nil {
					wp.errCh <- processedMsg
					continue
				}

				wp.outCh <- processedMsg
			}
		}()
	}
}

func (wp *WorkerPool) Stop() error {
	select {
	case <-wp.doneCh:
		return ErrWorkerPoolStillRunning
	default:
		close(wp.doneCh)
	}

	close(wp.inCh)

	wp.wg.Wait()

	close(wp.outCh)
	close(wp.errCh)

	return nil
}
