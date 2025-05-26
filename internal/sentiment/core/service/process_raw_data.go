package service

import (
	"context"
	"fmt"

	"github.com/charmingruby/pipo/internal/sentiment/core/model"
)

type ProcessRawDataInput struct {
	RawSentiment model.RawSentiment
}

type ProcessRawDataOutput struct {
}

func (s *Service) ProcessRawData(
	ctx context.Context,
	in ProcessRawDataInput,
) (ProcessRawDataOutput, error) {
	fmt.Println(in.RawSentiment)

	return ProcessRawDataOutput{}, nil
}
