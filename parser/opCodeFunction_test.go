package parser

import (
	"cpu/lexer"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_handleMov(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []byte
		want  []Node
	}{
		{
			"simple move",
			[]byte("MOV EBX, EBP"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
		},
		{
			"move with immediate value",
			[]byte("MOV EAX, 42"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
		},
		{
			"move with hex immediate",
			[]byte("MOV ECX, 0xFF"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
		},
		{
			"move with memory operand",
			[]byte("MOV EAX, [EBX]"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
							Base: RegisterEBX,
						},
					},
				},
			},
		},
		{
			"move with memory offset",
			[]byte("MOV EDX, [EBP+8]"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
							Base:   RegisterEBP,
							Offset: 8,
						},
					},
				},
			},
		},
		{
			"move to memory from register",
			[]byte("MOV [ESI], EAX"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
					Operands: []Node{
						&MemoryOperand{
							BaseOperand: BaseOperand{
								Line:   0,
								Column: 4,
							},
							Base: RegisterESI,
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
			"move with negative offset",
			[]byte("MOV EAX, [EBP-4]"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
							Base:   RegisterEBP,
							Offset: -4,
						},
					},
				},
			},
		},
		{
			"move with SIB addressing",
			[]byte("MOV EAX, [EBX+ECX*4]"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
							Base:  RegisterEBX,
							Index: RegisterECX,
							Scale: 4,
						},
					},
				},
			},
		},
		{
			"move with extra whitespace",
			[]byte("MOV    EAX,    EBX"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
			[]byte("mov eax, ebx"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
			[]byte("MOV ESI, 0"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
			[]byte("MOV EDI, 0xFFFFFFFF"),
			[]Node{
				&Instruction{
					Column: 0,
					Line:   0,
					OpCode: OpCodeMOV,
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
	for _, tt := range tests {
		lexerObj := lexer.NewLexer()

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lexerObj.ResetWithInput(tt.input)
			got, err := handleMov(lexerObj)

			if err != nil {
				t.Fatalf("handleMov() error = %v", err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("handleMov() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
