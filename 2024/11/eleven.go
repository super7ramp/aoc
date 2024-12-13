package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Stone int

type Stones []Stone

func StonesFrom(value string) Stones {
	var stones Stones
	for _, field := range strings.Split(value, " ") {
		number, _ := strconv.Atoi(field)
		stones = append(stones, Stone(number))
	}
	return stones
}

func (s *Stones) Blink(times int) {
	rules := []Rule{IfZeroThenOne, IfEvenNumberOfDigitsThenTwoStones, ElseMultiplyBy2024}
	for i := range times {
		fmt.Println("Current length: ", len(*s))
		fmt.Println("Blink", i+1)
		for i := 0; i < len(*s); i++ {
			for _, rule := range rules {
				if j, applied := rule(i, s); applied {
					i = j
					break
				}
			}
		}
	}
}

type Rule func(stoneIndex int, stones *Stones) (newIndex int, applied bool)

func IfZeroThenOne(stoneIndex int, stones *Stones) (int, bool) {
	if (*stones)[stoneIndex] != 0 {
		return stoneIndex, false
	}
	(*stones)[stoneIndex] = 1
	return stoneIndex, true
}

func IfEvenNumberOfDigitsThenTwoStones(stoneIndex int, stones *Stones) (int, bool) {
	stone := (*stones)[stoneIndex]
	digits := strconv.Itoa(int(stone))
	if len(digits)%2 != 0 {
		return stoneIndex, false
	}

	left := digits[:len(digits)/2]
	leftNumber, _ := strconv.Atoi(left)
	(*stones)[stoneIndex] = Stone(leftNumber)

	right := digits[len(digits)/2:]
	rightNumber, _ := strconv.Atoi(right)
	*stones = slices.Insert(*stones, stoneIndex+1, Stone(rightNumber))

	return stoneIndex + 1, true
}

func ElseMultiplyBy2024(stoneIndex int, stones *Stones) (int, bool) {
	(*stones)[stoneIndex] *= 2024
	return stoneIndex, true
}

//go:embed input.txt
var input string

func main() {
	stones := StonesFrom(input)
	fmt.Println("Initial stones:", stones)
	stones.Blink(25)
	fmt.Println(len(stones), "stones after blinking 25 times")
	stones.Blink(50)
	fmt.Println(len(stones), "stones after blinking 75 times")
}
