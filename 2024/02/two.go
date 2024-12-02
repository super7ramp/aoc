package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func main() {
	reports := strings.Split(input, "\n")
	reports = slices.DeleteFunc(reports, isReportUnsafe)
	fmt.Printf("%v safe report(s): %q", len(reports), reports)
}

func isReportUnsafe(report string) bool {
	return !isReportSafe(report)
}

func isReportSafe(report string) bool {
	levels := strings.Split(report, " ")
	firstLevel, _ := strconv.Atoi(levels[0])
	secondLevel, _ := strconv.Atoi(levels[1])
	if !isDiffSafe(firstLevel, secondLevel) {
		return false
	}

	increasing := secondLevel > firstLevel
	previousLevel := secondLevel
	for i := 2; i < len(levels); i++ {
		currentLevel, _ := strconv.Atoi(levels[i])
		safeDiff := isDiffSafe(previousLevel, currentLevel)
		keepsIncreasing := increasing && currentLevel > previousLevel
		keepsDecreasing := !increasing && currentLevel < previousLevel
		if !safeDiff || !(keepsIncreasing || keepsDecreasing) {
			return false
		}
		previousLevel = currentLevel
	}
	return true
}

func isDiffSafe(level1, level2 int) bool {
	diff := abs(level2 - level1)
	return diff > 0 && diff <= 3
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
