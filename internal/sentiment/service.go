package sentiment

import (
	"context"
	"strconv"

	"github.com/charmingruby/pipo/pkg/csv"
	"github.com/charmingruby/pipo/pkg/logger"
)

type Service struct {
	logger *logger.Logger
}

func NewService(logger *logger.Logger) *Service {
	return &Service{
		logger: logger,
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
			Comment:   record[2],
			Sentiment: sentiment,
		}
	}

	s.logger.Info("failed sentiment data parsing", "data", failedData)

	return nil
}
