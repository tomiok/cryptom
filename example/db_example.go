package main

import "cryptom/chain"

func main() {
	bc := chain.NewBlockchain()
	bc.AddBlock("le doy un vergacoin al vergante")
	bc.AddBlock("le doy otro vergacoin!")

	bc.PrintChain()

	// delete
	bc.Clean()
}
