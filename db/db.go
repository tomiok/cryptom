package db

import (
	"cryptom/model"
	"github.com/boltdb/bolt"
)

const BlockChainFile = "blockchain.db"

func OpenForBc() *bolt.DB {
	db, _ := bolt.Open(BlockChainFile, 0600, nil)
	return db
}

func save() {
	db, err := bolt.Open(BlockChainFile, 0600, nil)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err := db.Update(update)


}

func update(tx *bolt.Tx) error {
	b := tx.Bucket([]byte(BlockChainFile))

	if b == nil {
		genesis := model.NewGenesis()
		b, err := tx.CreateBucket([]byte(BlockChainFile))

		if err != nil {
			panic(err)
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		err = b.Put([]byte("l"), genesis.Hash)
		tip = genesis.Hash
	} else {
		tip = b.Get([]byte("l"))
	}

	return nil
}




