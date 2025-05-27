package repository

import (
	"context"

	"github.com/charmingruby/pipo/internal/sentiment/core/model"
)

type SentimentRepository interface {
	CreateMany(ctx context.Context, sentiments []model.Sentiment) (int64, error)
}
