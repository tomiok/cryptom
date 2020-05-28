package model

import (
	"bytes"
	"crypto/sha256"
	"cryptom/transaction"
	"encoding/gob"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Block struct {
	Identifier    string
	Data          []byte
	Transactions  []transaction.Tx
	Hash          []byte
	PrevBlockHash []byte
	Timestamp     int64
	Header        Header
	Nonce         int
}

// NewBlock is the function that creates a new block and add it into the chain
func NewBlock(data string, transactions []transaction.Tx, prevBlockHash []byte) *Block {
	identifier, _ := uuid.NewUUID()

	block := &Block{
		Identifier:    identifier.String(),
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

func NewGenesis(coinBase transaction.Tx) *Block {
	fmt.Println("Creating the GENESIS block")
	return NewBlock("", []transaction.Tx{coinBase}, []byte{})
}

// Serialize transform the block's data to slice of bytes
func (block *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(block)

	if err != nil {
		panic(err)
	}

	return res.Bytes()
}

// Deserialize transform an slice of bytes into a block
func Deserialize(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)

	if err != nil {
		panic(err)
	}

	return &block
}

// hash all the transactions in the block
func (b *Block) HashTransactions() []byte {
	var (
		hashes [][]byte
		hash   [32]byte
	)

	for _, tx := range b.Transactions {
		hashes = append(hashes, tx.ID)
	}
	hash = sha256.Sum256(bytes.Join(hashes, []byte{}))
	return hash[:]
}
