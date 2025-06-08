package service

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/charmingruby/pipo/apps/processor/internal/core/model"
)

const (
	// DefaultComment is the default comment, used when the comment is invalid UTF-8.
	DefaultComment = "Invalid UTF-8 text"

	// NegativeEmotion is the negative emotion, represented by 0.
	NegativeEmotion = 0
	// NeutralEmotion is the neutral emotion, represented by 1.
	NeutralEmotion = 1
	// PositiveEmotion is the positive emotion, represented by 2.
	PositiveEmotion = 2
)

var (
	// ErrInvalidSentiment is the error for an invalid sentiment.
	ErrInvalidSentiment = errors.New("invalid sentiment")
	// ErrBlankComment is the error for a blank comment.
	ErrBlankComment = errors.New("comment should not be empty")
	// ErrFailedToPersistSentiment is the error for a failed to persist sentiment.
	ErrFailedToPersistSentiment = errors.New("failed to persist sentiment")
)

// ProcessRawDataInput is the input for the process raw data.
type ProcessRawDataInput struct {
	// RawSentiments is the raw sentiments to process.
	RawSentiments []model.RawSentiment
}

// ProcessRawDataOutput is the output for the process raw data.
type ProcessRawDataOutput struct {
	// Sentiments is the sentiments created.
	Sentiments []model.Sentiment
	// SuccessCount is the number of sentiments created.
	SuccessCount int64
}

// ProcessRawData processes the raw data.
//
// ctx is the context.
// in is the input for the process raw data.
//
// Returns the output for the process raw data and an error if one occurs.
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
