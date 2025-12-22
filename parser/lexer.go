// Lexer code will be responsible of tokenizing the input, could return errors based on the input grammar
// but not on the sense of it, all this file does is transforming suits of characters into suits of token,
// if a character does not exist in the grammar, then this part of the code will return an error etc ...
//
// But it does not interpret the meaning of the tokens together, it just returns them.
//
// This file corresponds to the Lexical analysis.
package parser

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

const MinLexingBuffSize = 1024

// LexTimeout is the timeout for the lexing process.
const LexTimeout = time.Second * 10

var LexerError = errors.New("lexer error")

type tokenType string

const (
	NUMBER tokenType = "STATEMENT"
	// ...
)

type token struct {
	value []rune
	tokenType
}

type Lexer interface {
	io.RuneScanner
	Emitter
}

// Emitter is the lexer emitter, basically it just sends a token within an internal chan so that the token
// can then be parsed.
type Emitter interface {
	Emit(token token)
}

type lexerFunc func(ctx context.Context, l Lexer) (lexerFunc, error)

type LexerImpl struct {
	reader bufio.Reader

	tokens chan<- *token
}

var _ Lexer = (*LexerImpl)(nil)

func (l *LexerImpl) ReadRune() (r rune, size int, err error) {
	r, s, err := l.reader.ReadRune()
	if err != nil {
		return 0, 0, fmt.Errorf("could not read rune : %w\n", err)
	}

	return r, s, err
}

func (l *LexerImpl) UnreadRune() error {
	if err := l.reader.UnreadRune(); err != nil {
		return fmt.Errorf("failed to unread rune: %w", err)
	}
	return nil
}

func (l *LexerImpl) Emit(token token) {
	l.tokens <- &token
	return
}

func NewLexer(input io.Reader) *LexerImpl {
	r := bufio.NewReaderSize(input, MinLexingBuffSize)

	return &LexerImpl{
		reader: *r,
	}
}

func Lex(ctx context.Context, l Lexer, initState lexerFunc) error {
	var err error
	state := initState
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(LexTimeout))

	defer cancel()

	for state != nil {
		state, err = state(ctx, l)

		if err != nil {
			return fmt.Errorf("failed to lex input: %w", err)
		}
	}

	return nil
}
