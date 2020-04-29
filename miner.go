package main

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

type Miner int

func (m Miner) Start(v *Validator) {
	attempts := 0
	nonce := rand.Uint32()
	blockNum := v.CurrentBlockNumber()
	bytebuffer := make([]byte, 32)

	for {
		if blockNum < v.CurrentBlockNumber() {
			blockNum = v.CurrentBlockNumber()
			attempts = 0
			nonce = rand.Uint32()
		}

		attempts++
		x := rand.Uint32()
		sum := nonce + x
		binary.BigEndian.PutUint32(bytebuffer, sum)
		hash := sha256.Sum256(bytebuffer)
		if v.CheckHash(hash) {
			v.AddTicket(m, blockNum, sum, attempts)
		}
	}
}
