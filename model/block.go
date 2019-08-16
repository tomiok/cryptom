package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Block struct {
	Identifier    string
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
	Timestamp     int64
	Header        Header
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	identifier, _ := uuid.NewUUID()

	newBlock := &Block{identifier.String(), []byte(data), []byte{}, prevBlockHash, time.Now().Unix(),
		Header{}, 0}

	pow := NewPow(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce
	return newBlock
}

func NewGenesis() *Block {
	fmt.Println("Creating the GENESIS block")
	var block Block
	return &block
}

func (block *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(block)

	if err != nil {
		panic(err)
	}

	return res.Bytes()
}

func desearlize(b []byte) *Block {

	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)

	if err != nil {
		panic(err)
	}

	return &block
}
