package model

import (
	"bytes"
	"cryptom/utils"
	"math"
	"math/big"
	"time"
)

const targetBits = 16

const maxNonce = math.MaxInt64

type ProofOfWork struct {
	Block  *Block
	target *big.Int
}

func NewPow(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{block, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			utils.IntToHex(pow.Block.Timestamp),
			utils.IntToHex(int64(targetBits)),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}
