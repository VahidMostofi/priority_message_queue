package main

import (
	"encoding/json"
	"fmt"
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
}

// NewLoadGenerator Construct new LoadGenerator
func NewLoadGenerator(rate, duration int, publisher IPublisher) *LoadGenerator {
	lg := &LoadGenerator{
		Rate:      rate,
		Duration:  duration,
		Done:      make(chan bool),
		WG:        sync.WaitGroup{},
		Publisher: publisher,
	}
	return lg
}

func (l *LoadGenerator) NewRequest() {
	msg := l.makeRequest()

	b, err := json.Marshal(msg)
	if err != nil {
		panic(fmt.Errorf("problem while marshaling a message: %w", err))
	}
	fmt.Println(string(b))

	// l.Publisher.Publish()
}

// StartGenerator prepares the load generator. rate is the number of requests per second and duration is in seconds
func (l *LoadGenerator) StartGenerator() {
	time.AfterFunc(time.Second*time.Duration(l.Duration), func() {
		l.StopGenerator()
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
}

// StopGenerator stops generator
func (l *LoadGenerator) StopGenerator() {
	l.Done <- true
	l.WG.Done()

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

// func main() {
// 	const rate = 5     // rate per second
// 	const duration = 1 //seconds
// 	requests = make([]int64, 0)

// 	done := make(chan bool)

// 	wg := sync.WaitGroup{}
// 	wg.Add(1)
// 	go func() {
// 		for {
// 			select {
// 			case <-done:
// 				wg.Done()
// 				return
// 			default:
// 				go makeRequest()
// 			}
// 			time.Sleep(time.Duration(1e9 / rate))
// 		}
// 	}()

// 	wg.Wait()
// 	diffs := make([]int64, len(requests)-1)
// 	for i := 1; i < len(requests); i++ {
// 		diffs[i-1] = (requests[i] - requests[i-1]) / 1000000
// 	}
// 	fmt.Println(mean(diffs), len(requests))
// }
