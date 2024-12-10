package main

import (
	_ "embed"
	"fmt"
	"maps"
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
	dx := aa.antenna2.x - aa.antenna1.x
	dy := aa.antenna2.y - aa.antenna1.y
	antiNode1 := Pos{aa.antenna1.x - dx, aa.antenna1.y - dy}
	antiNode2 := Pos{aa.antenna2.x + dx, aa.antenna2.y + dy}
	var antiNodes []Pos
	if antiNode1.x >= 0 && antiNode1.x <= maxX && antiNode1.y >= 0 && antiNode1.y <= maxY {
		antiNodes = append(antiNodes, antiNode1)
	}
	if antiNode2.x >= 0 && antiNode2.x <= maxX && antiNode2.y >= 0 && antiNode2.y <= maxY {
		antiNodes = append(antiNodes, antiNode2)
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
	if len(currentRow) > 0 {
		tiles = append(tiles, currentRow)
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

func (m *AntennaMap) PrintAntiNodes() {
	antennaGroups := m.AntennaGroups()
	uniqueAntiNodes := make(map[Pos]struct{})
	for _, group := range antennaGroups {
		fmt.Printf("Antenna %c\n", group.frequency)
		fmt.Println("- Positions:", group.positions)
		fmt.Println("- Alignments:", group.Alignments())
		antiNodes := group.AntiNodes(m.Width()-1, m.Height()-1)
		fmt.Println("- Anti-nodes:", antiNodes)
		for _, antiNode := range antiNodes {
			uniqueAntiNodes[antiNode] = struct{}{}
		}
	}
	fmt.Println(len(uniqueAntiNodes), "distinct anti-nodes:", slices.Collect(maps.Keys(uniqueAntiNodes)))
	for y, row := range m.tiles {
		for x, tile := range row {
			if _, ok := uniqueAntiNodes[Pos{x, y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(string(tile))
			}
		}
		fmt.Println()
	}
}

//go:embed input.txt
var input []byte

func main() {
	antennaMap := AntennaMapFrom(input)
	antennaMap.PrintAntiNodes()
}
