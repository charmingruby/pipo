package model

import "time"

type Sentiment struct {
	ID        int       `json:"id"`
	Excerpt   string    `json:"excerpt"`
	Comment   string    `json:"comment"`
	Emotion   string    `json:"emotion"`
	CreatedAt time.Time `json:"created_at"`
}

type SentimentInput struct {
	ID      int    `json:"id"`
	Excerpt string `json:"excerpt"`
	Comment string `json:"comment"`
	Emotion string `json:"emotion"`
}

func NewSentiment(in SentimentInput) *Sentiment {
	return &Sentiment{
		ID:        in.ID,
		Excerpt:   in.Excerpt,
		Comment:   in.Comment,
		Emotion:   in.Emotion,
		CreatedAt: time.Now(),
	}
}
