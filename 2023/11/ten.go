package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

type Element rune

const (
	galaxy = Element('#')
)

type Position struct {
	rowIndex    int
	columnIndex int
}

func (p *Position) stepDistanceTo(other Position) int {
	return abs(p.rowIndex-other.rowIndex) + abs(p.columnIndex-other.columnIndex)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Universe struct {
	space       [][]Element
	columnCount int
	rowCount    int
}

func (s *Universe) String() string {
	var builder strings.Builder
	for _, row := range s.space {
		for _, element := range row {
			builder.WriteRune(rune(element))
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}

func (s *Universe) Expand() *Universe {
	rowsWithoutGalaxy := s.RowsWithoutGalaxy()
	expandedSpaceRowCount := s.rowCount + len(rowsWithoutGalaxy)

	columnsWithoutGalaxy := s.ColumnsWithoutGalaxy()
	expandedSpaceColumnCount := s.columnCount + len(columnsWithoutGalaxy)

	expandedSpace := make([][]Element, expandedSpaceRowCount)
	for rowIndex, addedRows := 0, 0; rowIndex < s.rowCount; rowIndex++ {
		expandedSpace[rowIndex+addedRows] = make([]Element, expandedSpaceColumnCount)
		for columnIndex, addedColumns := 0, 0; columnIndex < s.columnCount; columnIndex++ {
			expandedSpace[rowIndex+addedRows][columnIndex+addedColumns] = s.space[rowIndex][columnIndex]
			if slices.Contains(columnsWithoutGalaxy, columnIndex) {
				expandedSpace[rowIndex+addedRows][columnIndex+addedColumns+1] = s.space[rowIndex][columnIndex]
				addedColumns++
			}
		}
		if slices.Contains(rowsWithoutGalaxy, rowIndex) {
			expandedSpace[rowIndex+addedRows+1] = make([]Element, expandedSpaceColumnCount)
			copy(expandedSpace[rowIndex+addedRows+1], expandedSpace[rowIndex+addedRows])
			addedRows++
		}
	}
	return &Universe{
		space:       expandedSpace,
		columnCount: s.columnCount + len(columnsWithoutGalaxy),
		rowCount:    s.rowCount + len(rowsWithoutGalaxy),
	}

}

func (s *Universe) RowsWithoutGalaxy() []int {
	var rowsWithoutGalaxy []int
	for rowIndex, row := range s.space {
		if !slices.Contains(row, galaxy) {
			rowsWithoutGalaxy = append(rowsWithoutGalaxy, rowIndex)
		}
	}
	return rowsWithoutGalaxy
}

func (s *Universe) ColumnsWithoutGalaxy() []int {
	var columnsWithoutGalaxy []int
	for columnIndex := 0; columnIndex < s.columnCount; columnIndex++ {
		column := make([]Element, s.rowCount)
		for rowIndex := 0; rowIndex < s.rowCount; rowIndex++ {
			column[rowIndex] = s.space[rowIndex][columnIndex]
		}
		if !slices.Contains(column, galaxy) {
			columnsWithoutGalaxy = append(columnsWithoutGalaxy, columnIndex)
		}
	}
	return columnsWithoutGalaxy
}

func (s *Universe) GalaxyPositions() [](Position) {
	var galaxies []Position
	for rowIndex, row := range s.space {
		for columnIndex, element := range row {
			if element == galaxy {
				galaxies = append(galaxies, Position{rowIndex, columnIndex})
			}
		}
	}
	return galaxies
}

func UniverseFrom(input string) *Universe {
	lines := strings.Split(input, "\n")
	space := make([][]Element, len(lines))
	for rowIndex, line := range lines {
		space[rowIndex] = make([]Element, len(line))
		for columnIndex, char := range line {
			space[rowIndex][columnIndex] = Element(char)
		}
	}
	width := len(lines[0])
	height := len(lines)
	return &Universe{space, width, height}
}

type FastExpansionUniverse struct {
	galaxyPositions []Position
}

func FastExpansionUniverseFrom(universe *Universe) *FastExpansionUniverse {
	galaxyPositions := universe.GalaxyPositions()
	columnsWithoutGalaxy := universe.ColumnsWithoutGalaxy()
	rowsWithoutGalaxy := universe.RowsWithoutGalaxy()

	for i, galaxyPosition := range galaxyPositions {
		var columnsWithoutGalaxyBeforeGalaxy int
		for _, columnIndex := range columnsWithoutGalaxy {
			if columnIndex < galaxyPosition.columnIndex {
				columnsWithoutGalaxyBeforeGalaxy++
			}
		}
		var rowsWithoutGalaxyBeforeGalaxy int
		for _, rowIndex := range rowsWithoutGalaxy {
			if rowIndex < galaxyPosition.rowIndex {
				rowsWithoutGalaxyBeforeGalaxy++
			}
		}
		galaxyPositions[i] = Position{
			galaxyPosition.rowIndex + rowsWithoutGalaxyBeforeGalaxy*999_999,
			galaxyPosition.columnIndex + columnsWithoutGalaxyBeforeGalaxy*999_999,
		}
	}

	return &FastExpansionUniverse{galaxyPositions}
}

//go:embed input.txt
var input string

func main() {
	universe := UniverseFrom(input)
	fmt.Println("Initial universe:")
	fmt.Println(universe.String())

	expandedSpace := universe.Expand()
	fmt.Println("Expanded universe:")
	fmt.Println(expandedSpace.String())

	galaxyPositions := expandedSpace.GalaxyPositions()
	fmt.Println("Galaxy positions:", galaxyPositions)

	var galaxyStepDistanceSum int
	for i := 0; i < len(galaxyPositions); i++ {
		for j := i + 1; j < len(galaxyPositions); j++ {
			galaxyStepDistanceSum += galaxyPositions[i].stepDistanceTo(galaxyPositions[j])
		}
	}
	fmt.Println("Galaxy step distance sum (part 1): ", galaxyStepDistanceSum)

	fastExpansionUniverse := FastExpansionUniverseFrom(universe)
	galaxyPositions = fastExpansionUniverse.galaxyPositions
	galaxyStepDistanceSum = 0
	for i := 0; i < len(galaxyPositions); i++ {
		for j := i + 1; j < len(galaxyPositions); j++ {
			galaxyStepDistanceSum += galaxyPositions[i].stepDistanceTo(galaxyPositions[j])
		}
	}
	fmt.Println("Galaxy step distance sum (part 2): ", galaxyStepDistanceSum)
}
