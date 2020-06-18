package main

func main() {
	InitializeBasedOnRoles()
	// const amqpURL = "amqp://guest:guest@localhost:5672/"
	// const targetQueue = "targetQueue"
	// const sourceQueue = "sourceQueue"
	// queueAdapter := NewQueueAdapter(amqpURL)
	// queueAdapter.CreateQueue(targetQueue)

	// test1 := []byte("vahid")
	// test2 := []byte("saeed")
	// test3 := []byte("arman")
	// queueAdapter.Publish(test1, 1, targetQueue)
	// queueAdapter.Publish(test2, 2, targetQueue)
	// queueAdapter.Publish(test3, 3, targetQueue)
	// time.Sleep(15 * time.Second)
	// handler := make(chan []byte)
	// handledData := 0
	// go queueAdapter.Consume(targetQueue, handler)
	// for {
	// 	select {
	// 	case data := <-handler:
	// 		fmt.Println(string(data))
	// 		handledData += 1
	// 		if handledData == 3 {
	// 			queueAdapter.Close()
	// 			return
	// 		}
	// 	}
	// }

}
