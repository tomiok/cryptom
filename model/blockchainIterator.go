package model

import (
	"cryptom/db"
	"github.com/boltdb/bolt"
)

type BlockChainIterator struct {
	currentHash []byte
}

func (iterator *BlockChainIterator)Next(database *bolt.DB) *Block {

	var block *Block

	err := database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.BlockChainFile))
		encodedBlock := b.Get(iterator.currentHash)
		block = Desearlize(encodedBlock)

		return nil
	})

	if err != nil {
		panic(err)
	}
	iterator.currentHash = block.PrevBlockHash

	return block
}
