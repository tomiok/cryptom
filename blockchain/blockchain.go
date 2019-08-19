package blockchain

import (
	"cryptom/db"
)

type BlockChain struct {
	Tip []byte
}

func (bc *BlockChain) NewBlockChain() {
	dbInstance := db.OpenForBc()
	err := db.Save(dbInstance)

	if err != nil {
		panic(err)
	}
}

func (bc *BlockChain) AddBlock(data string) error {
	dbInstance := db.OpenForBc()
	err := db.AddBlock(data, dbInstance, bc)

	return err
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{CurrentHash: bc.Tip}
}
