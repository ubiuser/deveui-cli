package device

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	allowedChars = "ABCDEF0123456789" // accepted chars used to make up DevEUI
	devEuiLength = 16                 // valid DevEUI is string of length 16
)

type Device struct {
	identifier string
	code       string
}

// NewDevice Build a new device with DevEUI identifier and code values.
// Example: 1CEB0080F074F750 4F750
func NewDevice() (*Device, error) {
	hex, err := generateHexString()
	if err != nil {
		return nil, fmt.Errorf("failed to generate DevEUI: %w", err)
	}

	return &Device{
		identifier: hex,
		code:       hex[len(hex)-5:],
	}, nil
}

func (d Device) GetIdentifier() string {
	return d.identifier
}

func (d Device) GetCode() string {
	return d.code
}

// String returns a string representation of the device (see Stringer interface at https://go.dev/tour/methods/17)
func (d Device) String() string {
	return fmt.Sprintf("device has identifier: %s and code: %s", d.identifier, d.code)
}

// Generate valid DevEUI identifier value.
// Example:	1CEB0080F074F750
func generateHexString() (string, error) {
	max := big.NewInt(int64(len(allowedChars)))
	b := make([]byte, devEuiLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("failed to generate random int: %w", err)
		}
		b[i] = allowedChars[n.Int64()]
	}
	return string(b), nil
}
