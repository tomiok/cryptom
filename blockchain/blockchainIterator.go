package blockchain

import (
	"cryptom/db"
	"cryptom/model"
	"github.com/boltdb/bolt"
)

type BlockChainIterator struct {
	CurrentHash []byte
}

func (iterator *BlockChainIterator) Next(database *bolt.DB) *model.Block {

	var block *model.Block

	err := database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.BlockChainFile))
		encodedBlock := b.Get(iterator.CurrentHash)
		block = model.Deserialize(encodedBlock)

		return nil
	})

	if err != nil {
		panic(err)
	}
	iterator.CurrentHash = block.PrevBlockHash

	return block
}
