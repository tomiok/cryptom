package blocks

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"cryptom/internal"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
)

/**
Outputs are where “coins” are stored. Each output comes with an unlocking script, which determines the logic
of unlocking the output. Every new transaction must have at least one input and output. An input references an output
from a previous transaction and provides data (the ScriptSig field) that is used in the output’s unlocking script
to unlock it and use its value to create new outputs.
*/
const genesisBaseData = "we fight for a better world"

type Tx struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// The input of the transaction
type TxInput struct {
	TxID      []byte
	Vout      int
	Signature []byte //should have the signature from the wallet
	PubKey    []byte
}

type TxOutput struct {
	Value      int
	PubHashKey []byte
}

func (ti *TxInput) UsesKey(pubKey []byte) bool {
	lockingHash := HashPubKey(ti.PubKey)
	return bytes.Compare(lockingHash, pubKey) == 0
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

// NewCoinBaseTx is the "egg" for the transactions. Is the beginning of the transaction history.
func NewCoinBaseTx(to, data string) *Tx {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	in := TxInput{
		TxID:   []byte{},
		Vout:   -1,
		PubKey: []byte(data),
	}
	out := NewTXOutput(10, to)

	tx := Tx{
		Vin:  []TxInput{in},
		Vout: []TxOutput{*out},
	}
	tx.ID = tx.Hash()

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

// Hash returns the hash of the Transaction
func (tx *Tx) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Tx {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

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
			input := TxInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := Tx{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (Tx, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Tx{}, errors.New("transaction is not found")
}

// Verify verifies signatures of Transaction inputs
func (tx *Tx) Verify(prevTXs map[string]Tx) bool {
	if tx.IsCoinBase() {
		return true
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.TxID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.TxID)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubHashKey
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx *Tx) TrimmedCopy() Tx {
	var inputs []TxInput
	var outputs []TxOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TxInput{vin.TxID, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TxOutput{vout.Value, vout.PubHashKey})
	}

	txCopy := Tx{tx.ID, inputs, outputs}

	return txCopy
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (to *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(to.PubHashKey, pubKeyHash) == 0
}

// Lock signs the output
func (to *TxOutput) Lock(address []byte) {
	pubHashKey := internal.Base58Decode(address)
	pubHashKey = pubHashKey[1 : len(pubHashKey)-4]
	to.PubHashKey = pubHashKey
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

// Sign signs each input of a Transaction
func (tx *Tx) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Tx) {
	if tx.IsCoinBase() {
		return
	}

	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.TxID)].ID == nil {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.TxID)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubHashKey
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}
}
