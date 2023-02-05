package main

import (
	"flag"
	"fmt"

	"github.com/egoholic/proxx/proxx"
)

func main() {
	var size, bhNum int
	flag.IntVar(&size, "size", 15, "size of the game board")
	flag.IntVar(&bhNum, "bhNum", 15, "number of black holes")
	flag.Parse()
	game, err := proxx.Create(size, bhNum)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	game.Run()
}
