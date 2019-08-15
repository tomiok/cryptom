package model

import (
	"crypto/sha256"
	"time"
)

type Block struct {
	Identifier          string
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
	Timestamp     time.Time
	Header        Header
}

func GenerateHash(block *Block) string {
	record := block.Identifier + string(block.Hash) + block.Timestamp.String() + string(block.PrevBlockHash)
	hashFunc := sha256.New()
	hashFunc.Write([]byte(record))
}
