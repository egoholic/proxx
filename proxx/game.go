package proxx

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Game struct {
	board *board
}

// Creates new game, but do not run it.
func Create(size, bhNum int) (*Game, error) {
	if size < 9 {
		return nil, fmt.Errorf("size %dx%d is too small for a board", size, size)
	}
	if bhNum > size*2 { // here I assume the max number of black holes to be no greater than square root of cells number * 2
		return nil, fmt.Errorf("%d blackholes is too much for %dx%d board", bhNum, size, size)
	}
	board, err := createBoard(size, bhNum)
	if err != nil {
		return nil, err
	}
	return &Game{
		board: board,
	}, nil
}

// Runs game loop.
func (g *Game) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(g.board.toString())
		if g.board.checkWin() {
			fmt.Println("You won!")
			return
		}
		fmt.Print("Enter row and column numbers of a cell to open in format: <row>:<column>: ")
		_rowCol, _ := reader.ReadString('\n')
		_rowCol = strings.TrimSuffix(_rowCol, "\n")
		rowCol := strings.Split(_rowCol, ":")
		if len(rowCol) != 2 {
			fmt.Println("Error: wrong row and column numbers. Please try again:")
			continue
		}
		rIdx, err := strconv.Atoi(rowCol[0])
		if err != nil {
			fmt.Printf("Error: wrong row number(%s). Please try again:\n", err.Error())
			continue
		}
		cIdx, err := strconv.Atoi(rowCol[1])
		if err != nil {
			fmt.Printf("Error: wrong column number(%s). Please try again:\n", err.Error())
			continue
		}
		isBlackHole, err := g.board.openCell(rIdx, cIdx)
		if err != nil {
			fmt.Print(g.board.toString())
			fmt.Printf("Error: %s. Please try again:\n", err.Error())
			continue
		}
		if isBlackHole {
			g.board.openAll()
			fmt.Print(g.board.toString())
			fmt.Printf("You lose!")
			return
		}
	}
}
