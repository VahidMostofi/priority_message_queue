package main

type Trace struct {
	ReceivedAt           int64 `json:"receivedAt"`
	ProcessingStartedAt  int64 `json:"processingStartedAt"`
	QueueDuration        int64 `json:"queueDuration"`
	FinishedProcessingAt int64 `json:"finishedProcessingAt"`
	PublishedAt          int64 `json:"publishedAt"`
}

type Message struct {
	Data       int     `json:"data"`
	ID         string  `json:"id"`
	Traces     []Trace `json:"traces"`
	Priorities []uint8 `json:"priorities"`
}
