package parser

import (
	"bytes"
	"cpu/lexer"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func toPtr[T any](i T) *T {
	return &i
}

func Test_handleParseMov(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		lexerObj := lexer.NewLexer()

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lexerObj.ResetWithInput(append([]byte("MOV "), tt.input...))

			validateInstruction(t, lexerObj, []byte("MOV"))

			got, err := handleParseMov(lexerObj)

			if err != nil {
				t.Fatalf("handleParseMov() error = %v", err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("handleParseMov() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func validateInstruction(t *testing.T, lexerObj *lexer.ImplLexer, instructionLiteral []byte) {
	t.Helper()

	token := lexerObj.NextToken()
	if token.Type != lexer.TokenInstruction {
		t.Fatalf("got token that is not an instruction")
	}
	if !bytes.Equal(token.Literal, instructionLiteral) {
		t.Fatalf("got token that is not %s instruction", instructionLiteral)
	}
}

var tests = []struct {
	name  string
	input []byte
	want  []Node
}{
	{
		"simple instruction",
		[]byte("EBX, EBP"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterEBX,
					},
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Name: RegisterEBP,
					},
				},
			},
		},
	}, {
		"with immediate value",
		[]byte("EAX, 42"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterEAX,
					},
					&ImmediateOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Value: 42,
					},
				},
			},
		},
	}, {
		"with hex immediate",
		[]byte("ECX, 0xFF"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterECX,
					},
					&ImmediateOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Value: 0xFF,
					},
				},
			},
		},
	}, {
		"with memory operand",
		[]byte("EAX, [EBX]"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterEAX,
					},
					&MemoryOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Base: toPtr(RegisterEBX),
					},
				},
			},
		},
	},
	{
		"with memory offset",
		[]byte("EDX, [EBP+8]"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterEDX,
					},
					&MemoryOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Base:         toPtr(RegisterEBP),
						Displacement: 8,
					},
				},
			},
		},
	},
	/*{
		"memory from register",
		[]byte("[ESI], EAX"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&MemoryOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Base: toPtr(RegisterESI),
					},
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 11,
						},
						Name: RegisterEAX,
					},
				},
			},
		},
	},*/
	{
		"SIB addressing",
		[]byte("EAX, [EBX+(ECX*4)]"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterEAX,
					},
					&MemoryOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Base:        toPtr(RegisterEBX),
						Index:       toPtr(RegisterECX),
						ScaleFactor: 4,
					},
				},
			},
		},
	},
	{
		"extra whitespace",
		[]byte("   EAX,    EBX"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 7,
						},
						Name: RegisterEAX,
					},
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 15,
						},
						Name: RegisterEBX,
					},
				},
			},
		},
	},
	{
		"move lowercase",
		[]byte("eax, ebx"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterEAX,
					},
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Name: RegisterEBX,
					},
				},
			},
		},
	},
	{
		"move with zero immediate",
		[]byte("ESI, 0"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterESI,
					},
					&ImmediateOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Value: 0,
					},
				},
			},
		},
	},
	{
		"move with large immediate",
		[]byte("EDI, 0xFFFFFFFF"),
		[]Node{
			&Instruction{
				Column: 0,
				Line:   0,
				Operands: []Node{
					&RegisterOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 4,
						},
						Name: RegisterEDI,
					},
					&ImmediateOperand{
						BaseOperand: BaseOperand{
							Line:   0,
							Column: 9,
						},
						Value: 0xFFFFFFFF,
					},
				},
			},
		},
	},
}
