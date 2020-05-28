package model

import (
	"bytes"
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
		identifier.String(),
		[]byte(data),
		transactions,
		[]byte{},
		prevBlockHash,
		time.Now().Unix(),
		Header{},
		0}

	pow := NewPow(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesis(coinbase transaction.Tx) *Block {
	fmt.Println("Creating the GENESIS block")
	return NewBlock("", []transaction.Tx{coinbase}, []byte{})
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
