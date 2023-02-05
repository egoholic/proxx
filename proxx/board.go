package proxx

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type (
	board struct {
		size  int       // size * size board size
		cells [][]*cell // 2D slice to store cells within rows
		// here we save black holes indeces for performance optimisation
		// instead of searching black holes near every cell, we will only update cells near black holes
		blackHoles []*cell
		opened     int // counter for opened cells, if opened == len(cells) - len(blackHoles) then the player win (see: *Board.CheckWin())
	}
	operation func(*board, *cell)
)

const (
	minBoardSize = 9
	maxBoardSize = 24
)

// Here I assume the min number of black holes to be no less than square root of cells number (board size)
func minBlackHolesNum(s int) int {
	return s
}

// Here I assume the max number of black holes to be no greater than square root of cells number (board size) * 2
func maxBlackHolesNum(s int) int {
	return s * 2
}

// Incremets cell's neigboring black holes counter.
// The function implements oparetion type to be used within doWithNeighboringCells.
func incrementNeighboringBlackHolesOp(_ *board, c *cell) {
	c.incrementNeighboringBlackHoles()
	return
}

// Call it only for cells that have no black holes around.
func recursivelyOpenSafeNeighbors(b *board, c *cell) {
	res := c.open()
	if res != wasOpened {
		b.opened++
	}
	if c.code == 0 {
		doWithNeighboringCells(b, c, recursivelyOpenSafeNeighbors)
	}
}

// Performs given operation on neighboring cells.
// Ignores already opened cells to avoid infinite loops and because it has no sense from the busines/domain PoV.
func doWithNeighboringCells(b *board, c *cell, op operation) {
	// Special cases come first.
	var cel *cell
	// board's top edge cell
	if c.row == 0 {
		// board's top left cell
		if c.col == 0 {
			cel = b.cells[0][1]
			if !cel.opened {
				op(b, cel)
			}
			cel = b.cells[1][0]
			if !cel.opened {
				op(b, cel)
			}
			cel = b.cells[1][1]
			if !cel.opened {
				op(b, cel)
			}
			return
		}
		// board's top right cell
		if c.col == b.size-1 {
			cel = b.cells[0][b.size-2]
			if !cel.opened {
				op(b, cel)
			}
			cel = b.cells[1][b.size-1]
			if !cel.opened {
				op(b, cel)
			}
			cel = b.cells[1][b.size-2]
			if !cel.opened {
				op(b, cel)
			}
			return
		}
		// board's top middle cell
		cel = b.cells[0][c.col-1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[0][c.col+1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[1][c.col-1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[1][c.col]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[1][c.col+1]
		if !cel.opened {
			op(b, cel)
		}
		return
	}
	// board's bottom edge cell
	if c.row == b.size-1 {
		// board's bottom left cell
		if c.col == 0 {
			cel = b.cells[b.size-1][1]
			if !cel.opened {
				op(b, cel)
			}
			cel = b.cells[b.size-2][0]
			if !cel.opened {
				op(b, cel)
			}
			cel = b.cells[b.size-2][1]
			if !cel.opened {
				op(b, cel)
			}
			return
		}
		// board's bottom right cell
		if c.col == b.size-1 {
			cel = b.cells[b.size-1][b.size-2]
			if !cel.opened {
				op(b, cel)
			}
			cel = b.cells[b.size-2][b.size-1]
			if !cel.opened {
				op(b, cel)
			}
			cel = b.cells[b.size-2][b.size-2]
			if !cel.opened {
				op(b, cel)
			}
			return
		}
		// border's bottom middle cell
		cel = b.cells[b.size-1][c.col-1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[b.size-1][c.col+1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[b.size-2][c.col-1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[b.size-2][c.col]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[b.size-2][c.col+1]
		if !cel.opened {
			op(b, cel)
		}
		return
	}
	// board's left edge cell
	if c.col == 0 {
		cel = b.cells[c.row][c.col+1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[c.row-1][c.col]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[c.row-1][c.col+1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[c.row+1][c.col]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[c.row+1][c.col+1]
		if !cel.opened {
			op(b, cel)
		}
		return
	}
	// board's right edge cell
	if c.col == b.size-1 {
		cel = b.cells[c.row][c.col-1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[c.row-1][c.col]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[c.row-1][c.col-1]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[c.row+1][c.col]
		if !cel.opened {
			op(b, cel)
		}
		cel = b.cells[c.row+1][c.col-1]
		if !cel.opened {
			op(b, cel)
		}
		return
	}
	// board's middle cells
	cel = b.cells[c.row-1][c.col-1]
	if !cel.opened {
		op(b, cel)
	}
	cel = b.cells[c.row-1][c.col]
	if !cel.opened {
		op(b, cel)
	}
	cel = b.cells[c.row-1][c.col+1]
	if !cel.opened {
		op(b, cel)
	}
	cel = b.cells[c.row][c.col-1]
	if !cel.opened {
		op(b, cel)
	}
	cel = b.cells[c.row][c.col+1]
	if !cel.opened {
		op(b, cel)
	}
	cel = b.cells[c.row+1][c.col-1]
	if !cel.opened {
		op(b, cel)
	}
	cel = b.cells[c.row+1][c.col]
	if !cel.opened {
		op(b, cel)
	}
	cel = b.cells[c.row+1][c.col+1]
	if !cel.opened {
		op(b, cel)
	}
}

// Prepares complete game board in its initial (all cells are closed) state.
func createBoard(size, bhNum int) (*board, error) {
	if size < minBoardSize {
		return nil, fmt.Errorf("size %dx%d is too small for a board", size, size)
	}
	if size > maxBoardSize {
		return nil, fmt.Errorf("size %dx%d is too big for a board", size, size)
	}
	if bhNum < minBlackHolesNum(size) {
		return nil, fmt.Errorf("%d black holes is not enough for %dx%d board", bhNum, size, size)
	}
	if bhNum > maxBlackHolesNum(size) {
		return nil, fmt.Errorf("%d black holes is too much for %dx%d board", bhNum, size, size)
	}

	// first, let's prepare a board
	rows := make([][]*cell, size)
	for rIdx := range rows {
		row := make([]*cell, size)
		for cIdx := range row {
			row[cIdx] = createCell(0, rIdx, cIdx)
		}
		rows[rIdx] = row
	}
	brd := &board{
		size:       size,
		cells:      rows,
		blackHoles: make([]*cell, bhNum),
	}

	// ... then, let's place black holes
	cellsNum := size * size
	maxCellIdx := cellsNum - 1

	rs := rand.NewSource(time.Now().UnixNano())
	rg := rand.New(rs)
	for placedBH := 0; placedBH != bhNum; {
		idx := rg.Intn(maxCellIdx)
		row := idx / size
		col := idx % size
		cel := brd.cells[row][col]
		if cel.markBlackHole() {
			brd.blackHoles[placedBH] = cel
			placedBH++
		}
	}

	// now let's update cells with correct codes
	for _, blackHoleCell := range brd.blackHoles {
		doWithNeighboringCells(brd, blackHoleCell, incrementNeighboringBlackHolesOp)
	}
	return brd, nil
}

// Opens cell by its address and returns true if cell is a black hole.
func (b *board) openCell(rIdx, cIdx int) (bool, error) {
	if rIdx >= b.size || cIdx >= b.size {
		return false, fmt.Errorf("cell %d:%d does not exist", rIdx, cIdx)
	}
	row := b.cells[rIdx]
	cel := row[cIdx]
	res := cel.open()
	if res == wasOpened {
		return false, nil
	}
	b.opened++
	switch res {
	case blackHole:
		return true, nil
	case hasNoNeighboringBlackHoles:
		// Here is the hardest part.
		// Is cell has no neighboring black holes - we should open all it's neighbourghs for better UX.
		doWithNeighboringCells(b, cel, recursivelyOpenSafeNeighbors)
	case hasNeighboringBlackHoles:
		return false, nil
	}
	return false, nil
}

const (
	cellOpening           = " "
	cellLongOpening       = "  "
	cellHeaderOpening     = " "
	cellHeaderLongOpening = "  "
	cellClosing           = "  |"
	cellHeaderClosing     = " ||"
	rowSepFragment        = '_'
)

// Makes proper cell presentation and writes it to a string builder.
// TODO: For more complex code this functionality could be placed in a separate unit or so-called presentation/view layer.
//
// Usually bool args (header bool) sucks, but this function is for internal use only, so...
func alignAndWriteCell(c string, header bool, sb *strings.Builder) {
	l := len(c)
	if header {
		switch l {
		case 1:
			sb.WriteString(cellHeaderLongOpening)
		case 2:
			sb.WriteString(cellHeaderOpening)
		default:
			// can't happen
			panic(fmt.Errorf("wrong len of cell content, expected: len to be 1 or 2, got: %d", l))
		}
		sb.WriteString(c)

		sb.WriteString(cellHeaderClosing)
		return
	}
	switch l {
	case 1:
		sb.WriteString(cellLongOpening)
	case 2:
		sb.WriteString(cellOpening)
	default:
		// can't happen
		panic(fmt.Errorf("wrong len of cell content, expected: len to be 1 or 2, got: %d", l))
	}
	sb.WriteString(c)
	sb.WriteString(cellClosing)
}

// Presents board as a string for CLI.
// TODO: For more complex code this functionality could be placed in a separate unit or so-called presentation/view layer.
func (b *board) toString() string {
	var sb strings.Builder
	alignAndWriteCell("\\", true, &sb)
	for i := 0; i < b.size; i++ {
		alignAndWriteCell(strconv.Itoa(i), true, &sb)
	}
	sb.WriteRune('\n')

	rowSepLen := sb.Len()
	for i := 0; i < rowSepLen; i++ {
		sb.WriteRune(rowSepFragment)
	}
	sb.WriteRune('\n')
	for i, row := range b.cells {
		alignAndWriteCell(strconv.Itoa(i), true, &sb)
		for _, cel := range row {
			alignAndWriteCell(cel.toString(), false, &sb)
		}
		sb.WriteRune('\n')
		for i := 0; i < rowSepLen; i++ {
			sb.WriteRune(rowSepFragment)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

// Opens all the cells. We use it to present a board at the end of game when a player lose, so no need for *Board.opened incrementation.
func (b *board) openAll() {
	for _, row := range b.cells {
		for _, cel := range row {
			cel.open()
		}
	}
}
func (b *board) checkWin() bool {
	return b.opened+len(b.blackHoles) == b.size*b.size
}
