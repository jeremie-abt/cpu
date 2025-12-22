// This file code will receive tokens, verify that they are semantically correct and return a syntax tree.
package parser

type Parser interface {
	Interpret(t *token)
}

type Tree interface {
	// Node push the token into the current node as a child, moove the cursor on the child and make the link
	// between the parent and the child.
	Node(t *token)

	// Next returns the next token from the lexer.
	Next() *token
}

type ParserFunc func(t Tree) ParserFunc
type ParserFunc1 func(t *token) ParserFunc // ?
