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
			if len(msg.Traces) == 3 {
				if msg.Traces[len(msg.Traces)-1].QueueDuration > 40 {
					priority = 2
					fmt.Println("PRIORITY 2")
				} else {
					priority = 1
					fmt.Println("PRIORITY 1")
				}
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
