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
	"log"
	"time"
)

const MinLexingBuffSize = 1024

// LexTimeout is the timeout for the lexing process.
const LexTimeout = time.Second * 10
const maxBackPressureTolerated = time.Millisecond * 300
const lexerChanSize = 32

var LexerError = errors.New("lexer error")

type tokenType string

const (
	TokenTypeStatement tokenType = "STATEMENT"
	TokenTypeWord      tokenType = "WORD"
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
	Emit(ctx context.Context, token token) error
}

type lexerFunc func(ctx context.Context, l Lexer) (lexerFunc, error)

type LexerImpl struct {
	reader bufio.Reader

	tokens chan *token
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

// TODO: Faire un logger propre
func (l *LexerImpl) Emit(ctx context.Context, token token) error {
	var cancel context.CancelFunc
	timeoutCh := time.After(maxBackPressureTolerated)

	for {
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(time.Millisecond*100))

		select {
		case <-ctx.Done():
			log.Println("[WARN] context deadline exceeded while emitting token, " +
				"be careful there might be some back pressure to investigate")
		case <-timeoutCh:
			cancel()
			return fmt.Errorf("exceeded max back pressure tolerance of %x", maxBackPressureTolerated)
		case l.tokens <- &token:
			cancel()
			return nil
		}

		cancel()
	}
}

func NewLexer(input io.Reader) *LexerImpl {
	r := bufio.NewReaderSize(input, MinLexingBuffSize)

	return &LexerImpl{
		reader: *r,
		tokens: make(chan *token, lexerChanSize),
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
			return fmt.Errorf("%w failed to lex input: %w", LexerError, err)
		}
	}

	return nil
}
