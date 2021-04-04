package cli

import (
	"cryptom/blocks"
	"fmt"
)

type CLI struct {
}

func (cli *CLI) GetBalance(address string) {
	bc := blocks.NewBlockchain(address)

	defer bc.DB.CloseDB()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) Send(from, to string, amount int) {
	bc := blocks.NewBlockchain(from)
	defer bc.DB.CloseDB()

	tx := blocks.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*blocks.Tx{tx})
	fmt.Println("Success!")
}
