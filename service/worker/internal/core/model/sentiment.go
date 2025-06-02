package model

import (
	"time"

	"github.com/charmingruby/pipo-lib/core"
)

// Sentiment is the model for the sentiment.
type Sentiment struct {
	// CreatedAt is the created at time of the sentiment.
	CreatedAt time.Time `json:"created_at"`
	// ID is the id of the sentiment.
	ID string `json:"id"`
	// Excerpt is the excerpt of the sentiment.
	Excerpt string `json:"excerpt"`
	// Comment is the comment of the sentiment.
	Comment string `json:"comment"`
	// Emotion is the emotion of the sentiment.
	Emotion string `json:"emotion"`
	// DocumentID is the id of the document.
	DocumentID int `json:"document_id"`
}

// SentimentInput is the input for the sentiment.
type SentimentInput struct {
	Excerpt    string `json:"excerpt"`
	Comment    string `json:"comment"`
	Emotion    string `json:"emotion"`
	DocumentID int    `json:"document_id"`
}

// NewSentiment constructs a new Sentiment.
//
// in is the input for the sentiment.
//
// Returns a new Sentiment.
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
