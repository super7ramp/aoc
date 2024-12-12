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

type TrailHead struct {
	trails []Trail
}

func (t *TrailHead) AddTrail(trail Trail) {
	t.trails = append(t.trails, trail)
}

func (t *TrailHead) Score() int {
	distinctArrivals := make(map[Pos]struct{})
	for _, trail := range t.trails {
		distinctArrivals[trail[len(trail)-1]] = struct{}{}
	}
	return len(distinctArrivals)
}

func (t *TrailHead) Rating() int {
	return len(t.trails)
}

type TopographicMap struct {
	tiles []byte
	width int
}

func ParseTopographicMap(input []byte) *TopographicMap {
	width := bytes.IndexByte(input, '\n')
	return &TopographicMap{input, width}
}

func (t *TopographicMap) TrailHeads() map[Pos]TrailHead {
	trailHeads := make(map[Pos]TrailHead)
	for x := 0; x < t.width; x++ {
		for y := 0; y < t.Height(); y++ {
			if t.Get(x, y) == 0 {
				trails := t.nextTrails(x, y)
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

func (t *TopographicMap) nextTrails(x, y int) []Trail {
	currentPos := Pos{x, y}
	current := t.Get(x, y)
	if current == 9 {
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
		if t.Get(nextPos.x, nextPos.y) == current+1 {
			subNextTrails := t.nextTrails(nextPos.x, nextPos.y)
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

func (t *TopographicMap) Get(x, y int) byte {
	return t.tiles[y*(t.width+1)+x] - '0'
}

func (t *TopographicMap) Width() int {
	return t.width
}

func (t *TopographicMap) Height() int {
	return len(t.tiles) / t.width
}

func (t *TopographicMap) String() string {
	return string(t.tiles)
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
