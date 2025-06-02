package data

import (
	"embed"
)

// SentimentCSV is the embeddable sentiment data.
//
//go:embed sentiment_data.csv
var SentimentCSV embed.FS
