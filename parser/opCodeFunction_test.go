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
