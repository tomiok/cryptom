package model

import (
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
	identifier,_ :=uuid.NewUUID()

	newBlock := &Block{identifier.String(), []byte(data), []byte{}, prevBlockHash,time.Now().Unix(),
		Header{}, 0}

	pow := NewPow(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce
	return newBlock
}
