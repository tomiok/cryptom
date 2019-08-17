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

func Save(db *bolt.DB) error {

	err := db.Update(update)

	if err != nil {
		panic(err)
	}

	return err
}

func update(tx *bolt.Tx) error {
	var tip []byte
	return doUpdate(tx, tip)
}

func doUpdate(tx *bolt.Tx, tip []byte) error {
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
