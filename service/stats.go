package main

import (
	"fmt"
	"sort"
)

// PrintResponseTimeDetails ...
func PrintResponseTimeDetails(messages []*Message) {
	for _, msg := range messages {
		PrintResponseTimeDetail(msg)
	}
}

// PrintResponseTimeDetail ...
func PrintResponseTimeDetail(message *Message) {
	fmt.Printf("queue time: %8d, service time: %8d, response time: %8d\n", message.GetTotalQueueTime(), message.GetTotalServiceTime(), message.GetResponseTime())
}

// PrintPercentileHistogram ...
func PrintPercentileHistogram(messages []*Message) {
	// counts, binMax := GetPercentileHistogram(messages)
	//TODO
}

// GetPercentileHistogram returns an array containing historgram for percentiles of values (not cumulative)
func GetPercentileHistogram(messages []*Message) ([]int, []int64) {
	const binSize = 5
	sort.Slice(messages, func(i, j int) bool { return messages[i].GetResponseTime() < messages[j].GetResponseTime() })
	if 100%binSize != 0 {
		panic(fmt.Errorf("100 must be dividable by %d", 5))
	}

	counts := make([]int, int(100/binSize))
	binMax := make([]int64, int(100/binSize))
	for i, msg := range messages {
		counts[int(i/binSize)]++
		binMax[int(i/binSize)] = msg.GetResponseTime()
	}

	return counts, binMax
}

// PrintPercentiles ,,,
func PrintPercentiles(messages []*Message){ //TODO: fix!
	fmt.Println("90",GetResponseTimePercentile(messages, 90))
	fmt.Println("95",GetResponseTimePercentile(messages, 95))
	fmt.Println("99",GetResponseTimePercentile(messages, 99))
}

// GetResponseTimePercentile ...
func GetResponseTimePercentile(messages []*Message, p int) int64{ //TODO: fix!
	sort.Slice(messages, func(i, j int) bool { return messages[i].GetResponseTime() < messages[j].GetResponseTime() })
	index := (len(messages) * p) / 100

	return messages[index].GetResponseTime()
}
