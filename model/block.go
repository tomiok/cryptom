package model

import (
	"bytes"
	"crypto/sha256"
	"cryptom/internal"
	"cryptom/transaction"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Timestamp     int64
	Identifier    string
	Data          []byte
	Transactions  []*transaction.Tx
	Hash          []byte
	PrevBlockHash []byte
	Header        Header
	Nonce         int
}

// NewBlock is the function that creates a new block and add it into the chain
func NewBlock(data string, transactions []*transaction.Tx, prevBlockHash []byte) *Block {
	block := &Block{
		Identifier:    internal.GenerateID(),
		Data:          []byte(data),
		Transactions:  transactions,
		Hash:          []byte{},
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().UnixNano(),
		Header:        Header{},
		Nonce:         0}

	pow := NewPow(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesis(base *transaction.Tx) *Block {
	fmt.Println("Creating the GENESIS block")
	return NewBlock("GENESIS", []*transaction.Tx{base}, []byte{})
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
