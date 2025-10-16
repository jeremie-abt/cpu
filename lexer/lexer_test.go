package lexer

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const maxTokenAuthorizedForTest = 1000000

func TestNewLexer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  []byte
		wanted []Token
	}{
		{
			name:  "simple one liner",
			input: []byte("mov r0, r1"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("r1")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "instruction with immediate value",
			input: []byte("mov r0, 42"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("42")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "hexadecimal number",
			input: []byte("mov r0, 0xFF"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("0xFF")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "binary number",
			input: []byte("mov r0, 0b1010"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("0b1010")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "label definition",
			input: []byte("loop:"),
			wanted: []Token{
				{Type: TokenIdentifier, Literal: []byte("loop")},
				{Type: TokenColon, Literal: []byte(":")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "label with numbers",
			input: []byte("loop1:"),
			wanted: []Token{
				{Type: TokenIdentifier, Literal: []byte("loop1")},
				{Type: TokenColon, Literal: []byte(":")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "label with underscore",
			input: []byte("my_loop_2:"),
			wanted: []Token{
				{Type: TokenIdentifier, Literal: []byte("my_loop_2")},
				{Type: TokenColon, Literal: []byte(":")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "jump to label",
			input: []byte("jmp loop"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("jmp")},
				{Type: TokenIdentifier, Literal: []byte("loop")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "arithmetic operations",
			input: []byte("add r0, r1"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("add")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("r1")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "stack operations",
			input: []byte("push r0"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("push")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "memory access with brackets",
			input: []byte("load r0, [r1]"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("load")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenOpenBracket, Literal: []byte("[")},
				{Type: TokenRegister, Literal: []byte("r1")},
				{Type: TokenClosedBracket, Literal: []byte("]")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "memory access with address",
			input: []byte("store [0x1000], r0"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("store")},
				{Type: TokenOpenBracket, Literal: []byte("[")},
				{Type: TokenNumber, Literal: []byte("0x1000")},
				{Type: TokenClosedBracket, Literal: []byte("]")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "multiple instructions with newlines",
			input: []byte("mov r0, 42\nadd r0, 1"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("42")},
				{Type: TokenNewline, Literal: []byte("\n")},
				{Type: TokenInstruction, Literal: []byte("add")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("1")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "comment single line",
			input: []byte("mov r0, 42 ; this is a comment"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("42")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "comment full line",
			input: []byte("; full line comment\nmov r0, 1"),
			wanted: []Token{
				{Type: TokenNewline, Literal: []byte("\n")},
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("1")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "special registers",
			input: []byte("mov sp, r0"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("sp")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "compare instruction",
			input: []byte("cmp r0, r1"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("cmp")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("r1")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "conditional jump",
			input: []byte("je label"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("je")},
				{Type: TokenIdentifier, Literal: []byte("label")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "function call",
			input: []byte("call function"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("call")},
				{Type: TokenIdentifier, Literal: []byte("function")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "return instruction",
			input: []byte("ret"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("ret")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "halt instruction",
			input: []byte("halt"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("halt")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "nop instruction",
			input: []byte("nop"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("nop")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "multiple whitespaces",
			input: []byte("mov    r0,    42"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("42")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "tabs and spaces mixed",
			input: []byte("mov\tr0,\t42"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("42")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "case insensitive instructions",
			input: []byte("MOV R0, R1"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("MOV")},
				{Type: TokenRegister, Literal: []byte("R0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("R1")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "mixed case",
			input: []byte("MoV r0, R1"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("MoV")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("R1")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "zero value",
			input: []byte("mov r0, 0"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("0")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "large number",
			input: []byte("mov r0, 999999"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("999999")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "all registers",
			input: []byte("mov r7, r6"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("r7")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("r6")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "program counter register",
			input: []byte("mov pc, r0"),
			wanted: []Token{
				{Type: TokenInstruction, Literal: []byte("mov")},
				{Type: TokenRegister, Literal: []byte("pc")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "empty input",
			input: []byte(""),
			wanted: []Token{
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "only whitespace",
			input: []byte("   \t\n  "),
			wanted: []Token{
				{Type: TokenNewline, Literal: []byte("\n")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "only comment",
			input: []byte("; just a comment"),
			wanted: []Token{
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "label and instruction on same line",
			input: []byte("loop: add r0, 1"),
			wanted: []Token{
				{Type: TokenIdentifier, Literal: []byte("loop")},
				{Type: TokenColon, Literal: []byte(":")},
				{Type: TokenInstruction, Literal: []byte("add")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("1")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
		{
			name:  "simple program",
			input: []byte("loop:\n  add r0, 1\n  jmp loop"),
			wanted: []Token{
				{Type: TokenIdentifier, Literal: []byte("loop")},
				{Type: TokenColon, Literal: []byte(":")},
				{Type: TokenNewline, Literal: []byte("\n")},
				{Type: TokenInstruction, Literal: []byte("add")},
				{Type: TokenRegister, Literal: []byte("r0")},
				{Type: TokenComma, Literal: []byte(",")},
				{Type: TokenNumber, Literal: []byte("1")},
				{Type: TokenNewline, Literal: []byte("\n")},
				{Type: TokenInstruction, Literal: []byte("jmp")},
				{Type: TokenIdentifier, Literal: []byte("loop")},
				{Type: TokenEOF, Literal: []byte("")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			lexer := NewLexerWithInput(tt.input)

			tokensGotFromLexer := make([]*Token, 0)
			for i := 0; i < maxTokenAuthorizedForTest; i++ {
				newToken := lexer.NextToken()
				if newToken.Type == TokenEOF {
					break
				}

				tokensGotFromLexer = append(tokensGotFromLexer, newToken)
				if i == maxTokenAuthorizedForTest-1 {
					t.Fatalf("max token reached %d, probably stuck in infinite loop...\n"+
						"if you're testing very large input, you may need to bring up the max token limit",
						maxTokenAuthorizedForTest)
				}
			}

			cmp.Diff(tt.wanted, tokensGotFromLexer)
		})
	}
}

// Nothing done : BenchmarkAllocationWithSmallFile-16    	      72	  15531290 ns/op	 2000535 B/op	 1999994 allocs/op
func BenchmarkAllocationWithSmallFile(b *testing.B) {
	smallFile := bytes.Repeat([]byte("mov r1 r0\nMOV EBP EBX ECX\nSUB ebp ebx"), 3)
	lexer := NewLexer()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.ResetWithInput(smallFile)

		for i := 0; i < maxTokenAuthorizedForTest; i++ {
			token := lexer.NextToken()
			var _ = token
		}
	}
}

// nothing done : BenchmarkAllocationWithMediumFile-16    	      72	  15485821 ns/op	 2061101 B/op	 1999000 allocs/op
func BenchmarkAllocationWithMediumFile(b *testing.B) {
	mediumFile := bytes.Repeat([]byte("mov r1 r0\nMOV EBP EBX ECX\nSUB ebp ebx"), 500)
	lexer := NewLexer()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.ResetWithInput(mediumFile)

		for i := 0; i < maxTokenAuthorizedForTest; i++ {
			token := lexer.NextToken()
			var _ = token
		}
	}
}

// nothing done : BenchmarkAllocationWithLargeFile-16    	      34	  32743159 ns/op	13091403 B/op	 1818182 allocs/op
func BenchmarkAllocationWithLargeFile(b *testing.B) {
	largeFile := bytes.Repeat([]byte("mov r1 r0\nMOV EBP EBX ECX\nSUB ebp ebx"), 100000)
	lexer := NewLexer()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.ResetWithInput(largeFile)

		for i := 0; i < maxTokenAuthorizedForTest; i++ {
			token := lexer.NextToken()
			var _ = token
		}
	}
}

// Benchmark to optimize lookupIdentType.
func BenchmarkAllocationLookUpIdentTypeWithLargeFile(b *testing.B) {

	b.Run("lookupIdentType not optimized", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = lookupIdentType([]byte("mov"))
		}
	})

	b.Run("lookupIdentType optimized", func(b *testing.B) {
		initLookupIdentTypeOptimized()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = lookupIdentTypeOptimized([]byte("mov"))
		}
	})
}
