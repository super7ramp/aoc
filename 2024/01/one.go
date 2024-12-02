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
	column1, column2 := columns()

	// Part 1
	slices.Sort(column1)
	slices.Sort(column2)
	fmt.Println("Sorted column 1:", column1)
	fmt.Println("Sorted column 2:", column2)
	differenceSum := 0
	for i := range column1 {
		differenceSum += abs(column2[i] - column1[i])
	}
	fmt.Println("Sum of differences:", differenceSum)

	// Part 2
	similarityScore := 0
	for _, number := range column1 {
		similarityScore += number * countNumber(number, column2)
	}
	fmt.Println("Similarity score:", similarityScore)
}

func columns() ([]int, []int) {
	lines := strings.Split(input, "\n")
	column1 := make([]int, len(lines))
	column2 := make([]int, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, "   ")
		column1[i], _ = strconv.Atoi(parts[0])
		column2[i], _ = strconv.Atoi(parts[1])
	}
	return column1, column2
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func countNumber(number int, numbers []int) int {
	count := 0
	for _, n := range numbers {
		if n == number {
			count++
		}
	}
	return count
}
