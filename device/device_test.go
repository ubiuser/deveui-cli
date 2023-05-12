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

	identifier := device.GetIdentifier()
	code := device.GetCode()

	if len(code) != 5 {
		t.Errorf("code should  be 5 characters long, but is %d", len(code))
	}

	if identifier[len(identifier)-5:] != code {
		t.Errorf("code should be last 5 characters of identifier, but is %s", code)
	}

	hasChar := false
	for _, char := range allowedChars {
		if char == string(code[0]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("first char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(code[0]))
	}

	hasChar = false
	for _, char := range allowedChars {
		if char == string(code[1]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("second char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(code[1]))
	}

	hasChar = false
	for _, char := range allowedChars {
		if char == string(code[2]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("third char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(code[2]))
	}

	hasChar = false
	for _, char := range allowedChars {
		if char == string(code[3]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("fourth char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(code[3]))
	}

	hasChar = false
	for _, char := range allowedChars {
		if char == string(code[4]) {
			hasChar = true
		} else if hasChar {
			break
		}
	}

	if hasChar == false {
		t.Errorf("fifth char should be A, B, C, D, E, F, 0, 1, 2, 3, 4, 5, 6, 7, 8 or 9: but is: %s", string(code[4]))
	}
}
