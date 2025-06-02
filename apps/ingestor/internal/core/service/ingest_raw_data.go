package service

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/charmingruby/pipo-lib/concurrency"
	"github.com/charmingruby/pipo/apps/ingestor/data"
	"github.com/charmingruby/pipo/apps/ingestor/internal/core/model"
	"github.com/charmingruby/pipo/apps/ingestor/pkg/csv"
)

// IngestRawDataInput is the input for the IngestRawData function.
type IngestRawDataInput struct {
	// Records is the number of records to be ingested.
	Records int
}

// IngestRawDataOutput is the output for the IngestRawData function.
type IngestRawDataOutput struct {
	// Errors is the errors that occurred during the ingestion.
	Errors []error
	// IngestedDataCount is the number of records that were ingested.
	IngestedDataCount int
}

// ingestRawDataProcessorInput is the input for the ingestRawDataProcessor function.
type ingestRawDataProcessorInput = []string

// ingestRawDataProcessorOutput is the output for the ingestRawDataProcessor function.
type ingestRawDataProcessorOutput struct {
	// Data is the data that was ingested.
	Data model.RawSentiment
}

// IngestRawData processes the ingestion of raw data.
//
// ctx is the context for the operation.
// in is the input for the operation.
//
// Returns the output for the operation and an error if the operation failed.
func (s *Service) IngestRawData(
	ctx context.Context,
	in IngestRawDataInput,
) (IngestRawDataOutput, error) {
	file, err := data.SentimentCSV.Open("sentiment_data.csv") // #nosec G304
	if err != nil {
		return IngestRawDataOutput{}, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.ErrorContext(ctx, "error closing file", "err", err)
		}
	}()

	records, err := csv.ParseFile(file, in.Records)
	if err != nil {
		return IngestRawDataOutput{}, err
	}

	ingestedData := make([]model.RawSentiment, 0)
	processingErrors := make([]error, 0)

	wp := concurrency.NewWorkerPool(
		func(ctx context.Context, record ingestRawDataProcessorInput) (ingestRawDataProcessorOutput, error) {
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

	if err := wp.SendBatch(ctx, records); err != nil {
		return IngestRawDataOutput{}, err
	}

	if err := wp.Close(); err != nil {
		return IngestRawDataOutput{}, err
	}

	wg.Wait()

	s.logger.InfoContext(
		ctx,
		"ingested data",
		"total-records", len(records),
		"success-count", len(ingestedData),
		"errors-count", len(processingErrors),
	)

	return IngestRawDataOutput{
		IngestedDataCount: len(ingestedData),
		Errors:            processingErrors,
	}, nil
}
