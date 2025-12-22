package parser

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWord_GreenPath(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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

func TestReadWord_MultipleWords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantWords [][]byte
	}{
		{
			name:      "two words space",
			input:     "hello world",
			wantWords: [][]byte{[]byte("hello"), []byte("world")},
		},
		{
			name:      "three words space",
			input:     "hello beautiful world",
			wantWords: [][]byte{[]byte("hello"), []byte("beautiful"), []byte("world")},
		},
		{
			name:      "words with tabs",
			input:     "hello\tworld\ttest",
			wantWords: [][]byte{[]byte("hello"), []byte("world"), []byte("test")},
		},
		{
			name:      "words with newlines",
			input:     "hello\nworld\ntest",
			wantWords: [][]byte{[]byte("hello"), []byte("world"), []byte("test")},
		},
		{
			name:      "words with mixed whitespace",
			input:     "hello  \t\n  world",
			wantWords: [][]byte{[]byte("hello"), []byte("world")},
		},
		{
			name:      "unicode words",
			input:     "café привет 你好",
			wantWords: [][]byte{[]byte("café"), []byte("привет"), []byte("你好")},
		},
		{
			name:      "single word",
			input:     "hello",
			wantWords: [][]byte{[]byte("hello")},
		},
		{
			name:      "words with trailing spaces",
			input:     "hello world   ",
			wantWords: [][]byte{[]byte("hello"), []byte("world")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var gotWords = make([][]byte, 0, len(tt.wantWords))

			for i := 1; ; i++ {
				word, err := readWord(reader)

				if len(word) > 0 {
					gotWords = append(gotWords, word)
				}

				if errors.Is(err, io.EOF) {
					break
				}

				if err != nil {
					t.Fatalf("readWord() unexpected error at iteration %d: %v", i, err)
				}

				if i%20 == 0 {
					t.Fatalf("the test is probably in infinite loop, if you're testing long inputs," +
						" consider use smaller ones or upgrade the maximum iteration count")
				}
			}

			assert.Equal(t, tt.wantWords, gotWords)
		})
	}
}
