package chain

import (
	"cryptom/model"
)

type BChainIterator struct {
	CurrentHash []byte
	database    BCDB
}

func (iterator *BChainIterator) Next() *model.Block {
	var block *model.Block

	bytes := iterator.database.viewIterator(iterator)

	block = model.Deserialize(bytes)

	iterator.CurrentHash = block.PrevBlockHash

	return block
}
