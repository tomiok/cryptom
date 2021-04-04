package blocks

type BChainIterator struct {
	CurrentHash []byte
	database    BCDB
}

func (iterator *BChainIterator) Next() *Block {
	var block *Block

	bytes := iterator.database.ViewIterator(iterator)

	block = Deserialize(bytes)

	iterator.CurrentHash = block.PrevBlockHash

	return block
}
