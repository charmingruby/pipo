package model

import (
	"time"

	"github.com/charmingruby/pipo/lib/core"
)

type Sentiment struct {
	CreatedAt  time.Time `json:"created_at"`
	ID         string    `json:"id"`
	Excerpt    string    `json:"excerpt"`
	Comment    string    `json:"comment"`
	Emotion    string    `json:"emotion"`
	DocumentID int       `json:"document_id"`
}

type SentimentInput struct {
	Excerpt    string `json:"excerpt"`
	Comment    string `json:"comment"`
	Emotion    string `json:"emotion"`
	DocumentID int    `json:"document_id"`
}

func NewSentiment(in SentimentInput) *Sentiment {
	return &Sentiment{
		ID:         core.NewID(),
		DocumentID: in.DocumentID,
		Excerpt:    in.Excerpt,
		Comment:    in.Comment,
		Emotion:    in.Emotion,
		CreatedAt:  time.Now(),
	}
}
