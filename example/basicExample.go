package main

import (
	"cryptom/model"
	"fmt"
)

func main() {
	//need to move this to example package
	fmt.Println("Init ")
	block := model.NewGenesis()
	model.NewBlock("some data", block.PrevBlockHash)
	pow := model.NewPow(block)
	pow.Validate()
}
