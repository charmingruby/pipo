package model

type RawSentiment struct {
	ID        int    `json:"id"`
	Comment   string `json:"comment"`
	Sentiment int    `json:"sentiment"`
}

func NewRawSentiment(id int, comment string, sentiment int) *RawSentiment {
	return &RawSentiment{
		ID:        id,
		Comment:   comment,
		Sentiment: sentiment,
	}
}
