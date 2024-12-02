package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

//go:embed input-example.txt
var input string

func main() {
	reports := strings.Split(input, "\n")

	part1Reports := slices.Clone(reports)
	part1Reports = slices.DeleteFunc(part1Reports, isReportUnsafeWithoutTolerance)
	fmt.Printf("Part 1: %v safe report(s): %q\n", len(part1Reports), part1Reports)

	part2Reports := slices.Clone(reports)
	part2Reports = slices.DeleteFunc(part2Reports, isReportUnsafeWithToleranceOfOne)
	fmt.Printf("Part 2: %v safe report(s): %q\n", len(part2Reports), part2Reports)

}

func isReportUnsafeWithoutTolerance(report string) bool {
	return !isReportSafe(report, 0)
}

func isReportUnsafeWithToleranceOfOne(report string) bool {
	return !isReportSafe(report, 1)
}

func isReportSafe(report string, tolerance int) bool {
	levels := strings.Split(report, " ")
	safe := true
	currentTolerance := tolerance

	firstLevel, _ := strconv.Atoi(levels[0])
	secondLevel, _ := strconv.Atoi(levels[1])
	safe = isDiffSafe(firstLevel, secondLevel)

	increasing := secondLevel > firstLevel
	previousLevel := secondLevel
	for i := 2; i < len(levels) && safe; i++ {
		currentLevel, _ := strconv.Atoi(levels[i])
		safeDiff := isDiffSafe(previousLevel, currentLevel)
		keepsIncreasing := increasing && currentLevel > previousLevel
		keepsDecreasing := !increasing && currentLevel < previousLevel
		if !safeDiff || !(keepsIncreasing || keepsDecreasing) {
			if currentTolerance > 0 {
				currentTolerance--
				continue
			}
			safe = false
			break
		}
		previousLevel = currentLevel
	}

	if !safe && tolerance > 0 {
		reportWithoutFirstLevel := report[len(levels[0])+1:]
		return isReportSafe(reportWithoutFirstLevel, tolerance-1)
	}

	return safe
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
