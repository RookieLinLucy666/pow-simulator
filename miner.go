package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"
)

const DIFFICULTY = 3

func printHash(hash [32]byte) {
	for i := 0; i < 32; i++ {
		fmt.Printf("%02x", hash[i])
	}
}

func successful(hash [32]byte) bool {
	for i := 0; i < DIFFICULTY; i++ {
		if hash[i] != 0 {
			return false
		}
	}
	return true
}

func hashAndPrint(nonce uint32) {
	attempts := 0
	bytebuffer := make([]byte, 32)

	for {
		attempts++
		x := rand.Uint32()
		sum := nonce + x
		binary.BigEndian.PutUint32(bytebuffer, sum)
		hash := sha256.Sum256(bytebuffer)
		if successful(hash) {
			fmt.Printf(
				"Found summand to be %v, since sha256(%v + %v) = ", x, nonce, x)
			printHash(hash)
			fmt.Printf(", after %d attempts.\n", attempts)
			break
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	nonce := rand.Uint32()
	fmt.Printf("Mining nonce: %d\n", nonce)
	hashAndPrint(nonce)
}
