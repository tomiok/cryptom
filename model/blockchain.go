package model

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
