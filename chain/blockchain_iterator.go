package chain

import (
	"cryptom/model"
	"github.com/boltdb/bolt"
)

type BChainIterator struct {
	CurrentHash []byte
	database    *bolt.DB
}

func (iterator *BChainIterator) Next() *model.Block {
	var block *model.Block

	err := iterator.database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
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
