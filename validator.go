package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	blockchain "github.com/michaelhly/pow-simulator/blockchain"
)

type Ticket struct {
	Attempts, BlockNumber int
	MinerID               Miner
	Nonce                 uint32
	Hash                  [32]byte
	TicketTime            int64
}

type Validator struct {
	block         *blockchain.Block
	difficulty    int
	lastBlockTime int64
	WaitChan      chan Ticket
}

func NewValidator(startDiff int) *Validator {
	v := Validator{}
	v.block = blockchain.CreateBlock(startDiff)
	v.difficulty = startDiff
	v.lastBlockTime = 0
	v.WaitChan = make(chan Ticket)
	return &v
}

func (v *Validator) CheckHash(hash [32]byte) bool {
	for i := 0; i < v.difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}

func (v *Validator) Validate(ticket Ticket) bool {
	blockNumber := v.CurrentBlockNumber()
	if ticket.BlockNumber != blockNumber {
		return false
	}

	success := v.CheckHash(ticket.Hash)
	if success {
		if v.lastBlockTime > ticket.TicketTime {
			panic("Last block time came after this ticket. Should not create new block.")
		}

		v.lastBlockTime = ticket.TicketTime
		if (blockNumber+1)%1000 == 0 {
			v.difficulty++
		}
		v.block.ConfirmBlock(int(ticket.MinerID), ticket.TicketTime, ticket.Attempts)
		v.block = v.block.NextBlock(v.difficulty)
	}

	return success
}

func (v *Validator) AddTicket(m Miner, blockNumber int, nonce uint32, attempts int) {
	buffer := make([]byte, 32)
	binary.BigEndian.PutUint32(buffer, nonce)
	hash := sha256.Sum256(buffer)

	newTicket := Ticket{
		Attempts:    attempts,
		BlockNumber: blockNumber,
		Hash:        hash,
		MinerID:     m,
		Nonce:       nonce,
		TicketTime:  time.Now().UnixNano()}
	v.WaitChan <- newTicket
}

func (v *Validator) CurrentBlockNumber() int {
	return v.block.GetBlockHeight()
}

func (v *Validator) CurrentDifficulty() int {
	return v.difficulty
}

func (v *Validator) PrintLatestBlocks() {
	fmt.Println(
		"-----------------------------------Latest Blocks-----------------------------------")
	currBlock := v.block
	for currBlock != nil {
		fmt.Println(currBlock)
		currBlock = currBlock.GetPrevBlock()
	}
	fmt.Println(
		"-----------------------------------------------------------------------------------")
}
