package main

import (
	"fmt"
	"math/rand"
	"os"
)

const NoPriority = "no_priority"
const QueueTime = "queue_time"
const RandomPriority = "random"

func (msg *Message) SetPriority(strategy string) {
	var queueDuration int64 = 0
	var priority uint8 = 0

	if strategy == NoPriority {
		msg.Priorities = append(msg.Priorities, priority)
	} else if strategy == QueueTime {
		if len(msg.Traces) > 1 {
			for i := 1; i < len(msg.Traces); i++ {
				queueDuration += msg.Traces[i].QueueDuration
			}
			if queueDuration <= 160*1 {
				priority = 0
			} else if queueDuration <= 160*2 {
				priority = 1
			} else if queueDuration <= 160*3 {
				priority = 2
			} else if queueDuration <= 160*4 {
				priority = 3
			} else if queueDuration <= 160*5 {
				priority = 4
			} else if queueDuration <= 160*6 {
				priority = 5
			} else if queueDuration <= 160*7 {
				priority = 6
			} else if queueDuration <= 160*8 {
				priority = 7
			} else if queueDuration <= 160*9 {
				priority = 8
			} else {
				priority = 9
			}
		} else {
			priority = 0
		}
	} else if strategy == RandomPriority {
		priority = uint8(rand.Uint32() % 10)
	} else {
		panic("unknown strategy " + strategy)
	}
	if os.Getenv("DEBUG") == "TRUE" {
		fmt.Printf("%s: queue time: %10d, priority: %2d\n", os.Getenv("ROLE"), queueDuration, priority)
	}
	msg.Priorities = append(msg.Priorities, priority)
}
