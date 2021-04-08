package cli

import (
	"cryptom/blocks"
	"cryptom/internal"
	"fmt"
)

type CLI struct {
}

func (cli *CLI) Run() {
	w1 := createWallet()
	fmt.Printf("Your new address: %s\n", w1)

	w2 := createWallet()
	fmt.Printf("Your new address: %s\n", w2)

	createBlockchain(w1)
	getBalance(w1)
	send(w1, w2, 3)

	getBalance(w1)
	getBalance(w2)
}

func createBlockchain(address string) {
	blocks.CreateBlockchain(address)
}

func getBalance(address string) {
	valid := blocks.ValidateAddress(address)

	if !valid {
		panic("address in not valid")
	}

	bc := blocks.UpdateBlockchain(address)

	defer bc.DB.CloseDB()

	balance := 0
	pubHashKey := internal.Base58Decode([]byte(address))
	pubHashKey = pubHashKey[1 : len(pubHashKey)-4]
	UTXOs := bc.FindUTXO(pubHashKey)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func send(from, to string, amount int) {
	fmt.Print("sending... \n")
	bc := blocks.UpdateBlockchain(from)
	defer bc.DB.CloseDB()

	tx := blocks.NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*blocks.Tx{tx})
	fmt.Println("Success!")
}

func createWallet() string {
	wallets, _ := blocks.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	return address
}
