package encoding

import (
	"golang.org/x/text/encoding"
	"unicode/utf8"
)

// RFC: https://tools.ietf.org/html/rfc20

// 7-bit ASCII
var ASCII encoding.Encoding

const (
	// https://tools.ietf.org/html/rfc20#section-2
	ASCIISub = byte(0b0011010)
)

func init() {
	runeErrors := make(map[byte]rune)
	for i := 128; i <= 255; i++ {
		runeErrors[byte(i)] = utf8.RuneError
	}
	ASCII = NewCharacterSet(runeErrors, ASCIISub)
}
