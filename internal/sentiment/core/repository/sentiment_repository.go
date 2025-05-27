package repository

import (
	"context"

	"github.com/charmingruby/pipo/internal/sentiment/core/model"
)

type SentimentRepository interface {
	Create(ctx context.Context, sentiment model.Sentiment) error
}
