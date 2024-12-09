package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

const (
	Empty           = '.'
	Obstacle        = '#'
	GuardGoingUp    = '^'
	GuardGoingRight = '>'
	GuardGoingDown  = 'v'
	GuardGoingLeft  = '<'
)

type Direction byte

func (d Direction) next() Direction {
	switch d {
	case GuardGoingUp:
		return GuardGoingRight
	case GuardGoingRight:
		return GuardGoingDown
	case GuardGoingDown:
		return GuardGoingLeft
	case GuardGoingLeft:
		return GuardGoingUp
	}
	panic("Invalid direction")
}

type Pos struct {
	x, y int
}

func (p *Pos) Up() Pos {
	return Pos{p.x, p.y - 1}
}

func (p *Pos) Down() Pos {
	return Pos{p.x, p.y + 1}
}

func (p *Pos) Left() Pos {
	return Pos{p.x - 1, p.y}
}

func (p *Pos) Right() Pos {
	return Pos{p.x + 1, p.y}
}

type PatrolMap struct {
	tiles [][]byte
}

func PatrolMapFrom(bytes []byte) *PatrolMap {
	tiles := make([][]byte, 0)
	currentRow := make([]byte, 0)
	for _, b := range bytes {
		if b == '\n' {
			tiles = append(tiles, slices.Clone(currentRow))
			currentRow = slices.Delete(currentRow, 0, len(currentRow))
		} else {
			currentRow = append(currentRow, b)
		}
	}
	return &PatrolMap{tiles}
}

func (m *PatrolMap) VisitGuardPositions() []Pos {
	visited := make([]Pos, 0)
	for guardPosition := m.guardPosition(); m.contains(&guardPosition); guardPosition = m.nextGuardPosition(&guardPosition) {
		visited = append(visited, guardPosition)
		//fmt.Println("Visited:", guardPosition)
		//fmt.Println(m.String())
	}
	return visited
}

func (m *PatrolMap) PossibleObstructions() []Pos {
	visitedPositions := m.Clone().VisitGuardPositions()
	possibleObstructions := make([]Pos, 0)
	for _, pos := range visitedPositions {
		if m.DoesObstructionMakeGuardLoop(&pos) {
			possibleObstructions = append(possibleObstructions, pos)
		}
	}
	return possibleObstructions
}

func (m *PatrolMap) DoesObstructionMakeGuardLoop(obstruction *Pos) bool {
	if m.guardPosition() == *obstruction {
		return false
	}

	probeMap := m.Clone()
	probeMap.setTileAt(obstruction, Obstacle)
	visited := make(map[Pos][]Direction)

	for guardPosition := probeMap.guardPosition(); probeMap.contains(&guardPosition); guardPosition = probeMap.nextGuardPosition(&guardPosition) {
		previousGuardDirectionsOnThisPosition, seen := visited[guardPosition]
		currentGuardDirectionOnThisPosition := Direction(probeMap.getTileAt(&guardPosition))
		if seen && slices.Contains(previousGuardDirectionsOnThisPosition, currentGuardDirectionOnThisPosition) {
			return true
		}
		if !seen {
			previousGuardDirectionsOnThisPosition = make([]Direction, 0, 1)
		}
		previousGuardDirectionsOnThisPosition = append(previousGuardDirectionsOnThisPosition, currentGuardDirectionOnThisPosition)
		visited[guardPosition] = previousGuardDirectionsOnThisPosition
	}
	return false
}

func (m *PatrolMap) Clone() *PatrolMap {
	copiedTiles := make([][]byte, len(m.tiles))
	for i, row := range m.tiles {
		copiedTiles[i] = slices.Clone(row)
	}
	return &PatrolMap{copiedTiles}
}

func (m *PatrolMap) String() string {
	sb := strings.Builder{}
	for _, row := range m.tiles {
		for _, b := range row {
			_, _ = fmt.Fprintf(&sb, "%c", b)
		}
		_, _ = fmt.Fprintln(&sb)
	}
	return sb.String()
}

func (m *PatrolMap) guardPosition() Pos {
	for y, row := range m.tiles {
		for x, tile := range row {
			if tile == GuardGoingUp || tile == GuardGoingRight || tile == GuardGoingDown || tile == GuardGoingLeft {
				return Pos{x, y}
			}
		}
	}
	return Pos{-1, -1}
}

func (m *PatrolMap) getTileAt(pos *Pos) byte {
	return m.tiles[pos.y][pos.x]
}

func (m *PatrolMap) setTileAt(pos *Pos, tile byte) {
	m.tiles[pos.y][pos.x] = tile
}

func (m *PatrolMap) contains(pos *Pos) bool {
	return pos.x >= 0 && pos.x < m.width() && pos.y >= 0 && pos.y < m.height()
}

func (m *PatrolMap) width() int {
	return len(m.tiles[0])
}

func (m *PatrolMap) height() int {
	return len(m.tiles)
}

func (m *PatrolMap) nextGuardPosition(current *Pos) Pos {
	var next Pos
	guard := m.getTileAt(current)
	switch guard {
	case GuardGoingUp:
		next = current.Up()
	case GuardGoingRight:
		next = current.Right()
	case GuardGoingDown:
		next = current.Down()
	case GuardGoingLeft:
		next = current.Left()
	default:
		panic("Invalid guard position")
	}
	if !m.contains(&next) {
		m.setTileAt(current, Empty)
		return next
	}
	if m.getTileAt(&next) == Obstacle {
		newDirection := Direction(guard).next()
		guardWithNewDirection := byte(newDirection)
		m.setTileAt(current, guardWithNewDirection)
		return m.nextGuardPosition(current)
	}
	m.setTileAt(current, Empty)
	m.setTileAt(&next, guard)
	return next
}

//go:embed input.txt
var input []byte

func main() {
	patrolMap := PatrolMapFrom(input)
	visitedPositions := patrolMap.Clone().VisitGuardPositions()
	fmt.Println("(Part 1) Guard visited", len(visitedPositions), "positions:", visitedPositions)

	occurrences := make(map[Pos]struct{})
	distinctPositionCount := 0
	for _, pos := range visitedPositions {
		if _, seen := occurrences[pos]; !seen {
			distinctPositionCount++
			occurrences[pos] = struct{}{}
		}
	}
	fmt.Println("(Part 1) Guard visited", distinctPositionCount, "distinct positions")

	possibleObstructions := patrolMap.PossibleObstructions()
	fmt.Println("(Part 2) Possible obstructions:", possibleObstructions)

	clear(occurrences)
	distinctPossibleObstructionCount := 0
	for _, pos := range possibleObstructions {
		if _, seen := occurrences[pos]; !seen {
			distinctPossibleObstructionCount++
			occurrences[pos] = struct{}{}
		}
	}
	fmt.Println("(Part 2)", distinctPossibleObstructionCount, "distinct possible obstructions")
}
