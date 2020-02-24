package main

import "cryptom/chain"

func main() {
	bc := chain.NewBlockchain()

	bc.AddBlock("*****")
	bc.AddBlock("#####")

	bc.PrintChain()

	bc.Clean()
}
