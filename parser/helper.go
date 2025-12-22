package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
)

var ErrInvalidCharacter = errors.New("invalid character")

// readWord reads a word from the reader
func readWord(input io.RuneReader) ([]byte, error) {
	var r rune
	var err error
	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	r, _, err = input.ReadRune()
	if errors.Is(err, io.EOF) {
		return nil, io.EOF
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read rune: %w", err)
	}

	// In case we have space between to word, we skip it
	for unicode.IsSpace(r) {
		r, _, err = input.ReadRune()

		if errors.Is(err, io.EOF) {
			return nil, io.EOF
		}
	}

	for {
		if unicode.IsSpace(r) {
			break
		}

		if !unicode.IsLetter(r) {
			return nil, ErrInvalidCharacter
		}

		_, err = buf.WriteRune(r)
		if err != nil {
			return nil, fmt.Errorf("failed to write rune: %w", err)
		}

		r, _, err = input.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read rune: %w", err)
		}
	}

	if buf.Len() == 0 {
		return nil, io.EOF
	}
	return buf.Bytes(), nil
}
