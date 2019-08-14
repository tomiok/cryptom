package model

import "time"

type Block struct {
	Identifier    string
	Data          []byte
	Hash          []byte
	PrevBlockHash []byte
	Timestamp     time.Time
	Header        Header
}
