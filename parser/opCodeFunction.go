package parser

import (
	"bytes"
	"cpu/lexer"
	"fmt"
)

type parsingError struct {
	message     string
	line        int
	column      int
	instruction OpCode
}

func (e parsingError) Error() string {
	return fmt.Sprintf("could not parse %s instruction line:%d, column: %d %s\n",
		e.instruction, e.line, e.column, e.message)
}

// opCodeParserHandler are parsing OpCode functions that handle the parsing of all OpCode.
type opCodeParserHandler func(lexer.Lexer) ([]Node, error)

// handleParseMov is the MOV parser handler, it supposes that the lexer is in the state just after
// having read the MOV instruction, so we are parsing mov parameters now.
func handleParseMov(l lexer.Lexer) ([]Node, error) {
	operands := make([]Node, 0)

	destination := l.NextToken()
	if destination.Type == lexer.TokenInstruction && bytes.Equal(destination.Literal, []byte("MOV")) {
		destination = l.NextToken()
	}
	comma := l.NextToken()
	source := l.NextToken()

	if comma == nil {
		return nil, &parsingError{
			message:     "comma is nil, syntax not understandable",
			instruction: OpCodeSUB,
		}
	}
	if comma.Type != lexer.TokenComma {
		return nil, &parsingError{
			message:     "missing comma",
			instruction: OpCodeMOV,
		}
	}
	if source == nil || destination == nil {
		return nil, &parsingError{
			message:     "missing operand",
			instruction: OpCodeMOV,
		}
	} else if source.Type == lexer.TokenEOF || destination.Type == lexer.TokenEOF {
		return nil, &parsingError{
			column:      source.Column,
			line:        source.Line,
			message:     "missing operand",
			instruction: OpCodeMOV,
		}
	}

	if source.Type != lexer.TokenNumber && source.Type != lexer.TokenRegister {
		return nil, &parsingError{
			column:  source.Column,
			line:    source.Line,
			message: "invalid source operand, must be a number or a register",
		}
	}

	if destination.Type != lexer.TokenRegister {
		return nil, &parsingError{
			column:  source.Column,
			line:    source.Line,
			message: "invalid source operand, must be a number or a register",
		}
	}

	if destination.Type == lexer.TokenRegister {
		operands = append(operands, &RegisterOperand{
			Name: getRegisterCodeFromSlice(destination.Literal),
			BaseOperand: BaseOperand{
				Line:   destination.Line,
				Column: destination.Column,
			},
		})
	}

	if source.Type == lexer.TokenRegister {
		operands = append(operands, &RegisterOperand{
			Name: getRegisterCodeFromSlice(source.Literal),
			BaseOperand: BaseOperand{
				Line:   source.Line,
				Column: source.Column,
			},
		})
	} else if source.Type == lexer.TokenNumber {
		operands = append(operands, &ImmediateOperand{
			Value: source.Value,
			BaseOperand: BaseOperand{
				Line:   source.Line,
				Column: source.Column,
			},
		})
	}

	return []Node{
		&Instruction{
			Column:   0,
			Line:     destination.Line,
			OpCode:   OpCodeMOV,
			Operands: operands,
		},
	}, nil
}

// handleParseSub is the SUB parser handler, it supposes that the lexer is in the state just after
// having read the SUB instruction, so we are parsing sub parameters now.
func handleParseSub(l lexer.Lexer) ([]Node, error) {
	operands := make([]Node, 0)

	destination := l.NextToken()
	if destination.Type == lexer.TokenInstruction && bytes.Equal(destination.Literal, []byte("SUB")) {
		destination = l.NextToken()
	}
	comma := l.NextToken()
	source := l.NextToken()

	if comma == nil {
		return nil, &parsingError{
			message:     "comma is nil, syntax not understandable",
			instruction: OpCodeSUB,
		}
	}
	if comma.Type != lexer.TokenComma {
		return nil, &parsingError{
			message:     "missing comma",
			instruction: OpCodeSUB,
		}
	}
	if source == nil || destination == nil {
		return nil, &parsingError{
			message:     "missing operand",
			instruction: OpCodeSUB,
		}
	} else if source.Type == lexer.TokenEOF || destination.Type == lexer.TokenEOF {
		return nil, &parsingError{
			column:      source.Column,
			line:        source.Line,
			message:     "missing operand",
			instruction: OpCodeSUB,
		}
	}

	if source.Type != lexer.TokenNumber && source.Type != lexer.TokenRegister {
		return nil, &parsingError{
			column:  source.Column,
			line:    source.Line,
			message: "invalid source operand, must be a number or a register",
		}
	}

	if destination.Type != lexer.TokenRegister {
		return nil, &parsingError{
			column:  source.Column,
			line:    source.Line,
			message: "invalid source operand, must be a number or a register",
		}
	}

	if destination.Type == lexer.TokenRegister {
		operands = append(operands, &RegisterOperand{
			Name: getRegisterCodeFromSlice(destination.Literal),
			BaseOperand: BaseOperand{
				Line:   destination.Line,
				Column: destination.Column,
			},
		})
	}

	if source.Type == lexer.TokenRegister {
		operands = append(operands, &RegisterOperand{
			Name: getRegisterCodeFromSlice(source.Literal),
			BaseOperand: BaseOperand{
				Line:   source.Line,
				Column: source.Column,
			},
		})
	} else if source.Type == lexer.TokenNumber {
		operands = append(operands, &ImmediateOperand{
			Value: source.Value,
			BaseOperand: BaseOperand{
				Line:   source.Line,
				Column: source.Column,
			},
		})
	}

	return []Node{
		&Instruction{
			Column:   0,
			Line:     destination.Line,
			OpCode:   OpCodeSUB,
			Operands: operands,
		},
	}, nil
}
