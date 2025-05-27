package model

import (
	"time"

	"github.com/charmingruby/pipo/internal/shared/id"
)

type Sentiment struct {
	ID         string    `json:"id"`
	DocumentID int       `json:"document_id"`
	Excerpt    string    `json:"excerpt"`
	Comment    string    `json:"comment"`
	Emotion    string    `json:"emotion"`
	CreatedAt  time.Time `json:"created_at"`
}

type SentimentInput struct {
	DocumentID int    `json:"document_id"`
	Excerpt    string `json:"excerpt"`
	Comment    string `json:"comment"`
	Emotion    string `json:"emotion"`
}

func NewSentiment(in SentimentInput) *Sentiment {
	return &Sentiment{
		ID:         id.New(),
		DocumentID: in.DocumentID,
		Excerpt:    in.Excerpt,
		Comment:    in.Comment,
		Emotion:    in.Emotion,
		CreatedAt:  time.Now(),
	}
}
