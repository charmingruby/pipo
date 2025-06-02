package repository

import (
	"context"

	"github.com/charmingruby/pipo/apps/processor/internal/core/model"
)

// SentimentRepository is the repository for the sentiment.
type SentimentRepository interface {
	// CreateMany creates many sentiments.
	//
	// ctx is the context.
	// sentiments is the sentiments to create.
	//
	// Returns the number of sentiments created and an error if one occurs.
	CreateMany(ctx context.Context, sentiments []model.Sentiment) (int64, error)
}
