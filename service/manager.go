package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/VahidMostofi/pmq/service/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Manager struct {
	QA                 *QueueAdapter
	LGConnection       *grpc.ClientConn
	LGControllerClient api.ControllerClient
}

// NewManager creates a new Manager and returns it
func NewManager(qa *QueueAdapter) *Manager {
	m := &Manager{
		QA: qa,
	}
	time.Sleep(time.Second * 5)
	loadGeneratorURL := os.Getenv("LOAD_GENERATOR_URL")

	APIPort := os.Getenv("API_PORT")
	if len(APIPort) == 0 {
		panic(fmt.Errorf("API_PORT is empty"))
	}
	lgConn, err := grpc.Dial(loadGeneratorURL+":"+APIPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	log.Printf("connected to: %s\n", loadGeneratorURL+":"+APIPort)

	m.LGConnection = lgConn
	c := api.NewControllerClient(m.LGConnection)
	m.LGControllerClient = c

	return m
}

// WaitForLoadGenerator ....
func (m *Manager) WaitForLoadGenerator() {

	for {
		response, err := m.LGControllerClient.GetStatus(context.Background(), &api.Empty{})
		if err != nil {
			log.Printf("Cannot call GetStatus of LoadGenerator\n")
			log.Println(err)
			continue
		}
		if response.Status == "READY" {
			log.Println("Load Generator is ready")
			break
		} else {
			log.Println("Load Generator is not ready its " + response.Status)
		}
		time.Sleep(time.Second * 2)
	}
}

// StartLoadGenerator ...
func (m *Manager) StartLoadGenerator() string {
	response, err := m.LGControllerClient.Start(context.Background(), &api.Empty{})
	if err != nil {
		panic(err)
	}
	return response.Status
}

// StopLoadGenerator ...
func (m *Manager) StopLoadGenerator() string {
	response, err := m.LGControllerClient.Stop(context.Background(), &api.Empty{})
	if err != nil {
		panic(err)
	}
	return response.Status
}

// Close closes the connections to server (load generator api)
func (m *Manager) Close() {
	m.LGConnection.Close()
}

// ComputeStats comptues stats based on the messages on the TargetQueue
func (m *Manager) ComputeStats() {
	inputChan := make(chan []byte)
	log.Println("Consuming", os.Getenv("FINAL_QUEUE"))
	m.QA.Consume(os.Getenv("FINAL_QUEUE"), inputChan)

	reqCount, err := strconv.Atoi(os.Getenv("REQ_COUNT"))
	if err != nil {
		panic(fmt.Errorf("can't convert REQ_COUNT: '%s' to int", os.Getenv("REQ_COUNT")))
	}
	// rate, err := strconv.Atoi(os.Getenv("RATE"))
	// if err != nil {
	// 	panic(fmt.Errorf("can't convert RATE: '%s' to int", os.Getenv("RATE")))
	// }

	messages := make([]*Message, reqCount)
	count := 0
	flag := false
	for b := range inputChan {
		msg := &Message{}
		json.Unmarshal(b, msg)
		messages[count] = msg
		count++

		if os.Getenv("DEBUG") == "TRUE" {
			fmt.Println(os.Getenv("ROLE"), "Received", msg.ID)
		}

		if count > int(0.95*float32(reqCount)) && !flag {
			log.Println("95% of messages are consumed")
			flag = true
		}
		if flag && os.Getenv("DEBUG") == "TRUE" {
			log.Println(count)
		}
		if count == reqCount {
			break
		}
	}

	computeMessageCreationRate(messages)

	// printAll(messages)
	// PrintPercentileHistogram(messages)
	// PrintResponseTimeDetails(messages)
	// PrintPercentiles(messages)

	PrintQueueTimesStats(messages)
	PrintServiceTimesStats(messages)
	PrintResponseTimesStats(messages)
}

func computeMessageCreationRate(messages []*Message) {
	min := messages[0].CreatedAt
	max := messages[0].CreatedAt
	for _, msg := range messages {
		if msg.CreatedAt > max {
			max = msg.CreatedAt
		}
		if msg.CreatedAt < min {
			min = msg.CreatedAt
		}
	}
	duration := (max - min) / 1e3
	rate := len(messages) / int(duration)
	log.Println("Request Creation Rate is", rate)
}

func printAll(messages []*Message) {
	for _, msg := range messages {
		fmt.Println(msg)
	}
}
