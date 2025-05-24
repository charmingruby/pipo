package sentiment

import (
	"context"
	"fmt"

	"github.com/charmingruby/pipo/pkg/csv"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
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

	for _, record := range records {
		fmt.Println(record)
	}

	return nil
}
