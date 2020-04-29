package blockchain

import "fmt"

type Block struct {
	confirmed                       bool
	height, diff, minerID, attempts int
	prev                            *Block
	time                            int64
}

func CreateBlock(diff int) *Block {
	block := Block{}
	block.confirmed = false
	block.height = 0
	block.diff = diff
	block.prev = nil
	block.minerID = -1
	block.time = -1
	return &block
}

func (b Block) NextBlock(newDiff int) *Block {
	if !b.confirmed {
		panic("Cannot create next block on an unconfirmed block.")
	}

	nextBlock := CreateBlock(newDiff)
	nextBlock.height = b.GetBlockHeight() + 1
	nextBlock.prev = &b
	return nextBlock
}

func (b *Block) ConfirmBlock(miner int, time int64, attempts int) {
	if b.confirmed {
		panic("Cannot confirm an already confirmed block.")
	}

	b.time = time
	b.minerID = miner
	b.attempts = attempts
	b.confirmed = true
}

func (b Block) GetBlockHeight() int {
	return b.height
}

func (b Block) GetPrevBlock() *Block {
	return b.prev
}

func (b Block) String() string {
	if !b.confirmed {
		return fmt.Sprintf(
			"[Height: %d, Confirmed: %v, Difficulty: %d]", b.height, b.confirmed, b.diff)
	}

	return fmt.Sprintf(
		"[Height: %d, Confirmed: %v, Time: %v, Difficulty: %d, Miner: %d, Attempts: %d]",
		b.height, b.confirmed, b.time, b.diff, b.minerID, b.attempts)
}
