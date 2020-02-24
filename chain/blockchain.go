package chain

import (
	"cryptom/model"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
)

const (
	BlockChainFile = "chain.db"
	BlocksBucket   = "blocks"
)

type Blockchain struct {
	Tip []byte
	db  BCDB
}

func (bc *Blockchain) AddBlock(data string) {
	lastHash := bc.db.view()
	newBlock := model.NewBlock(data, lastHash)
	bc.Tip = bc.db.update(newBlock)
}

func (bc *Blockchain) Iterator() *BChainIterator {
	return &BChainIterator{bc.Tip, bc.db}
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(BlockChainFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := model.NewGenesis()

			b, err := tx.CreateBucket([]byte(BlocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, &InMemoryBCDB{db}}
}

func (bc *Blockchain) PrintChain() {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := model.NewPow(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
