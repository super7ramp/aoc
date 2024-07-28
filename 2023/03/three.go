package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	emptySymbol = '.'
	gearSymbol  = '*'
)

var digits = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

type Pos struct {
	rowIndex int
	colIndex int
}

type NumberSlot struct {
	rowIndex      int
	colStartIndex int
	colEndIndex   int
}

func (ns *NumberSlot) NeighborPositions(schema *Schema) []Pos {
	// XXXXX
	// X123X
	// XXXXX
	rowCount := schema.RowCount()
	colCount := schema.ColCount()
	neighborPositions := make([]Pos, 0)
	if ns.rowIndex > 0 {
		for colIndex := max(ns.colStartIndex-1, 0); colIndex <= min(ns.colEndIndex, colCount-1); colIndex++ {
			neighborPositions = append(neighborPositions, Pos{rowIndex: ns.rowIndex - 1, colIndex: colIndex})
		}
	}
	if ns.colStartIndex > 0 {
		neighborPositions = append(neighborPositions, Pos{rowIndex: ns.rowIndex, colIndex: ns.colStartIndex - 1})
	}
	if ns.colEndIndex < colCount-1 {
		neighborPositions = append(neighborPositions, Pos{rowIndex: ns.rowIndex, colIndex: ns.colEndIndex})
	}
	if ns.rowIndex < rowCount-1 {
		for colIndex := max(ns.colStartIndex-1, 0); colIndex <= min(ns.colEndIndex, colCount-1); colIndex++ {
			neighborPositions = append(neighborPositions, Pos{rowIndex: ns.rowIndex + 1, colIndex: colIndex})
		}
	}
	return neighborPositions
}

func (ns *NumberSlot) HasSymbolAsNeighbor(schema *Schema) bool {
	for _, pos := range ns.NeighborPositions(schema) {
		if isNonEmptySymbol(schema.cells[pos.rowIndex][pos.colIndex]) {
			return true
		}
	}
	return false
}

func (ns *NumberSlot) IsInNeighborhood(schema *Schema, pos Pos) bool {
	neighborPositions := ns.NeighborPositions(schema)
	return slices.Contains(neighborPositions, pos)
}

func (ns *NumberSlot) Value(schema *Schema) int {
	runes := schema.cells[ns.rowIndex][ns.colStartIndex:ns.colEndIndex]
	value, _ := strconv.Atoi(string(runes))
	return value
}

type Gear struct {
	pos   Pos
	ratio int
}

type Schema struct {
	cells [][]rune
}

func (s *Schema) ColCount() int {
	return len(s.cells[0])
}

func (s *Schema) RowCount() int {
	return len(s.cells)
}

func (s *Schema) NumberSlots() []NumberSlot {
	var numberSlots []NumberSlot
	for rowIndex, row := range s.cells {
		slot := NumberSlot{rowIndex: rowIndex, colStartIndex: -1, colEndIndex: -1}
		for columnIndex, cell := range row {
			if isDigit(cell) {
				if slot.colStartIndex == -1 {
					slot.colStartIndex = columnIndex
				}
				slot.colEndIndex = columnIndex + 1
			} else {
				if slot.colStartIndex >= 0 {
					numberSlots = append(numberSlots, slot)
					slot = NumberSlot{rowIndex: rowIndex, colStartIndex: -1, colEndIndex: -1}
				}
			}
		}
		if slot.colStartIndex >= 0 {
			slot.colEndIndex = len(row)
			numberSlots = append(numberSlots, slot)
		}
	}
	return numberSlots
}

func (s *Schema) PartNumbers() []int {
	numberSlots := s.NumberSlots()
	partNumbers := make([]int, 0)
	for _, numberSlot := range numberSlots {
		if numberSlot.HasSymbolAsNeighbor(s) {
			numberPart := numberSlot.Value(s)
			partNumbers = append(partNumbers, numberPart)
		}
	}
	return partNumbers
}

func (s *Schema) PartNumberSlots() []NumberSlot {
	numberSlots := s.NumberSlots()
	partNumbers := make([]NumberSlot, 0)
	for _, numberSlot := range numberSlots {
		if numberSlot.HasSymbolAsNeighbor(s) {
			partNumbers = append(partNumbers, numberSlot)
		}
	}
	return partNumbers
}

func (s *Schema) Gears() []Gear {
	possibleGearPositions := make([]Pos, 0)
	for rowIndex, row := range s.cells {
		for colIndex, cell := range row {
			if cell == gearSymbol {
				possibleGearPositions = append(possibleGearPositions, Pos{rowIndex, colIndex})
			}
		}
	}
	gears := make([]Gear, 0)
	for _, possibleGearPos := range possibleGearPositions {
		neighborPartNumbers := make([]NumberSlot, 0)
		for _, partNumber := range s.PartNumberSlots() {
			if partNumber.IsInNeighborhood(s, possibleGearPos) {
				neighborPartNumbers = append(neighborPartNumbers, partNumber)
			}
		}
		if len(neighborPartNumbers) == 2 {
			gear := Gear{possibleGearPos, neighborPartNumbers[0].Value(s) * neighborPartNumbers[1].Value(s)}
			gears = append(gears, gear)
		}
	}
	return gears
}

func isDigit(r rune) bool {
	return slices.Contains(digits, r)
}

func isNonEmptySymbol(r rune) bool {
	return r != emptySymbol && !isDigit(r)
}

func newSchema(input string) *Schema {
	lines := strings.Split(input, "\n")
	cells := make([][]rune, len(lines[0]))
	for i, row := range lines {
		cells[i] = []rune(row)
	}
	return &Schema{cells}
}

//go:embed input.txt
var input string

func main() {
	schema := newSchema(input)
	partNumbers := schema.PartNumbers()
	fmt.Printf("Part numbers: %v\n\n", partNumbers)

	partNumberSum := 0
	for _, partNumber := range partNumbers {
		partNumberSum += partNumber
	}
	fmt.Printf("Part number sum: %v\n\n", partNumberSum)

	gears := schema.Gears()
	fmt.Printf("Gears: %v\n\n", gears)

	gearRatioSum := 0
	for _, gear := range gears {
		gearRatioSum += gear.ratio
	}
	fmt.Printf("Gear ratio sum: %v\n", gearRatioSum)
}
