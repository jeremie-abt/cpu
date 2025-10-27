package common

import "testing"

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input    []byte
		expected int64
	}{
		{[]byte("42"), 42},
		{[]byte("0"), 0},
		{[]byte("123456"), 123456},
		{[]byte("0x1A"), 26},
		{[]byte("0XFF"), 255},
		{[]byte("0x0"), 0},
		{[]byte("0b1010"), 10},
		{[]byte("0B1111"), 15},
		{[]byte("0b0"), 0},
		{[]byte("0x"), 0},
		{[]byte("0b"), 0},
		{[]byte("0b1000010010010010101000101011001011010"), 71174477402},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			got := ParseNumber(tt.input)
			if got != tt.expected {
				t.Errorf("ParseNumber(%q) = %d; want %d", tt.input, got, tt.expected)
			}
		})
	}
}
