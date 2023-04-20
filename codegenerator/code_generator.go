package codegenerator

import (
	"crypto/rand"
	"math/big"
)

const allowedChars = "ABCDEF0123456789"

func Generate() (string, error) {
	hexStr, err := generateHexString(16)
	if err != nil {
		return "", err
	}
	return hexStr[len(hexStr)-5:], nil
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
