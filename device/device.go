package device

import (
	"crypto/rand"
	"log"
	"math/big"
)

const (
	ALLOWED_CHARS  = "ABCDEF0123456789" // accepted chars used to make up DevEUI
	DEV_EUI_LENGTH = 16                 // valid DevEUI is string of length 16
)

type Device struct {
	Identifier string
	Code       string
}

// Build new device with DevEUI identifier and code values.
//
// # Example
//
//	1CEB0080F074F750 4F750
func NewDevice() *Device {
	hex, err := generateHexString()
	if err != nil {
		log.Print(err)
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
