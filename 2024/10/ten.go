package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"maps"
)

type Pos struct {
	x, y int
}

func (p *Pos) AdjacentPositions() []Pos {
	return []Pos{{p.x - 1, p.y}, {p.x + 1, p.y}, {p.x, p.y - 1}, {p.x, p.y + 1}}
}

type Trail []Pos

func (t Trail) Arrival() Pos {
	return t[len(t)-1]
}

type TrailHead struct {
	trails []Trail
}

func (t *TrailHead) AddTrail(trail Trail) {
	t.trails = append(t.trails, trail)
}

func (t *TrailHead) Score() int {
	distinctArrivals := make(map[Pos]struct{})
	for _, trail := range t.trails {
		distinctArrivals[trail.Arrival()] = struct{}{}
	}
	return len(distinctArrivals)
}

func (t *TrailHead) Rating() int {
	return len(t.trails)
}

type TopographicMap struct {
	levels []byte
	width  int
}

const (
	StartLevel = 0
	EndLevel   = 9
)

func ParseTopographicMap(input []byte) *TopographicMap {
	width := bytes.IndexByte(input, '\n')
	return &TopographicMap{input, width}
}

func (t *TopographicMap) TrailHeads() map[Pos]TrailHead {
	trailHeads := make(map[Pos]TrailHead)
	for x := 0; x < t.width; x++ {
		for y := 0; y < t.Height(); y++ {
			if t.Level(x, y) == StartLevel {
				trails := t.trailsFrom(x, y)
				for _, trail := range trails {
					if trailHead, ok := trailHeads[Pos{x, y}]; ok {
						trailHead.AddTrail(trail)
						trailHeads[Pos{x, y}] = trailHead
					} else {
						trailHeads[Pos{x, y}] = TrailHead{[]Trail{trail}}
					}
				}
			}
		}
	}
	return trailHeads
}

func (t *TopographicMap) trailsFrom(x, y int) []Trail {
	currentPos := Pos{x, y}
	currentLevel := t.Level(x, y)
	if currentLevel == EndLevel {
		return []Trail{[]Pos{currentPos}}
	}
	var nextTrails []Trail
	for _, nextPos := range currentPos.AdjacentPositions() {
		if nextPos.x < 0 || nextPos.x >= t.width {
			continue
		}
		if nextPos.y < 0 || nextPos.y >= t.Height() {
			continue
		}
		if t.Level(nextPos.x, nextPos.y) == currentLevel+1 {
			subNextTrails := t.trailsFrom(nextPos.x, nextPos.y)
			for _, subNextTrail := range subNextTrails {
				positions := []Pos{{x, y}}
				positions = append(positions, subNextTrail...)
				nextTrail := Trail(positions)
				nextTrails = append(nextTrails, nextTrail)
			}
		}
	}
	return nextTrails
}

func (t *TopographicMap) Level(x, y int) byte {
	return t.levels[y*(t.width+1)+x] - '0'
}

func (t *TopographicMap) Width() int {
	return t.width
}

func (t *TopographicMap) Height() int {
	return len(t.levels) / t.width
}

func (t *TopographicMap) String() string {
	return string(t.levels)
}

//go:embed input.txt
var input []byte

func main() {
	tm := ParseTopographicMap(input)

	trailHeads := tm.TrailHeads()
	fmt.Println("(Part 1) Trail heads:", trailHeads)
	scoreSum := 0
	for trailHead := range maps.Values(trailHeads) {
		scoreSum += trailHead.Score()
	}
	fmt.Println("(Part 1) Score sum:", scoreSum)

	ratingSum := 0
	for trailHead := range maps.Values(trailHeads) {
		ratingSum += trailHead.Rating()
	}
	fmt.Println("(Part 2) Rating sum:", ratingSum)
}
