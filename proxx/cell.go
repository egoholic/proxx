package proxx

import (
	"strconv"
)

type (
	openResult int
	cellCode   int
	cell       struct {
		code   cellCode // -1 for blackhole, 0-8 for number of black holes around
		row    int
		col    int
		opened bool
	}
)

const blackHoleCode cellCode = -1

func createCell(code int, row int, col int) *cell {
	return &cell{
		code:   cellCode(code),
		row:    row,
		col:    col,
		opened: false,
	}
}

const (
	closedMark    = "?"
	blackHoleMark = "*"
)

// Provides cell's content.
func (c *cell) toString() string {
	if !c.opened {
		return closedMark
	}
	switch c.code {
	case -1:
		return blackHoleMark
	default:
		return strconv.Itoa(int(c.code))
	}
}

// Safely increments cell's neighboring black holes counter.
func (c *cell) incrementNeighboringBlackHoles() {
	if c.code == blackHoleCode {
		return
	}
	c.code++
}

// Safely marks the cell as a black hole.
func (c *cell) markBlackHole() bool {
	// the same cell can't have two black holes
	if c.code == blackHoleCode {
		return false
	}
	c.code = blackHoleCode
	return true
}

const (
	wasOpened openResult = iota
	hasNoNeighboringBlackHoles
	hasNeighboringBlackHoles
	blackHole
)

// Marks the cell open and returns a resulting code.
func (c *cell) open() openResult {
	if c.opened {
		return wasOpened
	}
	c.opened = true
	switch c.code {
	case blackHoleCode:
		return blackHole
	case 0:
		return hasNoNeighboringBlackHoles
	default:
		return hasNeighboringBlackHoles
	}
}
