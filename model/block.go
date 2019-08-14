package model

import "time"

type Block struct {
	Identifier string
	Data       []byte
	Timestamp  time.Time
	Header     Header
}
