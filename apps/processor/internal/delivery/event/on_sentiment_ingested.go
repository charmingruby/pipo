package event

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/charmingruby/pipo-lib/concurrency"
	"github.com/charmingruby/pipo/apps/processor/internal/core/model"
	"github.com/charmingruby/pipo/apps/processor/internal/core/service"
)

const batchSize = 1000

func (h *Handler) onSentimentIngested() []error {
	wp := concurrency.NewWorkerPool(h.service.ProcessRawData, 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp.Run(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	errors := make([]error, 0)
	batch := make([]model.RawSentiment, 0, batchSize)

	go func() {
		defer wg.Done()
		for err := range wp.Error() {
			h.logger.Error("error processing sentiment ingested", "error", err)
		}
	}()

	go func() {
		defer wg.Done()
		for op := range wp.Output() {
			h.logger.Info("processed batch",
				"batch-size", len(batch),
				"success-count", op.SuccessCount,
			)
		}
	}()

	if err := h.broker.Subscribe(context.Background(), h.topics.SentimentIngested, func(message []byte) error {
		h.logger.Info("received sentiment ingested")

		var rawSentiment model.RawSentiment
		if err := json.Unmarshal(message, &rawSentiment); err != nil {
			return err
		}

		h.logger.Info("parsed sentiment ingested", "raw-sentiment-id", rawSentiment.ID)

		batch = append(batch, rawSentiment)

		h.logger.Debug("batch", "batch-size", len(batch))

		if len(batch) >= batchSize {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case wp.Input() <- service.ProcessRawDataInput{
				RawSentiments: batch,
			}:
				h.logger.Debug("sent batch new batch")

				batch = make([]model.RawSentiment, 0, batchSize)
				return nil
			}
		}

		return nil
	}); err != nil {
		errors = append(errors, err)
	}

	wg.Wait()

	if len(batch) > 0 {
		select {
		case <-ctx.Done():
			return append(errors, ctx.Err())
		case wp.Input() <- service.ProcessRawDataInput{
			RawSentiments: batch,
		}:
		}
	}

	return errors
}
