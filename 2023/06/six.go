package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func main() {
	lines := strings.Split(input, "\n")
	times := parseLine(lines[0])
	distances := parseLine(lines[1])
	fmt.Println("Times: ", times)
	fmt.Println("Distances: ", distances)

	numberOfWays := computeNumberOfWays(times, distances)
	fmt.Println("Number of ways to beat the record (part 1): ", numberOfWays)

	numberOfWaysPart2 := computeNumberOfWays([]int64{58819676}, []int64{434104122191218})
	fmt.Println("Number of ways to beat the record (part 2): ", numberOfWaysPart2)
}

func parseLine(line string) []int64 {
	fields := strings.Fields(line)[1:]
	numbers := make([]int64, len(fields))
	for i, field := range fields {
		numbers[i], _ = strconv.ParseInt(field, 10, 64)
	}
	return numbers
}

func computeNumberOfWays(times []int64, distances []int64) interface{} {
	product := int64(1)
	for i := 0; i < len(times); i++ {
		product *= countSolutions(times[i], distances[i])
	}
	return product
}

func countSolutions(raceTimeInMs, minimumDistanceInMs int64) int64 {
	solutionCount := int64(0)
	for holdDurationInMs := int64(1); holdDurationInMs < raceTimeInMs; holdDurationInMs++ {
		distanceCoveredInMm := (raceTimeInMs - holdDurationInMs) * holdDurationInMs
		if distanceCoveredInMm > minimumDistanceInMs {
			solutionCount++
		}
	}
	return solutionCount
}
