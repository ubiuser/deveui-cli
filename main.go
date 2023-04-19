package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

const allowedChars = "ABCDEF0123456789"

func main() {
	hexStr, err := generateHexString(16)

	if err != nil {
		log.Print(err)
	}

	fmt.Print(hexStr)
}

func generateHexString(length int) (string, error) {
	max := big.NewInt(int64(len(allowedChars)))
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = allowedChars[n.Int64()]
	}
	return string(b), nil
}
