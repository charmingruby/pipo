package event

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/charmingruby/pipo/internal/sentiment/core/model"
	"github.com/charmingruby/pipo/internal/sentiment/core/service"
	"github.com/charmingruby/pipo/internal/shared/concurrency"
)

func (h *Handler) onSentimentIngested() error {
	wp := concurrency.NewWorkerPool(h.service.ProcessRawData, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp.Run(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for err := range wp.Error() {
			fmt.Printf("Error processing message: %v\n", err)
		}
	}()

	go func() {
		defer wg.Done()
		for range wp.Output() {
		}
	}()

	return h.broker.Subscribe(context.Background(), h.topics.SentimentIngested, func(message []byte) error {
		var rawSentiment model.RawSentiment
		if err := json.Unmarshal(message, &rawSentiment); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case wp.Input() <- service.ProcessRawDataInput{
			RawSentiment: rawSentiment,
		}:
			return nil
		}
	})
}
