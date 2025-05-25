package concurrency

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_WorkerPool(t *testing.T) {
	t.Run("should process messages in concurrently", func(t *testing.T) {
		wp := NewWorkerPool(&dummyProcessor{}, 10)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		wp.Run(ctx)

		amountOfMessages := 10000
		messages := make([]dummyInput, amountOfMessages)

		for i := range amountOfMessages {
			messages[i] = dummyInput{
				ID:      i,
				RawText: fmt.Sprintf("message-%d", i),
			}
		}

		var wg sync.WaitGroup
		wg.Add(2)

		processedCount := 0
		var mu sync.Mutex

		go func() {
			defer wg.Done()
			for msg := range wp.Output() {
				mu.Lock()
				processedCount++
				mu.Unlock()

				assert.Equal(t, msg.Status, "processed")
			}
		}()

		go func() {
			defer wg.Done()
			for err := range wp.Error() {
				assert.NoError(t, err)
			}
		}()

		err := wp.SendBatch(ctx, messages)
		assert.NoError(t, err)

		err = wp.Close()
		assert.NoError(t, err)

		wg.Wait()

		assert.True(t, wp.IsClosed())
		assert.Equal(t, amountOfMessages, processedCount, "all messages should be processed")
	})
}

type dummyProcessor struct{}

type dummyInput struct {
	ID      int
	RawText string
}

type dummyOutput struct {
	ID     int
	Text   string
	Status string
}

func (p *dummyProcessor) Process(msg dummyInput) (dummyOutput, error) {
	return dummyOutput{
		ID:     msg.ID,
		Text:   msg.RawText,
		Status: "processed",
	}, nil
}
