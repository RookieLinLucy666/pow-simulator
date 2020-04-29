package main

import (
	"crypto/sha256"
	"encoding/binary"
	"sync"
	"time"
)

type Ticket struct {
	Attempts, BlockNumber, MinerId int
	Nonce                          uint32
	Hash                           [32]byte
	TicketTime                     int64
}

type Validator struct {
	BlockNumber, Difficulty int
	WaitChan                chan Ticket
	lastBlockTime           int64
	mux                     sync.Mutex
}

func NewValidator(difficulty int) Validator {
	v := Validator{}
	v.BlockNumber = 0
	v.Difficulty = difficulty
	v.WaitChan = make(chan Ticket)
	v.lastBlockTime = 0
	return v
}

func (v Validator) CheckHash(hash [32]byte) bool {
	for i := 0; i < v.Difficulty; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}

func (v Validator) AddTicket(minerId int, blockNumber int, nonce uint32, attempts int) {
	buffer := make([]byte, 32)
	binary.BigEndian.PutUint32(buffer, nonce)
	hash := sha256.Sum256(buffer)

	newTicket := Ticket{
		Attempts:    attempts,
		BlockNumber: blockNumber,
		Hash:        hash,
		MinerId:     minerId,
		Nonce:       nonce,
		TicketTime:  time.Now().UnixNano()}
	v.WaitChan <- newTicket
}

func (v *Validator) Validate(ticket Ticket) bool {
	// Validate only one ticket at a time to avoid race conditions
	v.mux.Lock()
	if ticket.BlockNumber != v.BlockNumber {
		v.mux.Unlock()
		return false
	}

	success := v.CheckHash(ticket.Hash)
	if success {
		if v.lastBlockTime > ticket.TicketTime {
			panic("Last block time came after this ticket. Should not create new block.")
		}
		v.lastBlockTime = ticket.TicketTime
		v.BlockNumber++
		if v.BlockNumber%1000 == 0 {
			v.Difficulty++
		}
	}

	v.mux.Unlock()
	return success
}
