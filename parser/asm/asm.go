// Package asm is responsible to parse assembly code using the parser package.
//
// Note: I'm note sure about the organisation of the packages, will probably moove.
package asm

import "cpu/parser"

// Run will orchestrate the whole process of parsing which consists of:
// - Lexical analysis
// - Syntax analysis
// - Semantic analysis.
//
// Probably will orchestrate the compilation later too, we'll see.
func Run(input []byte, l parser.Lexer) {

}
