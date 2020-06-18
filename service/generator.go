package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

// LoadGenerator type
type LoadGenerator struct {
	Rate      int
	Duration  int
	Done      chan bool
	WG        sync.WaitGroup
	Publisher IPublisher
	status    string
}

// NewLoadGenerator Construct new LoadGenerator
func NewLoadGenerator(rate, duration int, publisher IPublisher) *LoadGenerator {
	lg := &LoadGenerator{
		Rate:      rate,
		Duration:  duration,
		Done:      make(chan bool),
		WG:        sync.WaitGroup{},
		Publisher: publisher,
		status:    "NOT_READY",
	}
	return lg
}

func (l *LoadGenerator) NewRequest() {
	msg := l.makeRequest()

	msg.SetPriority(NoPriority)

	b, err := json.Marshal(msg)
	if err != nil {
		panic(fmt.Errorf("problem while marshaling a message: %w", err))
	}
	// fmt.Println(string(b))

	l.Publisher.Publish(b, msg.Priorities[len(msg.Priorities)-1], os.Getenv("TARGET_QUEUE"))
}

// StartGenerator prepares the load generator. rate is the number of requests per second and duration is in seconds (for controlable)
func (l *LoadGenerator) Start() string {
	time.AfterFunc(time.Second*time.Duration(l.Duration), func() {
		l.Stop()
	})

	l.WG.Add(1)
	go func() {
		for {
			select {
			case <-l.Done:
				l.WG.Done()
				return
			default:
				go l.NewRequest()
			}
			time.Sleep(time.Duration(1e9 / l.Rate))
		}
	}()

	l.WG.Wait()
	return "Done"
}

// StopGenerator stops generator (for controlable)
func (l *LoadGenerator) Stop() string {
	l.Done <- true
	return "Done"
}

// GetStatus returns the status of the controlable
func (l *LoadGenerator) GetStatus() string {
	return l.status
}

// make a new Request and return it as a message
func (l *LoadGenerator) makeRequest() *Message {
	t, e := uuid.NewRandom()
	if e != nil {
		panic(fmt.Errorf("can create new random uuid: %w", e))
	}
	m := &Message{
		Data:       7,
		ID:         t.String(),
		Traces:     make([]Trace, 0),
		Priorities: make([]uint8, 0),
		CreatedAt:  time.Now().UnixNano(),
	}

	return m
}
