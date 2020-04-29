package blockchain

type Block struct {
	confirmed    bool
	number, diff int
	prev         *Block
}

func CreateBlock(diff int) *Block {
	block := Block{}
	block.confirmed = false
	block.number = 0
	block.diff = diff
	block.prev = nil
	return &block
}

func (b Block) NextBlock(newDiff int) *Block {
	if !b.confirmed {
		panic("Cannot create next block on an unconfirmed block.")
	}

	nextBlock := CreateBlock(newDiff)
	nextBlock.number = b.GetBlockNumber() + 1
	nextBlock.prev = &b
	return nextBlock
}

func (b *Block) ConfirmBlock() {
	if b.confirmed {
		panic("Cannot confirm an already confirmed block.")
	}

	b.confirmed = true
}

func (b Block) GetBlockNumber() int {
	return b.number
}
