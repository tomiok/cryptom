package main

import (
	"cryptom/blocks"
	"fmt"
)

func main() {
	//need to move this to example package
	fmt.Println("Init ")
	genesis := blocks.NewGenesis(blocks.MakeCoinBaseTx("", "here we are"))
	b1 := blocks.NewBlock("some data", nil, genesis.PrevBlockHash)
	pow1 := blocks.NewPow(genesis)
	pow2 := blocks.NewPow(b1)
	fmt.Println(pow1.Validate())
	fmt.Println(pow2.Validate())
}
