package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"log"
	"net"

	"github.com/VahidMostofi/pmq/service/api"
	"google.golang.org/grpc"
)

var queueAdapter *QueueAdapter
var loadGenerator *LoadGenerator
var serviceProvider *ServiceProvider
var manager *Manager
var waitController sync.WaitGroup

func prepareGenerator() {
	reqCount, err := strconv.Atoi(os.Getenv("REQ_COUNT"))
	if err != nil {
		panic(fmt.Errorf("can't convert REQ_COUNT: '%s' to int", os.Getenv("REQ_COUNT")))
	}
	arrivalRate, err := strconv.Atoi(os.Getenv("ARRIVAL_RATE"))
	if err != nil {
		panic(fmt.Errorf("can't convert ARRIVAL_RATE: '%s' to int", os.Getenv("ARRIVAL_RATE")))
	}
	loadGenerator = NewLoadGenerator(arrivalRate, reqCount, queueAdapter)
	loadGenerator.status = "READY"
}

func prepareServiceProvider() {
	serviceProvider = NewServiceProvider(queueAdapter)
}

func prepareManager() {
	manager = NewManager(queueAdapter)
	manager.WaitForLoadGenerator()
	status := manager.StartLoadGenerator()
	log.Println("manager.StartLoadGenerator():", status)
	manager.ComputeStats()
	manager.Close()
}

func prepareController(c api.Controlable) {
	waitController.Add(1)
	go func() {
		APIPort, err := strconv.ParseInt(os.Getenv("API_PORT"), 10, 32)
		if err != nil {
			waitController.Done()
			panic(fmt.Errorf("can't parse API_PORT:%s to int", os.Getenv("API_PORT")))
		}

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", APIPort))
		if err != nil {
			waitController.Done()
			log.Fatalf("failed to listen: %v", err)
		} // create a server instance

		s := api.Server{C: c, WG: &waitController} // create a gRPC server object
		grpcServer := grpc.NewServer()

		// attach the Ping service to the server
		api.RegisterControllerServer(grpcServer, &s)

		log.Println("listening to port " + os.Getenv("API_PORT"))
		// start the server
		if err := grpcServer.Serve(lis); err != nil {
			waitController.Done()
			log.Fatalf("failed to serve: %s", err)
		}
	}()
}

// InitializeBasedOnRoles Starts the service based on the environment variable ROLE
func InitializeBasedOnRoles() {
	queueAdapter = NewQueueAdapter(os.Getenv("AMQP_URL"))

	if os.Getenv("ROLE") == "GENERATOR" {
		targetQueue := os.Getenv("TARGET_QUEUE")
		queueAdapter.CreateQueue(targetQueue)

		prepareGenerator()
		prepareController(loadGenerator)
	} else if os.Getenv("ROLE") == "SERVICE" {

		prepareServiceProvider()
		prepareController(serviceProvider)
	} else if os.Getenv("ROLE") == "MANAGER" {
		finalQueue := os.Getenv("Final_QUEUE")
		queueAdapter.CreateQueue(finalQueue)
		prepareManager()
	}

	waitController.Wait()
}
