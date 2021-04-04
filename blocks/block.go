package blocks

import (
	"bytes"
	"crypto/sha256"
	"cryptom/internal"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	Identifier    string
	Data          []byte
	Transactions  []*Tx
	Hash          []byte
	PrevBlockHash []byte
	Header        Header
	Nonce         int
}

// NewBlock is the function that creates a new block and add it into the chain
func NewBlock(data string, transactions []*Tx, prevBlockHash []byte) *Block {
	block := &Block{
		Identifier:    internal.GenerateID(),
		Data:          []byte(data),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
	}

	pow := NewPow(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesis(base *Tx) *Block {
	fmt.Println("Creating the GENESIS block")
	return NewBlock("GENESIS", []*Tx{base}, []byte{})
}

// Serialize transform the block's data to slice of bytes
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)

	if err != nil {
		log.Println("cannot serialize " + err.Error())
		return nil
	}

	return res.Bytes()
}

// Deserialize transform an slice of bytes into a block
func Deserialize(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)

	if err != nil {
		log.Println("cannot deserialize " + err.Error())
		return nil
	}
	return &block
}

// hash all the transactions in the block
// TODO implement a merkle tree for store transaction hashes
func (b *Block) HashTransactions() []byte {
	var (
		hashes [][]byte
		hash   [32]byte
	)

	for _, tx := range b.Transactions {
		hashes = append(hashes, tx.Serialize())
	}
	hash = sha256.Sum256(bytes.Join(hashes, []byte{}))
	return hash[:]
}
