package parser

import (
	"cpu/lexer"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_handleParseMov(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		lexerObj := lexer.NewLexer()

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lexerObj.ResetWithInput(append([]byte("MOV "), tt.input...))
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
		// TODO: Implement memory handling
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
						Base: RegisterEBX.String(),
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
						Base:   RegisterEBP.String(),
						Offset: 8,
					},
				},
			},
		},
	},
	{
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
						Base: RegisterESI.String(),
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
	},
	{
		"negative offset",
		[]byte("EAX, [EBP-4]"),
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
						Base:   RegisterEBP.String(),
						Offset: -4,
					},
				},
			},
		},
	},
	{
		"SIB addressing",
		[]byte("EAX, [EBX+ECX*4]"),
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
						Base: RegisterEBX.String(),
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
