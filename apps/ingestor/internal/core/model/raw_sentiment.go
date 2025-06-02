package model

// RawSentiment is the model for a raw sentiment.
type RawSentiment struct {
	// Comment is the comment of the raw sentiment.
	Comment string `json:"comment"`
	// ID is the ID of the raw sentiment.
	ID int `json:"id"`
	// Sentiment is the sentiment of the raw sentiment.
	Sentiment int `json:"sentiment"`
}

// NewRawSentiment constructs a new RawSentiment.
//
// id is the ID of the raw sentiment.
// comment is the comment of the raw sentiment.
// sentiment is the sentiment of the raw sentiment.
//
// Returns a new RawSentiment.
func NewRawSentiment(id int, comment string, sentiment int) *RawSentiment {
	return &RawSentiment{
		ID:        id,
		Comment:   comment,
		Sentiment: sentiment,
	}
}
