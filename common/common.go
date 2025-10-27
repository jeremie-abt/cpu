package common

import (
	"bytes"
	"strconv"
	"unsafe"
)

func IsLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsHexDigit(ch byte) bool {
	return IsDigit(ch) || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func ParseNumber(literal []byte) int64 {
	if bytes.HasPrefix(literal, []byte("0x")) || bytes.HasPrefix(literal, []byte("0X")) {
		// TODO: Je pense que si mon fichier fini par 0X ca crash, détecter ca avec le fuzz.
		literalAsStr := unsafe.String(unsafe.SliceData(literal[2:]), len(literal[2:]))
		val, _ := strconv.ParseInt(literalAsStr, 16, 64)
		return val
	}

	if bytes.HasPrefix(literal, []byte("0b")) || bytes.HasPrefix(literal, []byte("0B")) {
		literalAsStr := unsafe.String(unsafe.SliceData(literal[2:]), len(literal[2:]))
		val, _ := strconv.ParseInt(literalAsStr, 2, 64)
		return val
	}

	literalAsStr := unsafe.String(unsafe.SliceData(literal), len(literal))
	val, _ := strconv.ParseInt(literalAsStr, 10, 64)
	return val
}
