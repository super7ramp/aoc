package main

import (
	_ "embed"
	"fmt"
	"maps"
	"math"
	"slices"
)

type Pos struct {
	x, y int
}

type AntennaGroup struct {
	frequency byte
	positions []Pos
}

func (group *AntennaGroup) Alignments() []AntennaAlignment {
	alignments := make([]AntennaAlignment, 0)
	for i, pos1 := range group.positions {
		for _, pos2 := range group.positions[i+1:] {
			alignments = append(alignments, AntennaAlignment{pos1, pos2})
		}
	}
	return alignments
}

func (group *AntennaGroup) AntiNodes(maxX, maxY int) []Pos {
	alignments := group.Alignments()
	antiNodes := make(map[Pos]struct{})
	for _, alignment := range alignments {
		for _, antiNode := range alignment.antiNodes(maxX, maxY) {
			antiNodes[antiNode] = struct{}{}
		}
	}
	return slices.Collect(maps.Keys(antiNodes))
}

type AntennaAlignment struct {
	antenna1, antenna2 Pos
}

func (aa *AntennaAlignment) antiNodes(maxX, maxY int) []Pos {
	a := float64(aa.antenna2.y-aa.antenna1.y) / float64(aa.antenna2.x-aa.antenna1.x)
	b := float64(aa.antenna1.y) - a*float64(aa.antenna1.x)
	antiNodes := make([]Pos, 0)
	for x := range maxX {
		y := a*float64(x) + b
		if y >= 0 && y < float64(maxY) {
			_, frac := math.Modf(y)
			if frac < 1e-5 {
				antiNodes = append(antiNodes, Pos{x, int(y)})
			}
		}
	}
	return antiNodes
}

type AntennaMap struct {
	tiles [][]byte
}

func AntennaMapFrom(input []byte) AntennaMap {
	tiles := make([][]byte, 0)
	currentRow := make([]byte, 0)
	for _, char := range input {
		if char == '\n' {
			tiles = append(tiles, slices.Clone(currentRow))
			currentRow = currentRow[:0]
		} else {
			currentRow = append(currentRow, char)
		}
	}
	return AntennaMap{tiles}
}

func (m *AntennaMap) Height() int {
	return len(m.tiles)
}

func (m *AntennaMap) Width() int {
	return len(m.tiles[0])
}

func (m *AntennaMap) AntennaGroups() []AntennaGroup {
	antenna := make(map[byte]AntennaGroup)
	for y, row := range m.tiles {
		for x, tile := range row {
			if tile != '.' {
				if group, ok := antenna[tile]; ok {
					group.positions = append(group.positions, Pos{x, y})
					antenna[tile] = group
				} else {
					antenna[tile] = AntennaGroup{tile, []Pos{{x, y}}}
				}
			}
		}
	}
	return slices.Collect(maps.Values(antenna))
}

//go:embed input-example.txt
var input []byte

func main() {
	antennaMap := AntennaMapFrom(input)
	antennaGroups := antennaMap.AntennaGroups()
	for _, group := range antennaGroups {
		fmt.Println("Antenna", group.frequency)
		fmt.Println("- Positions:", group.positions)
		fmt.Println("- Alignments:", group.Alignments())
		fmt.Println("- Anti-nodes:", group.AntiNodes(antennaMap.Width(), antennaMap.Height()))
	}
}
