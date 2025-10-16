// ImplLexer must

package lexer

import (
	"bytes"
	"maps"
	"strconv"
	"unique"
)

//go:generate stringer -type=TokenType
type TokenType int

var eof = byte(0)

const (
	TokenEOF TokenType = iota
	TokenNewline

	TokenRegister    // R0, R1, SP, PC
	TokenNumber      // 42, 0xFF, 0b1010
	TokenIdentifier  // labels, symboles
	TokenInstruction // MOV, ADD, PUSH etc ... Those are cpu instructions

	TokenOpenBracket   // [
	TokenClosedBracket // ]
	TokenComma         // ,
	TokenColon         // :

	TokenPlus   // +
	TOKEN_MINUS // -
	TOKEN_STAR  // *
	TOKEN_SLASH // /

	TokenIllegal // Error token
)

// TODO: Potentiel Hotpath sur le value
type Token struct {
	Type    TokenType
	Literal []byte
	Column  int
	Value   int32
	Line    int
}

func (t TokenType) GoString() string {
	return t.String()
}

type Lexer interface {
	NextToken() *Token
}

func NewLexerWithInput(input []byte) *ImplLexer {
	lexer := &ImplLexer{}
	lexer.setInput(input)
	return lexer
}

func NewLexer() *ImplLexer {
	return &ImplLexer{}
}

// ResetWithInput stop the current reads and restart it with the given input.
func (i *ImplLexer) ResetWithInput(input []byte) {
	i.setInput(input)
}

func (i *ImplLexer) setInput(input []byte) {
	i.rp = 0
	i.column = 0
	i.line = 0
	i.input = input

	if len(input) == 0 {
		i.value = eof
		return
	}
	i.value = input[0]
}

type ImplLexer struct {
	input  []byte
	rp     int // rp: current position we are reading
	value  byte
	line   int
	column int
}

var _ Lexer = (*ImplLexer)(nil)

func (i *ImplLexer) NextToken() *Token {
	var tok Token

	i.skipWhitespace()

	tok.Line = i.line
	tok.Column = i.column

	// TODO: étude de Hot path sur la conversion []byte{i.value}
	literal := []byte{i.value}
	switch i.value {
	case '\n':
		tok = Token{Type: TokenNewline, Literal: literal}
		i.line++
		i.column = 0

	case ',':
		tok = Token{Type: TokenComma, Literal: literal}

	case ':':
		tok = Token{Type: TokenColon, Literal: literal}

	case '[':
		tok = Token{Type: TokenOpenBracket, Literal: literal}

	case ']':
		tok = Token{Type: TokenClosedBracket, Literal: literal}

	case '+':
		tok = Token{Type: TokenPlus, Literal: literal}
	case ';':
		// TODO: Handle all type of comments
		i.skipComment()
		return i.NextToken()
	case 0:
		tok = Token{Type: TokenEOF, Literal: []byte{eof}}

	default:
		if isLetter(i.value) {
			tok.Literal = i.readTextBlock()
			tok.Type = lookupIdentType(tok.Literal)
			return &tok
		} else if isDigit(i.value) {
			tok.Literal = i.readNumber()
			tok.Value = i.parseNumber(tok.Literal)
			tok.Type = TokenNumber
			return &tok
		} else {
			tok = Token{Type: TokenIllegal, Literal: literal}
		}
	}

	i.nextChar()
	return &tok
}

// lEXER utils

// nextChar simple iteration for the current character.
func (i *ImplLexer) nextChar() {
	if i.rp >= len(i.input) {
		i.value = eof
	} else {
		i.value = accessInputSafe(i.input, i.rp+1)
	}
	i.rp++
	i.column++
}

func (i *ImplLexer) readChar() {
	if i.rp >= len(i.input) {
		i.value = eof
	} else {
		i.value = i.input[i.rp]
	}
	i.rp++
	i.column++
}

// peekChar returns the current character with overflow protection.
func (i *ImplLexer) peekChar() byte {
	if i.rp >= len(i.input) {
		return 0
	}
	return i.input[i.rp]
}

// peekChar returns the current character with overflow protection.
func (i *ImplLexer) peekCharPosition(idx int) byte {
	if i.rp+idx >= len(i.input) || i.rp+idx < 0 {
		return 0
	}
	return i.input[i.rp+idx]
}

// readTextBlock read a suite of text composed of letters, digits or underscore, it does not flag it as
// an identifier / instruction or other, the caller must do that.
func (i *ImplLexer) readTextBlock() []byte {
	position := i.rp

	for isLetter(i.value) || isDigit(i.value) || i.value == '_' {
		i.nextChar()
	}

	return i.input[position:i.rp]
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func (i *ImplLexer) readNumber() []byte {
	position := i.rp

	val := i.peekChar()
	if i.value == '0' && (val == 'x' || i.peekChar() == 'X') {
		i.nextChar() // '0'
		i.nextChar() // 'x'
		for isHexDigit(i.value) {
			i.nextChar()
		}
		return i.input[position:i.rp]
	}

	if i.value == '0' && (i.peekChar() == 'b' || i.peekChar() == 'B') {
		i.nextChar() // '0'
		i.nextChar() // 'b'
		for i.value == '0' || i.value == '1' {
			i.nextChar()
		}
		return i.input[position:i.rp]
	}

	for isDigit(i.value) {
		i.nextChar()
	}

	return i.input[position:i.rp]
}

func (i *ImplLexer) parseNumber(literal []byte) int32 {
	if bytes.HasPrefix(literal, []byte("0x")) || bytes.HasPrefix(literal, []byte("0X")) {
		// TODO: Je pense que si mon fichier fini par 0X ca crash, détecter ca avec le fuzz.
		val, _ := strconv.ParseInt(string(literal[2:]), 16, 32)
		return int32(val)
	}

	if bytes.HasPrefix(literal, []byte("0b")) || bytes.HasPrefix(literal, []byte("0B")) {
		val, _ := strconv.ParseInt(string(literal[2:]), 2, 32)
		return int32(val)
	}

	val, _ := strconv.ParseInt(string(literal), 10, 32)
	return int32(val)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func (i *ImplLexer) skipWhitespace() {
	for i.value == ' ' || i.value == '\t' || i.value == '\r' {
		i.nextChar()
	}
}

func (i *ImplLexer) skipComment() {
	for i.value != '\n' && i.value != eof {
		i.readChar()
	}
}

// TODO(hotpath): Essayer de faire de l'internalisation
var keywords = map[string]TokenType{
	"MOV":  TokenInstruction,
	"ADD":  TokenInstruction,
	"SUB":  TokenInstruction,
	"MUL":  TokenInstruction,
	"DIV":  TokenInstruction,
	"PUSH": TokenInstruction,
	"POP":  TokenInstruction,
	"JMP":  TokenInstruction,
	"JE":   TokenInstruction,
	"JNE":  TokenInstruction,
	"CALL": TokenInstruction,
	"RET":  TokenInstruction,
	"CMP":  TokenInstruction,
	"HALT": TokenInstruction,

	"R0": TokenRegister,
	"R1": TokenRegister,
	"R2": TokenRegister,
	"R3": TokenRegister,
	"R4": TokenRegister,
	"R5": TokenRegister,
	"R6": TokenRegister,
	"R7": TokenRegister,
	"SP": TokenRegister,
	"PC": TokenRegister,
}

func lookupIdentType(ident []byte) TokenType {
	identUpper := bytes.ToUpper(ident)

	if tokType, ok := keywords[string(identUpper)]; ok {
		return tokType
	}

	return TokenIdentifier
}

func accessInputSafe(input []byte, idx int) byte {
	if idx >= len(input) {
		return eof
	}

	return input[idx]
}

var keywordsOptimized = map[string]TokenType{}

func initLookupIdentTypeOptimized() {
	for key := range maps.Keys(keywords) {
		my := unique.Make(key)
		keywordsOptimized[my.Value()] = keywords[key]
	}
}

func lookupIdentTypeOptimized(ident []byte) TokenType {
	val := unique.Make(string(bytes.ToUpper(ident)))
	if tokType, ok := keywords[val.Value()]; ok {
		return tokType
	}

	return TokenIdentifier
}
