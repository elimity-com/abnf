package encoding

import (
	"testing"
	"unicode/utf8"
)

func TestASCII(t *testing.T) {
	for i := 0; i <= 255; i++ {
		r := rune(i)
		b := byte(i)
		if 128 <= i {
			b = ASCIISub
		}
		e := ASCII.NewEncoder()

		dst := make([]byte, 6)
		src := make([]byte, utf8.RuneLen(r))
		utf8.EncodeRune(src, r)

		nDst, nSrc, err := e.Transform(dst, src, true)
		if err != nil {
			t.Error(err)
		}
		if nSrc != len(src) {
			t.Errorf("incorrect length: %d, %d", nSrc, len(src))
		}
		if nDst != 1 {
			t.Errorf("length not one: %d", nDst)
		}
		if b != dst[0] {
			t.Errorf("invalid character: %d", dst[0])
		}
	}
}
