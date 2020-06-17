package main

import (
	"os"
	"strings"
)

var queueAdapter QueueAdapter

func prepareGenerator(){
	
}

func prepareService(){

}

func StartService(){
	targetQueue:= os.Getenv("TARGET_QUEUE")
	
	queueAdapter = NewQueueAdapter(os.Getenv("AMQP_URL"))
	queueAdapter.CreateQueue(targetQueue)

	if os.Getenv("ROLE") == "GENERATOR"{
		prepareGenerator()
	}else if os.Getenv("ROLE") == "SERVICE"{
		prepareService()
	}

	


}