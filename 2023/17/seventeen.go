package main

import (
	_ "embed"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

type Position struct {
	column int
	row    int
}

func Pos(column, row int) Position {
	return Position{column, row}
}

func PositionsAligned(a, b, c, d, e Position) bool {
	rowDifferences := abs(a.row-b.row) + abs(b.row-c.row) + abs(c.row-d.row) + abs(d.row-e.row)
	columnDifferences := abs(a.column-b.column) + abs(b.column-c.column) + abs(c.column-d.column) + abs(d.column-e.column)
	return rowDifferences == 0 || columnDifferences == 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type HeatLossMap [][]int

func NewPuzzleMap(input string) HeatLossMap {
	rows := strings.Split(input, "\n")
	rowCount := len(rows)
	columnCount := len(rows[0])
	heatLosses := make([][]int, rowCount)

	for rowIndex, row := range rows {
		heatLosses[rowIndex] = make([]int, columnCount)
		for columnIndex, char := range row {
			heatLoss, _ := strconv.Atoi(string(char))
			heatLosses[rowIndex][columnIndex] = heatLoss
		}
	}

	return heatLosses
}

func (heatLossMap *HeatLossMap) NeighborsOf(position Position) []Position {
	neighbors := make([]Position, 0)
	if position.row > 0 {
		neighbors = append(neighbors, Pos(position.column, position.row-1))
	}
	if position.row < heatLossMap.RowCount()-1 {
		neighbors = append(neighbors, Pos(position.column, position.row+1))
	}
	if position.column > 0 {
		neighbors = append(neighbors, Pos(position.column-1, position.row))
	}
	if position.column < heatLossMap.ColumnCount()-1 {
		neighbors = append(neighbors, Pos(position.column+1, position.row))
	}
	return neighbors
}

func (heatLossMap *HeatLossMap) Contains(position Position) bool {
	return position.row < heatLossMap.RowCount() && position.column < heatLossMap.ColumnCount()
}

func (heatLossMap *HeatLossMap) RowCount() int {
	return len(*heatLossMap)
}

func (heatLossMap *HeatLossMap) ColumnCount() int {
	return len((*heatLossMap)[0])
}

type NodeStatus struct {
	cumulatedHeatLoss int
	from              Position
	visited           bool
}

type DijsktraBasedShortestPathFinder struct {
	heatLossMap HeatLossMap
	statuses    map[Position]*NodeStatus
	current     Position
}

func NewDijsktraBasedShortestPathFinder(graph HeatLossMap) *DijsktraBasedShortestPathFinder {
	statuses := make(map[Position]*NodeStatus)
	current := Position{}
	return &DijsktraBasedShortestPathFinder{graph, statuses, current}
}

func (finder *DijsktraBasedShortestPathFinder) PathWithMinimalHeatLoss(from, to Position) (int, []Position) {
	defer clear(finder.statuses)
	*finder.status(from) = NodeStatus{cumulatedHeatLoss: 0, visited: true, from: from}
	finder.current = from

	for finder.current != to {
		for _, neighbor := range finder.nonVisitedNeighbors() {
			finder.updateNeighborCumulatedHeatLoss(neighbor)
		}
		newCurrent := finder.nonVisitedPositionWithMinimalCumulatedHeatLoss()
		finder.markVisited(newCurrent)
		finder.current = newCurrent
	}

	path := make([]Position, 0)
	for current := to; current != from; current = finder.previousOf(current) {
		path = append(path, current)
	}
	slices.Reverse(path)

	return finder.status(to).cumulatedHeatLoss, path
}

func (finder *DijsktraBasedShortestPathFinder) updateNeighborCumulatedHeatLoss(neighbor Position) {
	heatLossFromCurrent := finder.status(finder.current).cumulatedHeatLoss + finder.heatLossMap[neighbor.row][neighbor.column]
	neighborStatus := finder.status(neighbor)
	if heatLossFromCurrent < neighborStatus.cumulatedHeatLoss {
		neighborStatus.cumulatedHeatLoss = heatLossFromCurrent
		neighborStatus.from = finder.current
	}
}

func (finder *DijsktraBasedShortestPathFinder) nonVisitedNeighbors() []Position {
	neighbors := finder.heatLossMap.NeighborsOf(finder.current)
	slices.DeleteFunc(neighbors, func(position Position) bool {
		return finder.hasVisited(position) || finder.lastFourPositionsAlignedWith(position)
	})
	return neighbors
}

func (finder *DijsktraBasedShortestPathFinder) nonVisitedPositionWithMinimalCumulatedHeatLoss() Position {
	minimalHeatLoss := math.MaxInt
	var position Position
	for candidate, status := range finder.statuses {
		if !status.visited && status.cumulatedHeatLoss < minimalHeatLoss {
			minimalHeatLoss = status.cumulatedHeatLoss
			position = candidate
		}
	}
	return position
}

func (finder *DijsktraBasedShortestPathFinder) hasVisited(position Position) bool {
	return finder.status(position).visited
}

func (finder *DijsktraBasedShortestPathFinder) markVisited(position Position) {
	finder.status(position).visited = true
}

func (finder *DijsktraBasedShortestPathFinder) status(position Position) *NodeStatus {
	status, isPresent := finder.statuses[position]
	if !isPresent {
		status = &NodeStatus{cumulatedHeatLoss: math.MaxInt}
		finder.statuses[position] = status
	}
	return status
}

func (finder *DijsktraBasedShortestPathFinder) lastFourPositionsAlignedWith(position Position) bool {
	current := finder.current
	previous := finder.previousOf(current)
	beforePrevious := finder.previousOf(previous)
	beforeBeforePrevious := finder.previousOf(beforePrevious)
	if current == previous || previous == beforePrevious || beforePrevious == beforeBeforePrevious {
		// previous of start is itself, don't count it twice
		return false
	}
	return PositionsAligned(position, current, previous, beforePrevious, beforeBeforePrevious)
}

func (finder *DijsktraBasedShortestPathFinder) previousOf(position Position) Position {
	return finder.status(position).from
}

//go:embed input-example.txt
var input string

func main() {
	puzzleMap := NewPuzzleMap(input)
	dijsktra := NewDijsktraBasedShortestPathFinder(puzzleMap)
	from, to := Pos(0, 0), Pos(puzzleMap.ColumnCount()-1, puzzleMap.RowCount()-1)
	heatLoss, path := dijsktra.PathWithMinimalHeatLoss(from, to)
	fmt.Println("Heat loss: ", heatLoss)
	fmt.Println("Path: ", path)
	for i := 0; i < puzzleMap.RowCount(); i++ {
		for j := 0; j < puzzleMap.ColumnCount(); j++ {
			pos := Pos(i, j)
			if slices.Contains(path, pos) {
				fmt.Print(slices.Index(path, pos) % 10)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}
