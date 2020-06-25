package main

import (
	"fmt"
	"math"
	"sort"
)

// PrintResponseTimesStats ...
func PrintResponseTimesStats(messages []*Message) {
	responseTimes := make([]int64, len(messages))
	for i, v := range messages {
		responseTimes[i] = v.GetResponseTime()
	}
	fmt.Println("Response Times Stats:")
	printGeneralStats(responseTimes)
	fmt.Println("==========================")
}

// PrintServiceTimesStats ...
func PrintServiceTimesStats(messages []*Message) {
	serviceTimes := make([]int64, len(messages))
	for i, v := range messages {
		serviceTimes[i] = v.GetTotalServiceTime()
	}
	fmt.Println("Service Times Stats:")
	printGeneralStats(serviceTimes)
	fmt.Println("==========================")
}

// PrintQueueTimesStats ...
func PrintQueueTimesStats(messages []*Message) {
	queueTimes := make([]int64, len(messages))
	for i, v := range messages {
		queueTimes[i] = v.GetTotalQueueTime()
	}
	fmt.Println("Queue Times Stats:")
	printGeneralStats(queueTimes)
	fmt.Println("==========================")
}

func printGeneralStats(values []int64) {
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })

	minValue := float32(values[0])
	maxValue := float32(values[len(values)-1])
	percentile1 := float32(values[int(0.01*float32(len(values)))])
	percentile5 := float32(values[int(0.05*float32(len(values)))])
	percentile10 := float32(values[int(0.10*float32(len(values)))])
	percentile90 := float32(values[int(0.90*float32(len(values)))])
	percentile95 := float32(values[int(0.95*float32(len(values)))])
	percentile99 := float32(values[int(0.99*float32(len(values)))])

	var temp float64
	for _, num := range values {
		temp += float64(num)
	}

	mean := temp / float64(len(values))
	temp = 0
	for _, num := range values {
		n := float64(num)
		temp += (n - mean) * (n - mean)
	}
	std := math.Sqrt(temp / float64(len(values)))

	s := fmt.Sprintf("min:\t%08.3f\nmax:\t%08.3f\navg:\t%08.3f\nstd:\t%08.3f\n1:\t%08.3f\n5:\t%08.3f\n10:\t%08.3f\n90:\t%08.3f\n95:\t%08.3f\n99:\t%08.3f\n", minValue, maxValue, mean, std, percentile1, percentile5, percentile10, percentile90, percentile95, percentile99)
	fmt.Println(s)
}

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
		panic(fmt.Errorf("100 must be dividable by %08.3f", 5))
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
func PrintPercentiles(messages []*Message) { //TODO: fix!
	fmt.Println("90", GetResponseTimePercentile(messages, 90))
	fmt.Println("95", GetResponseTimePercentile(messages, 95))
	fmt.Println("99", GetResponseTimePercentile(messages, 99))
}

// GetResponseTimePercentile ...
func GetResponseTimePercentile(messages []*Message, p int) int64 { //TODO: fix!
	sort.Slice(messages, func(i, j int) bool { return messages[i].GetResponseTime() < messages[j].GetResponseTime() })
	index := (len(messages) * p) / 100

	return messages[index].GetResponseTime()
}

// PrintAllQueueTimes ...
func PrintAllQueueTimes(messages []*Message) {
	sums := make([]int64, len(messages[0].Traces)-1)

	for _, message := range messages {
		times := message.GetEachQueueTime()
		for j := 0; j < len(message.Traces)-1; j++ {
			sums[j] += times[j]
		}
	}
	fmt.Println("AVERAGE QUEUE TIMES FOR EACH SERVICE")
	for i := 0; i < len(sums); i++ {
		fmt.Println(i+1, float64(sums[i])/float64(len(messages)))
	}
}

// PrintAllServiceTimes ...
func PrintAllServiceTimes(messages []*Message) {
	sums := make([]int64, len(messages[0].Traces))

	for _, message := range messages {
		times := message.GetEachServiceTime()
		for j := 0; j < len(message.Traces); j++ {
			sums[j] += times[j]
		}
	}
	fmt.Println("AVERAGE SERVICE TIMES FOR EACH SERVICE")
	for i := 0; i < len(sums); i++ {
		fmt.Println(i+1, float64(sums[i])/float64(len(messages)))
	}
}
