package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

// ServiceProvider ...
type ServiceProvider struct {
	Status          string
	SourceQueueName string
	TargetQueueName string
	WorkWeight      int
	WG              sync.WaitGroup
	QA              *QueueAdapter
}

// NewServiceProvider ...
func NewServiceProvider(qa *QueueAdapter) *ServiceProvider {
	targetQueueName := os.Getenv("TARGET_QUEUE")
	sourceQueueName := os.Getenv("SOURCE_QUEUE")
	if len(targetQueueName) == 0 || len(sourceQueueName) == 0 {
		panic("TARGET_QUEUE or SOURCE_QUEUE is not provided")
	}
	workWeight, err := strconv.Atoi(os.Getenv("WORK_WEIGHT"))
	if err != nil {
		panic(fmt.Errorf("can't convert WORK_WEIGHT: '%s' to int", os.Getenv("WORK_WEIGHT")))
	}
	s := &ServiceProvider{
		Status:          "NOT_READY",
		SourceQueueName: sourceQueueName,
		TargetQueueName: targetQueueName,
		WorkWeight:      workWeight,
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
		message := &Message{}
		json.Unmarshal(b, message)
		message.Received()
		if os.Getenv("DEBUG") == "TRUE" {
			fmt.Println(os.Getenv("ROLE"), "Received", message.ID)
		}
		message.StartedProcessing()
		for i := 0; i < s.WorkWeight; i++ {
			isPrime(message.Data)
		}
		message.FinishedProcessing()
		message.SetPriority(NoPriority)
		message.Published()
		b, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		s.QA.Publish(b, message.Priorities[len(message.Priorities)-1], s.TargetQueueName)
		if os.Getenv("DEBUG") == "TRUE" {
			fmt.Println(os.Getenv("ROLE"), "Published", message.ID)
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
