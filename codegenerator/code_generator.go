package codegenerator

import (
	"crypto/rand"
	"math/big"
)

const allowedChars = "ABCDEF0123456789"

func Generate(hex string) (string, error) {

	return hex[len(hex)-5:], nil
}

func GenerateHexString(length int) (string, error) {
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
