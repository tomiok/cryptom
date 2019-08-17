package model

import "github.com/boltdb/bolt"

type BlockChain struct {
	tip []byte
	db *bolt.DB
}