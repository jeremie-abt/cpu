package parser

import (
	"context"
	"strings"
)

// orchestrator will keep all channels and make the communication and the all process works.
type orchestrator struct {
	tokens chan []rune
}

func (o *orchestrator) Run(ctx context.Context, input []byte) (err error) {
	return Lex(ctx, NewLexer(strings.NewReader(string(input))), defaultLexState)
}
