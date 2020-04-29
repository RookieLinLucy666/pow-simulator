package main

import (
	"fmt"
	"math/rand"
	"time"
)

// GenesisDifficulty - the difficulty of the genesis block
const GenesisDifficulty = 3

// NumMiners - number of miners mining for new blocks
const NumMiners = 5

var miners [NumMiners]Miner

func startNewRound(blockNumber int, diff int) {
	fmt.Printf(
		"Starting round %d (Difficulty: %d)\n", blockNumber, diff)
}

func printWinner(ticket Ticket) {
	fmt.Println("New Block Found!")
	fmt.Printf(
		"Block %d was found by miner %d with nonce %v after %d attempts\n",
		ticket.BlockNumber, ticket.MinerID, ticket.Nonce, ticket.Attempts)
	fmt.Print(fmt.Sprintf("Hash: 0x%x\n", ticket.Hash))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// Initalize validator
	validator := NewValidator(GenesisDifficulty)
	// Initialize and start miners
	for i := 0; i < NumMiners; i++ {
		miners[i] = Miner(i)
		go miners[i].Start(validator)
	}

	// Start the first round
	startNewRound(validator.CurrentBlockNumber(), validator.CurrentDifficulty())
	for {
		select {
		case ticket := <-validator.WaitChan:
			if validator.Validate(ticket) {
				printWinner(ticket)
				startNewRound(validator.CurrentBlockNumber(), validator.CurrentDifficulty())
			}
		default:
			if time.Now().Unix()%11 == 0 {
				validator.PrintLatestBlocks()
			}

			fmt.Println(time.Now().Unix())
			time.Sleep(1 * time.Second)
		}
	}
}
