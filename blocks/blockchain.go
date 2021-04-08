package blocks

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
)

const (
	blockChainFile = "chain.db"
	bucketBlocks   = "blocks"
)

type Blockchain struct {
	Tip []byte
	DB  BCDB
}

func (bc *Blockchain) Iterator() *BChainIterator {
	return &BChainIterator{bc.Tip, bc.DB}
}

// UpdateBlockchain creates a new Blockchain with genesis Block
func UpdateBlockchain(address string) *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(blockChainFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketBlocks))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{
		Tip: tip,
		DB:  &InMemoryBCDB{DB: db},
	}
}

func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Tx {
	var unspentTXs []Tx
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinBase() == false {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXO(pubHasKey []byte) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := bc.FindUnspentTransactions(pubHasKey)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubHasKey) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (bc *Blockchain) FindSpendableOutputs(pubHashKey []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubHashKey)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubHashKey) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Tx) {
	var lastHash []byte

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	bc.DB.ViewDB()

	newBlock := NewBlock("some data", transactions, lastHash)

	bc.DB.UpdateDB(newBlock)
}

func (bc *Blockchain) Clean() {
	bc.DB.CleanupDB()
}

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Tx) bool {
	prevTXs := make(map[string]Tx)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.TxID)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func dbExists() bool {
	if _, err := os.Stat(blockChainFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Tx, privateKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Tx)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.TxID)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privateKey, prevTXs)
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte

	cbtx := NewCoinBaseTx(address, genesisBaseData)
	genesis := NewGenesisBlock(cbtx)

	db, err := bolt.Open(blockChainFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(bucketBlocks))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &Blockchain{Tip: tip, DB: &InMemoryBCDB{DB: db}}
}
