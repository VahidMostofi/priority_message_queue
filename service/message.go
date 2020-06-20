package main

import (
	"time"
)

// Trace is the history for each step
type Trace struct {
	ReceivedAt           int64 `json:"receivedAt"`
	ProcessingStartedAt  int64 `json:"processingStartedAt"`
	QueueDuration        int64 `json:"queueDuration"`
	FinishedProcessingAt int64 `json:"finishedProcessingAt"`
	PublishedAt          int64 `json:"publishedAt"`
}

// Message is the structure that we send into the queue
type Message struct {
	Data       int     `json:"data"`
	ID         string  `json:"id"`
	Traces     []Trace `json:"traces"`
	Priorities []uint8 `json:"priorities"`
	CreatedAt  int64   `json:"createdAt"`
}

// Received needs to be called when a message is received from queue
func (msg *Message) Received() {
	msg.Traces = append(msg.Traces, Trace{})
	msg.Traces[len(msg.Traces)-1].ReceivedAt = time.Now().UnixNano() / 1000000
	if len(msg.Traces) > 1 {
		msg.Traces[len(msg.Traces)-1].QueueDuration = msg.Traces[len(msg.Traces)-1].ReceivedAt - msg.Traces[len(msg.Traces)-2].PublishedAt
		// msg.Traces[len(msg.Traces)-1].QueueDuration /= 1000000
	}
}

// StartedProcessing needs to be called when processing on a message is started
func (msg *Message) StartedProcessing() {
	msg.Traces[len(msg.Traces)-1].ProcessingStartedAt = time.Now().UnixNano() / 1000000
}

// FinishedProcessing needs be called when the processing of a messaged is done
func (msg *Message) FinishedProcessing() {
	msg.Traces[len(msg.Traces)-1].FinishedProcessingAt = time.Now().UnixNano() / 1000000
}

// Published needs to be called when a message is sent to queue
func (msg *Message) Published() {
	msg.Traces[len(msg.Traces)-1].PublishedAt = time.Now().UnixNano() / 1000000
}

// GetResponseTime returns difference between the time a request is created and the time it is published to the last queue
func (msg *Message) GetResponseTime() int64 {
	return msg.Traces[len(msg.Traces)-1].PublishedAt - msg.CreatedAt
}

// GetTotalServiceTime ...
func (msg *Message) GetTotalServiceTime() int64 {

	var total int64
	for _, trace := range msg.Traces {
		total += (trace.FinishedProcessingAt - trace.ProcessingStartedAt)
	}
	return total
}

// GetTotalQueueTime ...
func (msg *Message) GetTotalQueueTime() int64 {
	var total int64
	if len(msg.Traces) <= 1 {
		return 0
	}

	for i := 1; i < len(msg.Traces); i++ {
		total += msg.Traces[i].ReceivedAt - msg.Traces[i-1].PublishedAt
	}
	return total
}
