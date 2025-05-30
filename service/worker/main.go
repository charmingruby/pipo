package main

import (
	"fmt"
	"unsafe"
)

type RawSentiment struct {
	Comment   string `json:"comment"`
	ID        int    `json:"id"`
	Sentiment int    `json:"sentiment"`
}

type Sentiment struct {
	ID        int    `json:"id"`
	Comment   string `json:"comment"`
	Sentiment int    `json:"sentiment"`
}

// OptimizedSentiment has fields ordered from largest to smallest
type OptimizedSentiment struct {
	Comment   string `json:"comment"`   // 16 bytes
	ID        int    `json:"id"`        // 8 bytes
	Sentiment int    `json:"sentiment"` // 8 bytes
}

func main() {
	raw := RawSentiment{}
	fmt.Printf("RawSentiment size: %d bytes\n", unsafe.Sizeof(raw))

	sentiment := Sentiment{}
	fmt.Printf("Sentiment size: %d bytes\n", unsafe.Sizeof(sentiment))

	optimized := OptimizedSentiment{}
	fmt.Printf("OptimizedSentiment size: %d bytes\n", unsafe.Sizeof(optimized))
}
