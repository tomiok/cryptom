package blocks

import (
	bolt "go.etcd.io/bbolt"
)

type BChainIterator struct {
	currentHash []byte
	db          BCDB
}

func (i *BChainIterator) Next() *Block {
	var block *Block

	i.db.ViewChain(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketBlocks))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	i.currentHash = block.PrevBlockHash
	return block
}
