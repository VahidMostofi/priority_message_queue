package main

const NoPriority = "no_priority"
const QueueTime = "queue_time"

func (msg *Message) SetPriority(strategy string) {
	if strategy == NoPriority {
		msg.Priorities = append(msg.Priorities, 0)
	} else if strategy == QueueTime {

		var priority uint8 = 0
		if len(msg.Traces) > 1 {
			var queueDuration int64 = 0
			for i := 1; i < len(msg.Traces); i++ {
				queueDuration += msg.Traces[i].QueueDuration
			}
			if queueDuration <= 8 {
				priority = 0
			} else if queueDuration <= 16 {
				priority = 1
			} else if queueDuration <= 32 {
				priority = 2
			} else if queueDuration <= 64 {
				priority = 3
			} else if queueDuration <= 128 {
				priority = 4
			} else if queueDuration <= 256 {
				priority = 5
			} else if queueDuration <= 512 {
				priority = 6
			} else if queueDuration <= 1024 {
				priority = 7
			} else if queueDuration <= 2048 {
				priority = 8
			} else {
				priority = 9
			}
		} else {
			priority = 0
		}
		msg.Priorities = append(msg.Priorities, priority)
	} else {
		panic("unknown strategy " + strategy)
	}
}
