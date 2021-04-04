package main

import (
	"cryptom/model"
	tx "cryptom/transaction"
	"fmt"
)

func main() {
	//need to move this to example package
	fmt.Println("Init ")
	genesis := model.NewGenesis(tx.MakeCoinBaseTx("", "here we are"))
	b1 := model.NewBlock("some data", nil, genesis.PrevBlockHash)
	pow1 := model.NewPow(genesis)
	pow2 := model.NewPow(b1)
	fmt.Println(pow1.Validate())
	fmt.Println(pow2.Validate())
}
