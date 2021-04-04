package blocks

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
)

/**
Outputs are where “coins” are stored. Each output comes with an unlocking script, which determines the logic
of unlocking the output. Every new transaction must have at least one input and output. An input references an output
from a previous transaction and provides data (the ScriptSig field) that is used in the output’s unlocking script
to unlock it and use its value to create new outputs.
*/
const baseSignature = "we fight for a better world"

type Tx struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
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

func (ti *TxInput) CanUnlockOutputWith(unlocking string) bool {
	return ti.ScriptSignature == unlocking
}

func (to *TxOutput) CanBeUnlockedWith(unlocking string) bool {
	return to.ScriptPublicKey == unlocking
}

func (tx *Tx) SetID() {
	var (
		encoded bytes.Buffer
		hash    [32]byte
	)

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// MakeCoinBaseTx is the "egg" for the transactions. Is the beginning of the transaction history.
func MakeCoinBaseTx(to, data string) *Tx {
	if data == "" {
		data = baseSignature
	}

	txIn := TxInput{
		TxID:            []byte{},
		Vout:            -1,
		ScriptSignature: data,
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
	return &tx
}
func (tx *Tx) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TxID) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Tx) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(tx)

	if err != nil {
		panic(err)
	}

	return res.Bytes()
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Tx {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from}) // a change
	}

	tx := Tx{nil, inputs, outputs}
	tx.SetID()

	return &tx
}