package blocks

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

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

type BCDB interface {
	UpdateDB(block *Block) []byte
	ViewDB() []byte
	ViewIterator(iterator *BChainIterator) []byte
	CleanupDB()
	CloseDB()
	ViewChain(func(tx *bolt.Tx) error)
}

type InMemoryBCDB struct {
	*bolt.DB
}

func (i *InMemoryBCDB) CloseDB() {
	i.Close()
}

func (i *InMemoryBCDB) ViewChain(fn func(tx *bolt.Tx) error) {
	db, err := bolt.Open(blockChainFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.View(fn)
}

func (i *InMemoryBCDB) UpdateDB(block *Block) []byte {
	db, _ := bolt.Open(blockChainFile, 0600, nil)
	defer db.Close()
	var tip []byte
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketBlocks))
		err := b.Put(block.Hash, block.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			log.Panic(err)
		}

		tip = block.Hash

		return nil
	})
	return tip
}

func (i *InMemoryBCDB) ViewDB() []byte {
	db, _ := bolt.Open(blockChainFile, 0600, nil)
	defer db.Close()
	var lastHash []byte
	i.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketBlocks))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	return lastHash
}

func (i *InMemoryBCDB) ViewIterator(iterator *BChainIterator) []byte {
	db, _ := bolt.Open(blockChainFile, 0600, nil)
	defer db.Close()
	var block []byte
	i.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketBlocks))
		block = b.Get(iterator.currentHash)

		return nil
	})

	return block
}

func (i *InMemoryBCDB) CleanupDB() {
	i.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketBlocks))
	})
}
