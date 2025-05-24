package sentiment

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/charmingruby/pipo/internal/shared/messaging"
	"github.com/charmingruby/pipo/pkg/csv"
	"github.com/charmingruby/pipo/pkg/logger"
)

type Service struct {
	logger *logger.Logger
	broker messaging.Broker
}

func NewService(logger *logger.Logger, broker messaging.Broker) *Service {
	return &Service{
		logger: logger,
		broker: broker,
	}
}

type DispatchRawSentimentDataInput struct {
	FilePath string
	Records  int
}

func (s *Service) DispatchRawSentimentData(
	ctx context.Context,
	in DispatchRawSentimentDataInput,
) error {
	records, err := csv.ReadFile(in.FilePath, in.Records)
	if err != nil {
		return err
	}

	rawSentimentData := make([]RawSentiment, len(records))
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

		rawSentimentData[idx] = RawSentiment{
			ID:        id,
			Comment:   record[1],
			Sentiment: sentiment,
		}

		message, err := json.Marshal(rawSentimentData[idx])
		if err != nil {
			failedData = append(failedData, err.Error())
			continue
		}

		if err := s.broker.Publish(ctx, "sentiment", message); err != nil {
			failedData = append(failedData, err.Error())
			continue
		}
	}

	s.logger.Info("sentiments published", "count", len(rawSentimentData))

	s.logger.Error("failed sentiment data parsing", "data", failedData)

	return nil
}
