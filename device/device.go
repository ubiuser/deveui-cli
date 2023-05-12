package device

import (
	"crypto/rand"
	"log"
	"math/big"
)

const (
	AllowedChars = "ABCDEF0123456789" // accepted chars used to make up DevEUI
	DevEuiLength = 16                 // valid DevEUI is string of length 16
)

type Device struct {
	Identifier string
	Code       string
}

// NewDevice Build a new device with DevEUI identifier and code values.
//
// # Example
//
//	1CEB0080F074F750 4F750
func NewDevice() *Device {
	hex, err := generateHexString()
	if err != nil {
		log.Fatal(err)
	}

	return &Device{
		Identifier: hex,
		Code:       hex[len(hex)-5:],
	}
}

// Generate valid DevEUI identifier value.
//
// # Example
//
//	1CEB0080F074F750
func generateHexString() (string, error) {
	max := big.NewInt(int64(len(AllowedChars)))
	b := make([]byte, DevEuiLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = AllowedChars[n.Int64()]
	}
	return string(b), nil
}
