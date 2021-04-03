package main

import (
	"cryptom/model"
	tx "cryptom/transaction"
	"fmt"
)

func main() {
	//need to move this to example package
	fmt.Println("Init ")
	block := model.NewGenesis(tx.MakeCoinBaseTx("", "here we are"))
	model.NewBlock("some data", block.PrevBlockHash)
	pow := model.NewPow(block)
	pow.Validate()
}
