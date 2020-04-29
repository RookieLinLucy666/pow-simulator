package main

import (
	"crypto/sha256"
	"encoding/binary"
)

type Ticket struct {
	Attempts, BlockNumber, MinerId int
	Nonce                          uint32
	Hash                           [32]byte
}

type Validator struct {
	BlockNumber int
	Difficulty  int
	WaitChan    chan Ticket
}

func NewValidator(difficulty int) Validator {
	v := Validator{}
	v.BlockNumber = 0
	v.Difficulty = difficulty
	v.WaitChan = make(chan Ticket)
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
		Nonce:       nonce}
	v.WaitChan <- newTicket
}

func (v *Validator) Validate(ticket Ticket) bool {
	if ticket.BlockNumber != v.BlockNumber {
		return false
	}

	if v.CheckHash(ticket.Hash) {
		v.BlockNumber++
		// v.Difficulty++
		return true
	}

	return false
}
