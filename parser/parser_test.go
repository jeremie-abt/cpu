package parser

import (
	"cpu/lexer"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseMemoryOperand(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  Node
	}{
		{
			name:  "[base]",
			input: "[EAX]",
			want: &MemoryOperand{
				Base: toPtr(RegisterEAX),
			},
		},
		{
			name:  "[displacement hexa]",
			input: "[0XFF]",
			want: &MemoryOperand{
				Displacement: 255,
			},
		},
		{
			name:  "[displacement binary]",
			input: "[0b11111111]",
			want: &MemoryOperand{
				Displacement: 255,
			},
		},
		{
			name:  "[base + displacement]",
			input: "[EAX + 0XFF]",
			want: &MemoryOperand{
				Displacement: 255,
				Base:         toPtr(RegisterEAX),
			},
		},
		{
			name:  "[base + displacement]",
			input: "[EAX + 0XFF]",
			want: &MemoryOperand{
				Displacement: 255,
				Base:         toPtr(RegisterEAX),
			},
		},
		{
			name:  "[base + index + displacement]",
			input: "[EAX + EBP + 0XFF]",
			want: &MemoryOperand{
				Displacement: 255,
				Index:        toPtr(RegisterEBP),
				Base:         toPtr(RegisterEAX),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			lexerObj := lexer.NewLexer()
			lexerObj.ResetWithInput([]byte(tt.input))

			got, err := parseMemoryOperand(lexerObj)

			if err != nil {
				t.Fatalf("parseMemoryOperand() erreur inattendue = %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("parseMemoryOperand() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetOpCodeFromSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []byte
		expected OpCode
	}{
		{
			name:     "MOV opcode",
			input:    []byte("MOV"),
			expected: OpCodeMOV,
		},
		{
			name:     "ADD opcode",
			input:    []byte("ADD"),
			expected: OpCodeADD,
		},
		{
			name:     "SUB opcode",
			input:    []byte("SUB"),
			expected: OpCodeSUB,
		},
		{
			name:     "MUL opcode",
			input:    []byte("MUL"),
			expected: OpCodeMUL,
		},
		{
			name:     "DIV opcode",
			input:    []byte("DIV"),
			expected: OpCodeDIV,
		},
		{
			name:     "PUSH opcode",
			input:    []byte("PUSH"),
			expected: OpCodePUSH,
		},
		{
			name:     "POP opcode",
			input:    []byte("POP"),
			expected: OpCodePOP,
		},
		{
			name:     "JMP opcode",
			input:    []byte("JMP"),
			expected: OpCodeJMP,
		},
		{
			name:     "JE opcode",
			input:    []byte("JE"),
			expected: OpCodeJE,
		},
		{
			name:     "JNE opcode",
			input:    []byte("JNE"),
			expected: OpCodeJNE,
		},
		{
			name:     "CALL opcode",
			input:    []byte("CALL"),
			expected: OpCodeCALL,
		},
		{
			name:     "RET opcode",
			input:    []byte("RET"),
			expected: OpCodeRET,
		},
		{
			name:     "CMP opcode",
			input:    []byte("CMP"),
			expected: OpCodeCMP,
		},
		{
			name:     "HALT opcode",
			input:    []byte("HALT"),
			expected: OpCodeHALT,
		},
	}

	var opCodeLengthWanted = 14
	if int(OpCodeUnknown) != opCodeLengthWanted {
		t.Fatalf("length of the OpCode enum has changed, are you sure this test is testing all opcodes ?" +
			" if you've added a new opcode recently, you must update the test to be sure the function behave well" +
			" and then update the opCodeLengthWanted variable")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getOpCodeFromSlice(tt.input)
			if result != tt.expected {
				t.Errorf("getOpCodeFromSlice(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func Test_getRegisterCodeFromSlice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    []byte
		expected RegisterName
	}{
		{
			name:     "EAX",
			input:    []byte("EAX"),
			expected: RegisterEAX,
		},
		{
			name:     "ECX",
			input:    []byte("ECX"),
			expected: RegisterECX,
		},
		{
			name:     "EDX",
			input:    []byte("EDX"),
			expected: RegisterEDX,
		},
		{
			name:     "EBX",
			input:    []byte("EBX"),
			expected: RegisterEBX,
		},
		{
			name:     "ESI",
			input:    []byte("ESI"),
			expected: RegisterESI,
		},
		{
			name:     "EDI",
			input:    []byte("EDI"),
			expected: RegisterEDI,
		},
		{
			name:     "ESP",
			input:    []byte("ESP"),
			expected: RegisterESP,
		},
		{
			name:     "EBP",
			input:    []byte("EBP"),
			expected: RegisterEBP,
		},
		{
			name:     "R8D",
			input:    []byte("R8D"),
			expected: RegisterR8D,
		},
		{
			name:     "R9D",
			input:    []byte("R9D"),
			expected: RegisterR9D,
		},
		{
			name:     "R10D",
			input:    []byte("R10D"),
			expected: RegisterR10D,
		},
		{
			name:     "R11D",
			input:    []byte("R11D"),
			expected: RegisterR11D,
		},
		{
			name:     "R12D",
			input:    []byte("R12D"),
			expected: RegisterR12D,
		},
		{
			name:     "R13D",
			input:    []byte("R13D"),
			expected: RegisterR13D,
		},
		{
			name:     "R14D",
			input:    []byte("R14D"),
			expected: RegisterR14D,
		},
		{
			name:     "R15D",
			input:    []byte("R15D"),
			expected: RegisterR15D,
		},
	}

	var registerLengthWanted = 16
	if int(RegisterUnknown) != registerLengthWanted {
		t.Fatalf("length of the register enum has changed, are you sure this test is testing all registers ?" +
			" if you've added a new register recently, you must update the test to be sure the function behave well" +
			" and then update the registerLengthWanted variable")
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := getRegisterCodeFromSlice(tt.input)
			if result != tt.expected {
				t.Errorf("getRegisterCodeFromSlice(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}

}
