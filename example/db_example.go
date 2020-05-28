package main

import "cryptom/chain"

func main() {
	bc := chain.NewBlockchain("some address")
	bc.PrintChain()

	// delete
	bc.Clean()
}
