package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWord_GreenPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantWord []byte
		wantErr  error
	}{
		{"single word", "hello", []byte("hello"), nil},
		{"word with space", "hello world", []byte("hello"), nil},
		{"word with tab", "hello\tworld", []byte("hello"), nil},
		{"word with newline", "hello\nworld", []byte("hello"), nil},
		{"word with multiple spaces", "hello   world", []byte("hello"), nil},
		{"empty input", "", nil, io.EOF},
		{"only spaces", "   ", nil, io.EOF},
		{"unicode letters", "café", []byte("café"), nil},
		{"mixed case", "HeLLo", []byte("HeLLo"), nil},
		{"cyrillic characters", "привет", []byte("привет"), nil},
		{"chinese characters", "你好", []byte("你好"), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := readWord(reader)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, got, tt.wantWord)
		})
	}
}

func TestReadWord_WrongPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantWord []byte
		wantErr  error
	}{
		{"word with digit", "hello123", nil, ErrInvalidCharacter},
		{"word with special character", "hello!world", nil, ErrInvalidCharacter},
		{"digit at start", "123hello", nil, ErrInvalidCharacter},
		{"word with underscore", "hello_world", nil, ErrInvalidCharacter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := readWord(reader)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, got)
		})
	}
}
