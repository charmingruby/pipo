package service

import (
	"context"
	"errors"

	"github.com/charmingruby/pipo/internal/sentiment/core/model"
)

const (
	NEGATIVE_EMOTION = iota
	NEUTRAL_EMOTION
	POSITIVE_EMOTION
)

var (
	ErrInvalidSentiment         = errors.New("invalid sentiment")
	ErrBlankComment             = errors.New("comment should not be empty")
	ErrFailedToPersistSentiment = errors.New("failed to persist sentiment")
)

type ProcessRawDataInput struct {
	RawSentiment model.RawSentiment
}

type ProcessRawDataOutput struct {
	Sentiment model.Sentiment
}

func (s *Service) ProcessRawData(
	ctx context.Context,
	in ProcessRawDataInput,
) (ProcessRawDataOutput, error) {
	rawData := in.RawSentiment

	if err := s.validateRawData(rawData); err != nil {
		return ProcessRawDataOutput{}, err
	}

	sentiment, err := s.transformRawData(rawData)
	if err != nil {
		return ProcessRawDataOutput{}, err
	}

	// if err := s.persistSentiment(ctx, sentiment); err != nil {
	// 	return ProcessRawDataOutput{}, err
	// }

	return ProcessRawDataOutput{
		Sentiment: sentiment,
	}, nil
}

func (s *Service) validateRawData(rawData model.RawSentiment) error {
	if rawData.Sentiment < NEGATIVE_EMOTION || rawData.Sentiment > POSITIVE_EMOTION {
		return ErrInvalidSentiment
	}

	if rawData.Comment == "" {
		return ErrBlankComment
	}

	return nil
}

func (s *Service) transformRawData(rawData model.RawSentiment) (model.Sentiment, error) {
	emotion := s.mapEmotion(rawData.Sentiment)

	var excerpt string
	if rawData.Comment != "" && len(rawData.Comment) > 100 {
		minComment := rawData.Comment[:97]
		excerpt = minComment + "..."
	}

	sentiment := model.NewSentiment(model.SentimentInput{
		ID:      rawData.ID,
		Excerpt: excerpt,
		Comment: rawData.Comment,
		Emotion: emotion,
	})

	return *sentiment, nil
}

func (s *Service) persistSentiment(ctx context.Context, sentiment model.Sentiment) error {
	if err := s.sentimentRepository.Create(ctx, sentiment); err != nil {
		s.logger.Error("failed to persist sentiment", "error", err)

		return ErrFailedToPersistSentiment
	}

	return nil
}

func (s *Service) mapEmotion(sentiment int) string {
	switch sentiment {
	case NEGATIVE_EMOTION:
		return "negative"
	case NEUTRAL_EMOTION:
		return "neutral"
	case POSITIVE_EMOTION:
		return "positive"
	default:
		return "unknown"
	}
}
