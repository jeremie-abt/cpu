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

var emitToEachWord = func(ctx context.Context, l Lexer) (lexerFunc, error) {
	var w []byte
	var err error
	buf := bytes.NewBufferString("")

	for {
		w, err = readWord(l)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		buf.Reset()
		_, err = buf.Write(w)
		if err != nil {
			return nil, err
		}

		if err := l.Emit(ctx, token{
			value: bytes.Runes(buf.Bytes()),
		}); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// initStateFunc parse just every word composed of letter and return an error if it encounters the word 'error',
// this function does not emit.
var initStateFunc = func(ctx context.Context, l Lexer) (lexerFunc, error) {
	var w []byte
	var err error
	buf := bytes.NewBufferString("")

	for w, err = readWord(l); !errors.Is(err, io.EOF); w, err = readWord(l) {
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(w)
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

	for _, tt := range tests {
		lexer := NewLexer(strings.NewReader(tt.input))
		err = Lex(t.Context(), lexer, initStateFunc)

		assert.ErrorIs(t, err, tt.err)
	}
}

func TestLexerGreenPath(t *testing.T) {
	lexer := NewLexer(strings.NewReader("hello world"))
	err := Lex(t.Context(), lexer, emitToEachWord)

	assert.NoError(t, err)
	fmt.Println(len(lexer.tokens))

	select {
	case tok := <-lexer.tokens:
		assert.Equal(t, []rune("hello"), tok)
	case <-t.Context().Done():
		t.Fatalf("context deadline exceeded, could not get the hello word")
	}

	select {
	case tok := <-lexer.tokens:
		assert.Equal(t, []rune("world"), tok)
	case <-t.Context().Done():
		t.Fatalf("context deadline exceeded, could not get the world word")
	}
}
