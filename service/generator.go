package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/google/uuid"
)

// LoadGenerator type
type LoadGenerator struct {
	ArrivalRate    int
	ReqCount       int
	Done           chan bool
	WG             sync.WaitGroup
	TargetQueue    string
	Publisher      IPublisher
	status         string
	GeneratedCount int
}

// NewLoadGenerator Construct new LoadGenerator
func NewLoadGenerator(arrivalRate, reqCount int, publisher IPublisher) *LoadGenerator {
	lg := &LoadGenerator{
		ArrivalRate: arrivalRate,
		ReqCount:    reqCount,
		Done:        make(chan bool),
		WG:          sync.WaitGroup{},
		Publisher:   publisher,
		status:      "NOT_READY",
		TargetQueue: os.Getenv("TARGET_QUEUE"),
	}
	return lg
}

// NewRequest ...
func (l *LoadGenerator) NewRequest() {
	msg := l.makeRequest()
	msg.Received()

	msg.SetPriority(NoPriority)
	msg.Published()
	b, err := json.Marshal(msg)
	if err != nil {
		panic(fmt.Errorf("problem while marshaling a message: %w", err))
	}

	l.Publisher.Publish(b, msg.Priorities[len(msg.Priorities)-1], l.TargetQueue)
	l.GeneratedCount++

	if os.Getenv("DEBUG") == "TRUE" {
		fmt.Println(l.GeneratedCount, os.Getenv("ROLE"), "Published", msg.ID)
	}
}

// Start prepares the load generator. rate is the number of requests per second and duration is in seconds (for controlable)
func (l *LoadGenerator) Start() string {
	// if I stop this using time, it may not generat all messages!
	// time.AfterFunc(time.Second*time.Duration(l.Duration), func() {
	// 	l.Stop()
	// })
	var bar *pb.ProgressBar
	if os.Getenv("PROGRESS_BAR") == "TRUE" {
		bar = pb.StartNew(l.ReqCount)
	}

	l.WG.Add(1)
	go func() {
		for {
			select {
			case <-l.Done:
				log.Println("donning watgroup")
				l.WG.Done()
				return
			default:
				if l.GeneratedCount == l.ReqCount {
					go l.Stop()
					if bar != nil {
						bar.Finish()
					}
				} else {
					go l.NewRequest()
					if bar != nil {
						bar.Increment()
					}
				}
			}
			interval := int64(1000.0 * (rand.ExpFloat64() / float64(l.ArrivalRate)))
			// fmt.Println(interval)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()

	l.WG.Wait()
	log.Println("Start method is existing...")
	return "Done"
}

// Stop stops generator (for controlable)
func (l *LoadGenerator) Stop() string {
	log.Println("Stopping loadgenerator")
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
		Data:       6700417,
		ID:         t.String(),
		Traces:     make([]Trace, 0),
		Priorities: make([]uint8, 0),
		CreatedAt:  time.Now().UnixNano() / 1000000,
	}

	return m
}
