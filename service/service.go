package main

import (
	"fmt"
	"os"
	"strconv"

	"log"
	"net"

	"github.com/VahidMostofi/pmq/service/api"
	"google.golang.org/grpc"
)

var queueAdapter *QueueAdapter
var loadGenerator *LoadGenerator

func prepareGenerator() {
	duration, err := strconv.Atoi(os.Getenv("DURATION"))
	if err != nil {
		panic(fmt.Errorf("can't convert DURATION: '%s' to int", os.Getenv("DURATION")))
	}
	rate, err := strconv.Atoi(os.Getenv("RATE"))
	if err != nil {
		panic(fmt.Errorf("can't convert RATE: '%s' to int", os.Getenv("RATE")))
	}
	loadGenerator = NewLoadGenerator(rate, duration, queueAdapter)
}

func prepareService() {

}

func prepareManager() {

}

func prepareController() {
	APIPort, err := strconv.ParseInt(os.Getenv("API_PORT"), 10, 32)
	if err != nil {
		panic(fmt.Errorf("can't parse API_PORT:%s to int", os.Getenv("API_PORT")))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", APIPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} // create a server instance

	s := api.Server{} // create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Ping service to the server
	api.RegisterControllerServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// StartService Start the service based on the environment variable ROLE
func StartService() {
	fmt.Println(os.Getenv("AMQP_URL"))
	targetQueue := os.Getenv("TARGET_QUEUE")

	queueAdapter = NewQueueAdapter(os.Getenv("AMQP_URL"))
	queueAdapter.CreateQueue(targetQueue)

	if os.Getenv("ROLE") == "GENERATOR" {
		prepareGenerator()
		prepareController()
	} else if os.Getenv("ROLE") == "SERVICE" {
		prepareService()
		prepareController()
	} else if os.Getenv("ROLE") == "MANAGER" {
		StartManager()
	}

}
