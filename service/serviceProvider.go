package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// ServiceProvider ...
type ServiceProvider struct {
	Status          string
	SourceQueueName string
	TargetQueueName string
	ServiceRate     int
	WG              sync.WaitGroup
	QA              *QueueAdapter
	ReceivedCount   int
	PublishedCount  int
}

// NewServiceProvider ...
func NewServiceProvider(qa *QueueAdapter) *ServiceProvider {
	targetQueueName := os.Getenv("TARGET_QUEUE")
	sourceQueueName := os.Getenv("SOURCE_QUEUE")
	if len(targetQueueName) == 0 || len(sourceQueueName) == 0 {
		panic("TARGET_QUEUE or SOURCE_QUEUE is not provided")
	}
	serviceRate, err := strconv.Atoi(os.Getenv("SERVICE_RATE"))
	if err != nil {
		panic(fmt.Errorf("can't convert WORK_WEIGHT: '%s' to int", os.Getenv("WORK_WEIGHT")))
	}
	s := &ServiceProvider{
		Status:          "NOT_READY",
		SourceQueueName: sourceQueueName,
		TargetQueueName: targetQueueName,
		ServiceRate:     serviceRate,
		WG:              sync.WaitGroup{},
		QA:              qa,
	}
	s.QA.CreateQueue(s.SourceQueueName)
	s.QA.CreateQueue(s.TargetQueueName)
	s.WG.Add(1)
	go s.startHandling()
	return s
}

func (s *ServiceProvider) startHandling() {

	inputChan := make(chan []byte)
	s.QA.Consume(s.SourceQueueName, inputChan)
	s.Status = "READY"
	for b := range inputChan {
		s.ReceivedCount++
		message := &Message{}
		json.Unmarshal(b, message)
		message.Received()
		if os.Getenv("DEBUG") == "TRUE" {
			fmt.Println(os.Getenv("ROLE"), "Received", message.ID)
		}
		message.StartedProcessing()
		duration := int64(1000.0 * (rand.ExpFloat64() / float64(s.ServiceRate)))
		// start := time.Now().Nanosecond()
		end := time.Now().Add(time.Millisecond * time.Duration(duration))
		for time.Now().Before(end) {
		}
		// fmt.Println(duration, time.Now().Nanosecond()-start)
		message.FinishedProcessing()
		message.SetPriority(os.Getenv("PRIORITY_STRATEGY"))
		message.Published()
		b, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		s.QA.Publish(b, message.Priorities[len(message.Priorities)-1], s.TargetQueueName)
		s.PublishedCount++

		if os.Getenv("DEBUG") == "TRUE" {
			fmt.Println(s.ReceivedCount, s.PublishedCount, os.Getenv("ROLE"), "Published", message.ID)
		}
	}
	s.WG.Wait()
}

// Start ...
func (s *ServiceProvider) Start() string {
	return "Done"
}

// Stop ...
func (s *ServiceProvider) Stop() string {
	return ""
}

// GetStatus ...
func (s *ServiceProvider) GetStatus() string {
	return ""
}

func isPrime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}
