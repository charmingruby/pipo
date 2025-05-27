package event

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/charmingruby/pipo/internal/sentiment/core/model"
	"github.com/charmingruby/pipo/internal/sentiment/core/service"
	"github.com/charmingruby/pipo/internal/shared/concurrency"
)

func (h *Handler) onSentimentIngested() []error {
	wp := concurrency.NewWorkerPool(h.service.ProcessRawData, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp.Run(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	errors := make([]error, 0)

	go func() {
		defer wg.Done()
		for err := range wp.Error() {
			errors = append(errors, err)
		}
	}()

	go func() {
		defer wg.Done()
		for op := range wp.Output() {
			h.logger.Debug("processed message", "output", op)
		}
	}()

	if err := h.broker.Subscribe(context.Background(), h.topics.SentimentIngested, func(message []byte) error {
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
	}); err != nil {
		errors = append(errors, err)
	}

	wg.Wait()

	return errors
}
