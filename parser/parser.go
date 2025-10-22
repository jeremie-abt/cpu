package parser

import (
	"cpu/lexer"
	"errors"
	"fmt"
)

type NodeType int

const (
	// NodeProgram is the ast root.
	NodeProgram NodeType = iota
	NodeInstruction
	NodeMemoryOperand
	NodeLabel
	NodeDirective
)

type Node interface {
	Type() NodeType
	Position() (line, column int)
	Children() []Node
}

// ProgramAST represents the AST of a program, each Node representing a sequential branch of the
// program with all its children.
type ProgramAST struct {
	statements []Node
	lexer      lexer.Lexer
}

func (p *ProgramAST) Type() NodeType       { return NodeProgram }
func (p *ProgramAST) Position() (int, int) { return 0, 0 }
func (p *ProgramAST) Children() []Node     { return p.statements }

// Instruction represents ASM instruction.
type Instruction struct {
	Line     int
	Column   int
	Label    *string // Label optional
	OpCode   OpCode  // MOV, ADD, JMP, etc.
	Operands []Node
}

//go:generate stringer -type=OpCode,RegisterName -output=stringer.go
type OpCode int

const (
	OpCodeMOV OpCode = iota
	OpCodeADD
	OpCodeSUB
	OpCodeMUL
	OpCodeDIV
	OpCodePUSH
	OpCodePOP
	OpCodeJMP
	OpCodeJE
	OpCodeJNE
	OpCodeCALL
	OpCodeRET
	OpCodeCMP
	OpCodeHALT

	// OpCodeUnknown must stay the last value of the enum as some test rely on
	OpCodeUnknown
)

type RegisterName int

const (
	RegisterEAX RegisterName = iota
	RegisterECX
	RegisterEDX
	RegisterEBX
	RegisterESI
	RegisterEDI
	RegisterESP
	RegisterEBP
	RegisterR8D
	RegisterR9D
	RegisterR10D
	RegisterR11D
	RegisterR12D
	RegisterR13D
	RegisterR14D
	RegisterR15D

	RegisterUnknown
)

func (i *Instruction) Type() NodeType       { return NodeInstruction }
func (i *Instruction) Position() (int, int) { return i.Line, i.Column }

func (i *Instruction) Children() []Node {
	var children []Node

	if i == nil {
		return children
	}

	for _, operand := range i.Operands {
		children = append(children, operand)
	}
	return children
}

type OperandType int

const (
	OperandRegister OperandType = iota
	OperandImmediate
	OperandMemory
	OperandLabel
)

type BaseOperand struct {
	Line   int
	Column int
}

func (b *BaseOperand) Type() NodeType       { return NodeInstruction }
func (b *BaseOperand) Position() (int, int) { return b.Line, b.Column }
func (b *BaseOperand) Children() []Node     { return nil }

type RegisterOperand struct {
	BaseOperand
	Name RegisterName
}

// TODO: quelle terminologie / dois-je vraiment renvoyer un NodeInstruction ou faire un type spécifique ?
func (b *RegisterOperand) Type() NodeType       { return NodeInstruction }
func (b *RegisterOperand) Position() (int, int) { return b.Line, b.Column }
func (b *RegisterOperand) Children() []Node     { return nil }

type ImmediateOperand struct {
	BaseOperand
	Value int
}

func (b *ImmediateOperand) Type() NodeType       { return NodeInstruction }
func (b *ImmediateOperand) Position() (int, int) { return b.Line, b.Column }
func (b *ImmediateOperand) Children() []Node     { return nil }

// MemoryOperand is based on this calcul Base + (Index × Scale) + Displacement
// where pretty every field are optional.
// https://www.intel.com/content/www/us/en/developer/articles/technical/intel-sdm.html
type MemoryOperand struct {
	BaseOperand
	Base         *RegisterName
	Index        *RegisterName
	Displacement int32
	ScaleFactor  uint8
}

func (b *MemoryOperand) Type() NodeType       { return NodeInstruction }
func (b *MemoryOperand) Position() (int, int) { return b.Line, b.Column }
func (b *MemoryOperand) Children() []Node     { return nil }

type LabelOperand struct {
	BaseOperand
	Name string
}

type ParseError struct {
	Line    int
	Column  int
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("ligne %d, colonne %d: %s", e.Line, e.Column, e.Message)
}

func NewProgramAST(lexer lexer.Lexer) *ProgramAST {
	return &ProgramAST{
		lexer:      lexer,
		statements: make([]Node, 0),
	}
}

// GenerateAST the ast generation.
// INPUT: tokens du lexer
// OUTPUT: *ProgramAST (AST) et slice d'erreurs
func GenerateAST(lexerObj lexer.Lexer) (*ProgramAST, []ParseError) {
	astProgram := NewProgramAST(lexerObj)

	program := &ProgramAST{
		statements: make([]Node, 0),
	}

	for token := astProgram.lexer.NextToken(); token.Type != lexer.TokenEOF; token = lexerObj.NextToken() {
		stmt, err := astProgram.parseStatement()
		fmt.Printf("could not parse statement: %v\n", err)

		if stmt != nil {
			program.statements = append(program.statements, stmt)
		}
	}

	return program, nil
}

var errEofLexer = errors.New("EOF lexer")

// parseStatement parse une ligne complète (label, instruction, ou directive)
func (p *ProgramAST) parseStatement() (Node, error) {
	// Parse avec le nextLine until on a un Node complet

	token := p.lexer.NextToken()
	if token.Type == lexer.TokenEOF {
		return nil, errEofLexer
	}

	if token.Type != lexer.TokenInstruction {
		return nil, errors.New("only token instruction are implemented from now")
	}

	p.parseInstruction(token.Literal)
	return nil, nil
}

// parseInstruction parse une instruction (MOV eax, 5)
func (p *ProgramAST) parseInstruction(label []byte) *Instruction {
	return nil
}

func parseMemoryOperand(l lexer.Lexer) (*MemoryOperand, error) {

	baseOperand := BaseOperand{
		Line:   l.Line(),
		Column: l.Column(),
	}

	tok := l.NextToken()
	if tok.Type != lexer.TokenOpenBracket {
		return nil, &parsingError{
			message: "could not parse memory operand, missing open bracket",
		}
	}

	var base *RegisterName
	var index *RegisterName
	var displacement int32
	var scaleFactor uint8

	tok = l.NextToken()
	if tok.Type == lexer.TokenRegister {
		register := getRegisterCodeFromSlice(tok.Literal)
		if register == RegisterUnknown {
			return nil, &parsingError{
				message: "invalid register",
			}
		}
		base = &register

		tok = l.NextToken()
		if tok.Type != lexer.TokenClosedBracket {
			if tok.Type == lexer.TokenPlus {
				tok = l.NextToken()

				if tok.Type == lexer.TokenNumber {
					displacement = int32(tok.Value)
					tok = l.NextToken()
				} else if tok.Type == lexer.TokenRegister {
					idxRegister := getRegisterCodeFromSlice(tok.Literal)
					index = &idxRegister

					if register == RegisterUnknown {
						return nil, &parsingError{
							message: "invalid register",
						}
					}

					tok = l.NextToken()
					if tok.Type == lexer.TokenPlus {
						tok = l.NextToken()

						if tok.Type == lexer.TokenNumber {
							displacement = int32(tok.Value)
							tok = l.NextToken()
						}
					}
				}
			}
		}
	} else if tok.Type == lexer.TokenNumber {
		displacement = int32(tok.Value)
		tok = l.NextToken()
	}

	if tok.Type != lexer.TokenClosedBracket {
		return nil, &parsingError{
			message: "invalid memory operand, missing closed bracket",
		}
	}

	return &MemoryOperand{
		BaseOperand:  baseOperand,
		Base:         base,
		Index:        index,
		Displacement: displacement,
		ScaleFactor:  scaleFactor,
	}, nil
}

// getRegisterCodeFromSlice is a heuristic returning which register is represented by the byte slice, if you
// pass a non-valid register, then the behavior is not defined.
func getRegisterCodeFromSlice(register []byte) RegisterName {
	switch register[0] {
	case 'E':
		if register[1] == 'A' {
			return RegisterEAX
		}
		if register[1] == 'C' {
			return RegisterECX
		}
		if register[1] == 'D' {
			if register[2] == 'X' {
				return RegisterEDX
			}
			return RegisterEDI
		}
		if register[1] == 'S' {
			if register[2] == 'I' {
				return RegisterESI
			}
			return RegisterESP
		}
		if register[1] == 'B' {
			if register[2] == 'P' {
				return RegisterEBP
			}
			return RegisterEBX
		}
		return RegisterUnknown
	case 'R':
		switch register[1] {
		case '8':
			return RegisterR8D
		case '9':
			return RegisterR9D
		default:
			switch register[2] {
			case '0':
				return RegisterR10D
			case '1':
				return RegisterR11D
			case '2':
				return RegisterR12D
			case '3':
				return RegisterR13D
			case '4':
				return RegisterR14D
			case '5':
				return RegisterR15D
			default:
				return RegisterUnknown
			}
		}
	default:
		return RegisterUnknown
	}
}

// getOpCodeFromSlice is an heuristic returning which OpCode is represented by the opcode slice, if you
// pass a not valid opcode then the behavior is not defined.
func getOpCodeFromSlice(opCode []byte) OpCode {
	opCodeToReturn, _ := getOpCodeFunctionFromSlice(opCode)
	return opCodeToReturn
}

func getOpCodeFunctionFromSlice(opCode []byte) (OpCode, func()) {
	switch opCode[0] {
	case 'M':
		if opCode[1] == 'O' {
			return OpCodeMOV, nil
		}
		return OpCodeMUL, nil
	case 'A':
		return OpCodeADD, nil
	case 'S':
		return OpCodeSUB, nil
	case 'D':
		return OpCodeDIV, nil
	case 'P':
		if opCode[1] == 'U' {
			return OpCodePUSH, nil
		}
		return OpCodePOP, nil
	case 'J':
		if opCode[1] == 'M' {
			return OpCodeJMP, nil
		} else if opCode[1] == 'E' {
			return OpCodeJE, nil
		}
		return OpCodeJNE, nil
	case 'C':
		if opCode[1] == 'M' {
			return OpCodeCMP, nil
		}
		return OpCodeCALL, nil
	case 'R':
		return OpCodeRET, nil
	case 'H':
		return OpCodeHALT, nil
	default:
		return OpCodeUnknown, nil
	}
}
