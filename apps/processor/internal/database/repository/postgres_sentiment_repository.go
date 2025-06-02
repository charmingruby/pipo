package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmingruby/pipo/apps/processor/internal/core/model"
	"github.com/jmoiron/sqlx"
)

// PostgresSentimentRepository is the repository for the sentiment.
type PostgresSentimentRepository struct {
	db *sqlx.DB
}

// NewPostgresSentimentRepository constructs a new PostgresSentimentRepository.
//
// db is the database connection.
//
// Returns a new PostgresSentimentRepository.
func NewPostgresSentimentRepository(db *sqlx.DB) *PostgresSentimentRepository {
	return &PostgresSentimentRepository{
		db: db,
	}
}

// CreateMany creates many sentiments.
//
// ctx is the context.
// sentiments is the sentiments to create.
//
// Returns the number of sentiments created and an error if one occurs.
func (r *PostgresSentimentRepository) CreateMany(ctx context.Context, sentiments []model.Sentiment) (int64, error) {
	if len(sentiments) == 0 {
		return 0, nil
	}

	query := `
		INSERT INTO sentiments (id, document_id, excerpt, comment, emotion, created_at) 
		VALUES %s`

	now := time.Now()
	values := make([]interface{}, 0, len(sentiments)*6)
	placeholders := make([]string, 0, len(sentiments))

	for i := range sentiments {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
	}

	for _, sentiment := range sentiments {
		values = append(values,
			sentiment.ID,
			sentiment.DocumentID,
			sentiment.Excerpt,
			sentiment.Comment,
			sentiment.Emotion,
			now,
		)
	}

	query = fmt.Sprintf(query, strings.Join(placeholders, ","))

	result, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return 0, fmt.Errorf("failed to batch insert sentiments: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return affected, nil
}
