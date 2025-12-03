package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

// Orientation represents the current orientation (0-99).
type Orientation int

// Rotation represents a single rotation instruction, the sign indicating the rotation direction (negative for turn
// left, positive for turn right), the absolute value indicating the rotation distance.
type Rotation int

func ParseRotation(input string) Rotation {
	var sign int
	if input[0] == 'L' {
		sign = -1
	} else {
		sign = 1
	}
	value, _ := strconv.Atoi(input[1:])
	return Rotation(sign * value)
}

type Rotations []Rotation

func ParseRotations(input string) Rotations {
	var instructions Rotations
	for _, line := range strings.Split(input, "\n") {
		if len(line) > 1 {
			instruction := ParseRotation(line)
			instructions = append(instructions, instruction)
		}
	}
	return instructions
}

// CountPointedAtZeroFrom applies the rotations starting from the given initial orientation and returns the number of
// times the rotations pointed at orientation zero.
func (r Rotations) CountPointedAtZeroFrom(initialOrientation Orientation) int {
	orientation := int(initialOrientation)
	pointedAtZeroCount := 0
	for _, instruction := range r {
		orientation = mod(orientation+int(instruction), 100)
		if orientation == 0 {
			pointedAtZeroCount++
		}
	}
	return pointedAtZeroCount
}

// CountCrossedZeroFrom applies the rotations starting from the given initial orientation and returns the number of
// times the rotations crossed orientation zero.
func (r Rotations) CountCrossedZeroFrom(initialOrientation Orientation) int {
	orientation := int(initialOrientation)
	crossedZeroCount := 0
	for _, rotation := range r {
		crossedZeroCount += countCrossedZero(orientation, rotation)
		orientation = mod(orientation+int(rotation), 100)
	}
	return crossedZeroCount
}

func countCrossedZero(initialOrientation int, rotation Rotation) int {
	crossedZeroCount := abs(int(rotation)) / 100
	rotationLeftover := int(rotation) % 100
	unboundOrientation := initialOrientation + rotationLeftover
	if (unboundOrientation <= 0 && initialOrientation != 0) || unboundOrientation >= 100 {
		crossedZeroCount++
	}
	return crossedZeroCount
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	fmt.Println("Part 1:")
	rotations := ParseRotations(input)
	initialOrientation := Orientation(50)
	pointedAtZeroCount := rotations.CountPointedAtZeroFrom(initialOrientation)
	fmt.Printf("Pointed at orientation zero %v times\n", pointedAtZeroCount)

	fmt.Println("Part 2:")
	crossedZeroCount := rotations.CountCrossedZeroFrom(initialOrientation)
	fmt.Printf("Crossed orientation zero %v times\n", crossedZeroCount)
}
