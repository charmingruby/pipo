package service

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/charmingruby/pipo/service/worker/internal/core/model"
)

const (
	DefaultComment = "Invalid UTF-8 text"

	NegativeEmotion = iota
	NeutralEmotion
	PositiveEmotion
)

var (
	ErrInvalidSentiment         = errors.New("invalid sentiment")
	ErrBlankComment             = errors.New("comment should not be empty")
	ErrFailedToPersistSentiment = errors.New("failed to persist sentiment")
)

type ProcessRawDataInput struct {
	RawSentiments []model.RawSentiment
}

type ProcessRawDataOutput struct {
	Sentiments   []model.Sentiment
	SuccessCount int64
}

func (s *Service) ProcessRawData(
	ctx context.Context,
	in ProcessRawDataInput,
) (ProcessRawDataOutput, error) {
	sentiments := make([]model.Sentiment, 0, len(in.RawSentiments))

	for _, rawData := range in.RawSentiments {
		if err := s.validateRawData(rawData); err != nil {
			return ProcessRawDataOutput{}, err
		}

		sentiment := s.transformRawData(rawData)

		sentiments = append(sentiments, sentiment)
	}

	affected, err := s.persistSentiments(ctx, sentiments)
	if err != nil {
		return ProcessRawDataOutput{}, err
	}

	return ProcessRawDataOutput{
		Sentiments:   sentiments,
		SuccessCount: affected,
	}, nil
}

func (s *Service) validateRawData(rawData model.RawSentiment) error {
	if rawData.Sentiment < NegativeEmotion || rawData.Sentiment > PositiveEmotion {
		return ErrInvalidSentiment
	}

	if rawData.Comment == "" {
		return ErrBlankComment
	}

	return nil
}

func (s *Service) transformRawData(rawData model.RawSentiment) model.Sentiment {
	emotion := s.mapEmotion(rawData.Sentiment)

	comment := rawData.Comment

	var excerpt string
	if len(comment) > 100 {
		minComment := comment[:97]
		excerpt = minComment + "..."
	} else {
		excerpt = comment
	}

	if !utf8.ValidString(comment) || !utf8.ValidString(excerpt) {
		comment = DefaultComment
		excerpt = DefaultComment
	}

	sentiment := model.NewSentiment(model.SentimentInput{
		DocumentID: rawData.ID,
		Excerpt:    excerpt,
		Comment:    comment,
		Emotion:    emotion,
	})

	return *sentiment
}

func (s *Service) persistSentiments(ctx context.Context, sentiments []model.Sentiment) (int64, error) {
	affected, err := s.sentimentRepository.CreateMany(ctx, sentiments)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to persist sentiments", "error", err)
		return 0, ErrFailedToPersistSentiment
	}

	return affected, nil
}

func (s *Service) mapEmotion(sentiment int) string {
	switch sentiment {
	case NegativeEmotion:
		return "negative"
	case NeutralEmotion:
		return "neutral"
	case PositiveEmotion:
		return "positive"
	default:
		return "unknown"
	}
}
