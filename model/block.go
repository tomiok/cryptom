package model

import (
	"crypto/sha256"
	"encoding/hex"
)

type Block struct {
	Identifier    string
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
	Timestamp     int64
	Header        Header
	Nonce         int64
}

func GenerateHash(block *Block) string {
	record := block.Identifier + string(block.Hash) + string(block.Timestamp) + string(block.PrevBlockHash)
	hashFunc := sha256.New()
	hashFunc.Write([]byte(record))
	hashed := hashFunc.Sum(nil)
	return hex.EncodeToString(hashed)
}
