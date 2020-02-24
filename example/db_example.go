package main

import "cryptom/chain"

func main() {
	bc := chain.NewBlockchain()
	bc.AddBlock("*****")
	bc.AddBlock("#####")

	bc.PrintChain()

	// delete
	bc.Clean()
}