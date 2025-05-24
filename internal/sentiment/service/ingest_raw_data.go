package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/charmingruby/pipo/internal/sentiment/model"
	"github.com/charmingruby/pipo/pkg/csv"
)

type IngestRawDataInput struct {
	FilePath string
	Records  int
	Topic    string
}

func (s *Service) IngestRawData(
	ctx context.Context,
	in IngestRawDataInput,
) error {
	records, err := csv.ReadFile(in.FilePath, in.Records)
	if err != nil {
		return err
	}

	rawSentimentData := make([]model.RawSentiment, len(records))
	failedData := make([]string, 0)

	for idx, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			failedData = append(failedData, record[0])
			continue
		}

		sentiment, err := strconv.Atoi(record[2])
		if err != nil {
			failedData = append(failedData, record[2])
			continue
		}

		rawSentimentData[idx] = *model.NewRawSentiment(id, record[1], sentiment)

		message, err := json.Marshal(rawSentimentData[idx])
		if err != nil {
			failedData = append(failedData, err.Error())
			continue
		}

		if err := s.broker.Publish(ctx, in.Topic, message); err != nil {
			failedData = append(failedData, err.Error())
			continue
		}
	}

	s.logger.Info("sentiments published", "count", len(rawSentimentData))

	s.logger.Error("failed sentiment data parsing", "data", failedData)

	return nil
}
