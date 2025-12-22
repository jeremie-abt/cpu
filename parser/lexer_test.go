package parser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLexerEdgesCases must only test the lexer codes, not the reel state function.
func TestLexerEdgeCases(t *testing.T) {
	var err error

	var tests = []struct {
		input string
		err   error
	}{
		{
			"this must return an error",
			LexerError,
		},
	}

	initStateFunc := func(ctx context.Context, l Lexer) (lexerFunc, error) {
		var r rune
		buf := bytes.NewBuffer(make([]byte, 0, 1024))

		for r, _, err = l.ReadRune(); !errors.Is(err, io.EOF); r, _, err = l.ReadRune() {
			if err != nil {
				return nil, err
			}

			_, err = buf.WriteRune(r)
			if err != nil {
				return nil, err
			}
		}

		if bytes.Contains(buf.Bytes(), []byte("error")) {
			return nil, fmt.Errorf("contain the word 'error' which is not accepted " +
				"into my imaginary test grammar")
		}

		return nil, nil
	}

	for _, tt := range tests {
		lexer := NewLexer(strings.NewReader(tt.input))
		err = Lex(t.Context(), lexer, initStateFunc)

		assert.ErrorIs(t, err, tt.err)
	}
}
