package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

/**
Outputs are where “coins” are stored. Each output comes with an unlocking script, which determines the logic
of unlocking the output. Every new transaction must have at least one input and output. An input references an output
from a previous transaction and provides data (the ScriptSig field) that is used in the output’s unlocking script
to unlock it and use its value to create new outputs.
*/

type Tx struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

func (tx *Tx) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// The input of the transaction
type TxInput struct {
	TxID            []byte
	Vout            int
	ScriptSignature string //should have the signature from the wallet
}

type TxOutput struct {
	Value           int
	ScriptPublicKey string
}

// MakeCoinBaseTx is the "egg" for the transactions. Is the beginning of the transaction history.
func MakeCoinBaseTx(to, signature string) Tx {
	if signature == "" {
		signature = "some signature"
	}

	txIn := TxInput{
		TxID:            []byte{},
		Vout:            -1,
		ScriptSignature: signature,
	}

	txOut := TxOutput{
		Value:           999, //constant, for now
		ScriptPublicKey: to,
	}

	tx := Tx{
		ID:   nil,
		Vin:  []TxInput{txIn},
		Vout: []TxOutput{txOut},
	}
	return tx
}
