package main

import (
	"log"

	"github.com/VahidMostofi/pmq/service/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// StartManager starts manager
func StartManager() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7789", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewControllerClient(conn)
	response, err := c.IsReady(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s %s", response.Done, response.Ready)
}
