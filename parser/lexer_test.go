package parser

import (
	"bytes"
	"context"
	"fmt"
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
		buf := bytes.NewBuffer(make([]byte, 0, 1024))

		r, _, readErr := l.ReadRune()
		if readErr != nil {
			return nil, readErr
		}

		_, writeErr := buf.WriteRune(r)
		if writeErr != nil {
			return nil, writeErr
		}

		if bytes.Contains(buf.Bytes(), []byte("error")) {
			return nil, fmt.Errorf("contain the word 'error' which is not accepted into my imaginary test grammar")
		}

		return nil, nil
	}

	for _, tt := range tests {
		lexer := NewLexer(strings.NewReader(tt.input))
		err = Lex(t.Context(), lexer, initStateFunc)

		assert.ErrorIs(t, err, tt.err)
	}
}
