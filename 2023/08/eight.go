package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
)

const (
	left  = Direction(0)
	right = Direction(1)
)

var crossingRe = regexp.MustCompile(`(?P<location>\w+) = \((?P<left>\w+), (?P<right>\w+)\)`)

type Direction int

type Location string

type Crossing struct {
	location Location
	onLeft   Location
	onRight  Location
}

type Puzzle struct {
	directions []Direction
	crossings  map[Location]Crossing
}

func (puzzle *Puzzle) RequiredSteps() int {
	steps := 0
	for crossing := puzzle.crossings["AAA"]; crossing.location != "ZZZ"; steps++ {
		if puzzle.directions[steps%len(puzzle.directions)] == left {
			crossing = puzzle.crossings[crossing.onLeft]
		} else {
			crossing = puzzle.crossings[crossing.onRight]
		}
	}
	return steps
}

func (puzzle *Puzzle) RequiredStepsForAGhost() int64 {
	var startCrossings []Crossing
	for location, crossing := range puzzle.crossings {
		if strings.HasSuffix(string(location), "A") {
			startCrossings = append(startCrossings, crossing)
		}
	}

	stepsPerStartPoint := make([]int64, len(startCrossings))
	for crossings, steps, remaining := startCrossings, int64(0), len(startCrossings); remaining > 0; steps++ {
		direction := puzzle.directions[steps%int64(len(puzzle.directions))]
		for i, crossing := range crossings {
			if strings.HasSuffix(string(crossing.location), "Z") {
				if stepsPerStartPoint[i] == 0 {
					stepsPerStartPoint[i] = steps
					remaining--
				}
				continue
			}
			if direction == left {
				crossings[i] = puzzle.crossings[crossing.onLeft]
			} else {
				crossings[i] = puzzle.crossings[crossing.onRight]
			}
		}
	}

	if len(stepsPerStartPoint) == 1 {
		return stepsPerStartPoint[0]
	}

	return lcm(stepsPerStartPoint[0], stepsPerStartPoint[1], stepsPerStartPoint[2:]...)
}

func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int64, integers ...int64) int64 {
	result := a * b / gcd(a, b)
	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}
	return result
}

func ParsePuzzle(input string) *Puzzle {
	sections := strings.Split(input, "\n\n")
	directions := parseDirections(sections[0])
	crossings := parseCrossings(sections[1])
	return &Puzzle{directions, crossings}
}

func parseDirections(s string) []Direction {
	directions := make([]Direction, len(s))
	for i, r := range s {
		if r == 'L' {
			directions[i] = left
		} else {
			directions[i] = right
		}
	}
	return directions
}

func parseCrossings(s string) map[Location]Crossing {
	lines := strings.Split(s, "\n")
	crossings := make(map[Location]Crossing)
	for _, line := range lines {
		crossing := parseCrossing(line)
		crossings[crossing.location] = crossing
	}
	return crossings
}

func parseCrossing(s string) Crossing {
	matches := crossingRe.FindStringSubmatch(s)
	return Crossing{
		location: Location(matches[crossingRe.SubexpIndex("location")]),
		onLeft:   Location(matches[crossingRe.SubexpIndex("left")]),
		onRight:  Location(matches[crossingRe.SubexpIndex("right")]),
	}
}

//go:embed input.txt
var input string

func main() {
	puzzle := ParsePuzzle(input)
	fmt.Println("Puzzle: ", puzzle)
	fmt.Println("Required steps: ", puzzle.RequiredSteps())
	fmt.Println("Required steps for a ghost: ", puzzle.RequiredStepsForAGhost())
}
