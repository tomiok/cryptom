package db

import (
	"cryptom/model"
	"github.com/boltdb/bolt"
)

const BlockChainFile = "blockchain.db"

/*
In blocks, the key -> value pairs are:

'b' + 32-byte block hash -> block index record
'f' + 4-byte file number -> file information record
'l' -> 4-byte file number: the last block file number used
'R' -> 1-byte boolean: whether we're in the process of reindexing
'F' + 1-byte flag name length + flag name string -> 1 byte boolean: various flags that can be on or off
't' + 32-byte transaction hash -> transaction index record
In chainstate, the key -> value pairs are:

'c' + 32-byte transaction hash -> unspent transaction output record for that transaction
'B' -> 32-byte block hash: the block hash up to which the database represents the unspent transaction outputs
*/

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

func AddBlock(data string, db *bolt.DB, bc *model.BlockChain) error {
	var lastHash []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockChainFile))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	newBlock := model.NewBlock(data, lastHash)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockChainFile))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.Tip = newBlock.Hash

		return err
	})

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
