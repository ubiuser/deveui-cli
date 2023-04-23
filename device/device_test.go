package device

import (
	"testing"
)

func TestCanGenerateValidCode(t *testing.T) {
	allowedChars := []string{"A", "B", "C", "D", "E", "F", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	device := NewDevice()

	if device == nil {
		t.Errorf("deivce should not be nil, but is %s", device)
	}

	if len(device.Code) != 5 {
		t.Errorf("code should  be 5 characters long, but is %d", len(device.Code))
	}

	hasChar := false
	for _, char := range allowedChars {
		if char == string(device.Code[0]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("first char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(device.Code[0]))
	}

	hasChar = false
	for _, char := range allowedChars {
		if char == string(device.Code[1]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("second char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(device.Code[1]))
	}

	hasChar = false
	for _, char := range allowedChars {
		if char == string(device.Code[2]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("third char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(device.Code[2]))
	}

	hasChar = false
	for _, char := range allowedChars {
		if char == string(device.Code[3]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("fourth char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(device.Code[3]))
	}

	hasChar = false
	for _, char := range allowedChars {
		if char == string(device.Code[4]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("fifth char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(device.Code[4]))
	}
}
