package main

import (
	"fmt"
	"math/rand"
	"time"
)

const STARTING_DIFF = 3
const NUM_MINERS = 5

var miners [NUM_MINERS]Miner

func startNewRound(v *Validator) {
	fmt.Printf(
		"Starting round %d (Difficulty: %d)\n", v.BlockNumber, v.Difficulty)
}

func printWinner(ticket Ticket) {
	fmt.Println("New Block Found!")
	fmt.Printf(
		"Block %d was found by miner %d with nonce %v after %d attempts\n",
		ticket.BlockNumber, ticket.MinerId, ticket.Nonce, ticket.Attempts)
	fmt.Print(fmt.Sprintf("Hash: 0x%x\n", ticket.Hash))
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// Initalize validator
	v := NewValidator(STARTING_DIFF)
	// Initialize and start miners
	for i := 0; i < NUM_MINERS; i++ {
		miners[i] = Miner(i)
		go miners[i].Start(&v)
	}

	// Start the first round
	startNewRound(&v)
	for {
		select {
		case ticket := <-v.WaitChan:
			if v.Validate(ticket) {
				printWinner(ticket)
				startNewRound(&v)
			}
		default:
			fmt.Println(time.Now().Unix())
			time.Sleep(1 * time.Second)
		}
	}
}
