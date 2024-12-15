package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

type Pos struct {
	x, y int
}

func (p Pos) AdjacentPositions() []Pos {
	return []Pos{
		{p.x + 1, p.y}, // right
		{p.x, p.y + 1}, // down
		{p.x - 1, p.y}, // left
		{p.x, p.y - 1}, // up
	}
}

func (p Pos) IsWithin(width, height int) bool {
	return p.x >= 0 && p.y >= 0 && p.x < width && p.y < height
}

type Plant byte

type Region struct {
	plant Plant
	plots []Pos
}

func (r Region) Area() int {
	return len(r.plots)
}

func (r Region) Perimeter() int {
	perimeter := 0
	for _, plot := range r.plots {
		for _, adjacent := range plot.AdjacentPositions() {
			if !slices.Contains(r.plots, adjacent) {
				perimeter++
			}
		}
	}
	return perimeter
}

func (r Region) FencingPrice() int {
	return r.Perimeter() * r.Area()
}

func (r Region) String() string {
	return fmt.Sprintf("%c: %v", r.plant, r.plots)
}

type Garden [][]Plant

func GardenFrom(value string) Garden {
	var garden Garden
	for _, row := range strings.Split(value, "\n") {
		garden = append(garden, []Plant(row))
	}
	return garden
}

func (garden Garden) Regions() []Region {
	var positionsToVisit []Pos
	for y, row := range garden {
		for x := range row {
			positionsToVisit = append(positionsToVisit, Pos{x, y})
		}
	}
	var regions []Region
	for len(positionsToVisit) > 0 {
		plant := garden.PlantAt(positionsToVisit[0])
		region := Region{plant, []Pos{positionsToVisit[0]}}
		positionsToVisit = slices.Delete(positionsToVisit, 0, 1)
		for i := 0; i < len(region.plots); i++ {
			current := region.plots[i]
			for _, adjacent := range current.AdjacentPositions() {
				if !adjacent.IsWithin(garden.Width(), garden.Height()) || garden.PlantAt(adjacent) != plant {
					continue
				}
				if visitedIndex := slices.Index(positionsToVisit, adjacent); visitedIndex >= 0 {
					region.plots = append(region.plots, adjacent)
					positionsToVisit = slices.Delete(positionsToVisit, visitedIndex, visitedIndex+1)
				}
			}
		}
		regions = append(regions, region)
	}
	return regions
}

func (garden Garden) PlantAt(p Pos) Plant {
	return garden[p.y][p.x]
}

func (garden Garden) Width() int {
	return len(garden[0])
}

func (garden Garden) Height() int {
	return len(garden)
}

//go:embed input.txt
var input string

func main() {
	garden := GardenFrom(input)

	totalFencingPrice := 0
	for _, region := range garden.Regions() {
		fmt.Printf("(Part 1) A region of %c plants with price %d * %d = %d.\n", region.plant, region.Area(), region.Perimeter(), region.FencingPrice())
		totalFencingPrice += region.FencingPrice()
	}
	fmt.Println("(Part 1) The total fencing price is", totalFencingPrice)

}
