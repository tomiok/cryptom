package main

import (
	"cryptom/blocks"
)

func main() {
	bc := blocks.NewBlockchain("some address")
	bc.PrintChain()

	// delete
	bc.Clean()
}
