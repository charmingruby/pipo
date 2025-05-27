package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/charmingruby/pipo/lib/concurrency"
	"github.com/charmingruby/pipo/service/worker/internal/core/model"
	"github.com/charmingruby/pipo/service/worker/pkg/csv"
)

type IngestRawDataInput struct {
	FilePath string
	Records  int
}

type IngestRawDataOutput struct {
	ProcessedData []model.RawSentiment
	Errors        []error
}

func (s *Service) IngestRawData(
	ctx context.Context,
	in IngestRawDataInput,
) (IngestRawDataOutput, error) {
	records, err := csv.ReadFile(in.FilePath, in.Records)
	if err != nil {
		return IngestRawDataOutput{}, err
	}

	ingestedData := make([]model.RawSentiment, 0)
	processingErrors := make([]error, 0)

	wp := concurrency.NewWorkerPool(func(ctx context.Context, record ingestRawDataProcessorInput) (ingestRawDataProcessorOutput, error) {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return ingestRawDataProcessorOutput{}, err
		}

		sentiment, err := strconv.Atoi(record[2])
		if err != nil {
			return ingestRawDataProcessorOutput{}, err
		}

		rawSentimentData := model.NewRawSentiment(id, record[1], sentiment)

		message, err := json.Marshal(rawSentimentData)
		if err != nil {
			return ingestRawDataProcessorOutput{}, err
		}

		publishCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := s.broker.Publish(publishCtx, s.sentimentIngestTopic, message); err != nil {
			return ingestRawDataProcessorOutput{}, err
		}
		return ingestRawDataProcessorOutput{
			Data: *rawSentimentData,
		}, nil
	},
		10,
	)

	wp.Run(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for msg := range wp.Output() {
			ingestedData = append(ingestedData, msg.Data)
		}
	}()

	go func() {
		defer wg.Done()
		for err := range wp.Error() {
			processingErrors = append(processingErrors, err)
		}
	}()

	if err := wp.SendBatch(ctx, []ingestRawDataProcessorInput(records)); err != nil {
		return IngestRawDataOutput{}, err
	}

	if err := wp.Close(); err != nil {
		return IngestRawDataOutput{}, err
	}

	wg.Wait()

	s.logger.Info(
		"ingested data",
		"total-records", len(records),
		"success-count", len(ingestedData),
		"errors-count", len(processingErrors),
	)

	return IngestRawDataOutput{
		ProcessedData: ingestedData,
		Errors:        processingErrors,
	}, nil
}

type ingestRawDataProcessorInput = []string

type ingestRawDataProcessorOutput struct {
	Data model.RawSentiment
}
