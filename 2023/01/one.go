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
	calibrationValues := findCalibrationValues(lines)
	sum := 0
	for _, calibrationValue := range calibrationValues {
		sum += calibrationValue
	}
	fmt.Println(sum)
}

func findCalibrationValues(lines []string) []int {
	calibrationValues := make([]int, 0, len(lines))
	for _, line := range lines {
		calibrationValue := findCalibrationValue(line)
		calibrationValues = append(calibrationValues, calibrationValue)
		fmt.Printf("%v -> %v\n", line, calibrationValue)
	}
	return calibrationValues
}

func findCalibrationValue(line string) int {
	digits := findDigits(line)
	return digits[0]*10 + digits[len(digits)-1]
}

func findDigits(line string) []int {
	preparedLine := strings.ReplaceAll(line, "one", "o1e")
	preparedLine = strings.ReplaceAll(preparedLine, "two", "t2o")
	preparedLine = strings.ReplaceAll(preparedLine, "three", "t3e")
	preparedLine = strings.ReplaceAll(preparedLine, "four", "4")
	preparedLine = strings.ReplaceAll(preparedLine, "five", "5e")
	preparedLine = strings.ReplaceAll(preparedLine, "six", "6")
	preparedLine = strings.ReplaceAll(preparedLine, "seven", "7n")
	preparedLine = strings.ReplaceAll(preparedLine, "eight", "8t")
	preparedLine = strings.ReplaceAll(preparedLine, "nine", "n9e")
	return findDigitsPart1(preparedLine)
}

func findDigitsPart1(line string) []int {
	digits := make([]int, 0, 2)
	for _, r := range []rune(line) {
		n, err := strconv.Atoi(string(r))
		if err == nil {
			digits = append(digits, n)
		}
	}
	return digits
}
