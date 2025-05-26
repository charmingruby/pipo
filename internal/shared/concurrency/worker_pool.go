package concurrency

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrWorkerPoolStillRunning = errors.New("worker pool is still running")
	ErrWorkerPoolClosed       = errors.New("worker pool is closed")
)

// ProcessFunc defines a function type that processes a message of type T and returns a result of type R
type ProcessFunc[T any, R any] func(msg T) (R, error)

// WorkerPool implements a concurrent worker pool pattern for processing messages.
// It allows processing multiple messages concurrently using a specified number of workers.
// T represents the input message type and R represents the output result type.
//
// Example usage:
//
//	func main() {
//		processFunc := func(msg Input) (Output, error) {
//			return Output{ProcessedText: strings.ToUpper(msg.Text)}, nil
//		}
//
//		wp := NewWorkerPool(processFunc, 10)
//
//		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		defer cancel()
//
//		wp.Run(ctx)
//
//		amountOfMessages := 10000
//		messages := make([]Input, amountOfMessages)
//
//		for i := range amountOfMessages {
//			messages[i] = Input{
//				Text: fmt.Sprintf("message-%d", i),
//			}
//		}
//
//		var wg sync.WaitGroup
//		wg.Add(2)
//
//		go func() {
//			defer wg.Done()
//			for msg := range wp.Output() {
//				fmt.Println(msg.ProcessedText)
//			}
//		}()
//
//		go func() {
//			defer wg.Done()
//			for err := range wp.Error() {
//				fmt.Println(err)
//			}
//		}()
//
//		if err := wp.SendBatch(ctx, messages); err != nil {
//			fmt.Println(err)
//		}
//
//		if err := wp.Close(); err != nil {
//			fmt.Println(err)
//		}
//
//		wg.Wait()
//		fmt.Printf("Is closed: %t\n", wp.IsClosed())
//	}
type WorkerPool[T any, R any] struct {
	// Incomming messages channel
	inCh chan T
	// Outgoing messages channel
	outCh chan R
	// Errors channel
	errCh chan error
	// Done channel, used to signal when the worker pool is done
	doneCh chan struct{}
	// Mutex for thread safety
	mu sync.Mutex
	// WaitGroup for the workers
	wg sync.WaitGroup
	// Processing function
	processFunc ProcessFunc[T, R]
	// Concurrency level, represents the number of workers
	concurrency int
	// Closed flag, used to signal when the worker pool is closed
	closed bool
}

func NewWorkerPool[T any, R any](processFunc ProcessFunc[T, R], concurrency int) *WorkerPool[T, R] {
	return &WorkerPool[T, R]{
		inCh:        make(chan T, concurrency*2),
		outCh:       make(chan R, concurrency*2),
		errCh:       make(chan error, concurrency*2),
		doneCh:      make(chan struct{}),
		mu:          sync.Mutex{},
		wg:          sync.WaitGroup{},
		closed:      false,
		processFunc: processFunc,
		concurrency: concurrency,
	}
}

func (wp *WorkerPool[T, R]) Run(ctx context.Context) {
	for range wp.concurrency {
		wp.wg.Add(1)

		go func() {
			defer wp.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case <-wp.doneCh:
					for {
						select {
						case msg, ok := <-wp.inCh:
							if !ok {
								return
							}

							processedMsg, err := wp.processFunc(msg)
							if err != nil {
								wp.errCh <- err
								continue
							}

							wp.outCh <- processedMsg
						default:
							return
						}
					}
				case msg, ok := <-wp.inCh:
					if !ok {
						return
					}

					processedMsg, err := wp.processFunc(msg)
					if err != nil {
						wp.errCh <- err
						continue
					}

					wp.outCh <- processedMsg
				}
			}
		}()
	}
}

func (wp *WorkerPool[T, R]) SendBatch(ctx context.Context, messages []T) error {
	wp.mu.Lock()
	if wp.closed {
		wp.mu.Unlock()
		return ErrWorkerPoolClosed
	}
	wp.mu.Unlock()

	for _, msg := range messages {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-wp.doneCh:
			return ErrWorkerPoolClosed
		case wp.inCh <- msg:
		}
	}
	return nil
}

func (wp *WorkerPool[T, R]) Input() chan<- T {
	return wp.inCh
}

func (wp *WorkerPool[T, R]) Output() <-chan R {
	return wp.outCh
}

func (wp *WorkerPool[T, R]) Error() <-chan error {
	return wp.errCh
}

func (wp *WorkerPool[T, R]) Close() error {
	wp.mu.Lock()
	if wp.closed {
		wp.mu.Unlock()
		return ErrWorkerPoolClosed
	}

	wp.closed = true
	close(wp.doneCh)
	wp.mu.Unlock()

	wp.wg.Wait()

	close(wp.inCh)
	close(wp.outCh)
	close(wp.errCh)

	return nil
}

func (wp *WorkerPool[T, R]) IsClosed() bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.closed
}
