package blocks

import (
	"bytes"
	"crypto/sha256"
	"cryptom/internal"
	"fmt"
	"math"
	"math/big"
)

const (
	targetBits = 16 // arbitrary number, 24 will work for staging or prod (bigger is more difficult)
	maxBits    = 256
	maxNonce   = math.MaxInt64
)

type ProofOfWork struct {
	Block  *Block
	target *big.Int
}

func NewPow(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(maxBits-targetBits))

	pow := &ProofOfWork{block, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			pow.Block.HashTransactions(),
			internal.IntToHex(pow.Block.Timestamp),
			internal.IntToHex(int64(targetBits)),
			internal.IntToHex(int64(nonce)),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int

	fmt.Printf("Mining the block containing \"%s\"\n", pow.Block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
