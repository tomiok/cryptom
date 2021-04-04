package internal

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
	"log"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func GenerateID() string {
	return uuid.NewString()
}
