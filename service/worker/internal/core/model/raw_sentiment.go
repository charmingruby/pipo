package model

type RawSentiment struct {
	Comment   string `json:"comment"`
	ID        int    `json:"id"`
	Sentiment int    `json:"sentiment"`
}

func NewRawSentiment(id int, comment string, sentiment int) *RawSentiment {
	return &RawSentiment{
		ID:        id,
		Comment:   comment,
		Sentiment: sentiment,
	}
}
