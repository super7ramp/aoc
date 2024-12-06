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
	visitedPositions := patrolMap.VisitGuardPositions()
	fmt.Println("Guard visited", len(visitedPositions), "positions:", visitedPositions)

	occurrences := make(map[Pos]struct{})
	distinctPositionCount := 0
	for _, pos := range visitedPositions {
		if _, seen := occurrences[pos]; !seen {
			distinctPositionCount++
			occurrences[pos] = struct{}{}
		}
	}
	fmt.Println("Guard visited", distinctPositionCount, "distinct positions")
}
