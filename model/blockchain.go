package model

import (
	"cryptom/db"
)

type BlockChain struct {
	tip []byte
}

func (bc *BlockChain) NewBlockChain() *BlockChain {
	var tip []byte
	db = db.OpenForBc()
}
