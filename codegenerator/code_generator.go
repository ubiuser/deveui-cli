package codegenerator

import (
	"crypto/rand"
	"math/big"
)

const (
	ALLOWED_CHARS  = "ABCDEF0123456789" // accepted chars used to make up DevEUI
	DEV_EUI_LENGTH = 16                 // valid DevEUI is string of length 16
)

// Generate valid DevEUI code from DevEUI identifier.
//
// # Example
//
//	1CEB0080F074F750 will return 4F750
func GenerateCode(hex string) (string, error) {
	return hex[len(hex)-5:], nil
}

// Generate valid DevEUI identifier value.
//
// # Example
//
//	1CEB0080F074F750
func GenerateHexString() (string, error) {
	max := big.NewInt(int64(len(ALLOWED_CHARS)))
	b := make([]byte, DEV_EUI_LENGTH)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = ALLOWED_CHARS[n.Int64()]
	}
	return string(b), nil
}
