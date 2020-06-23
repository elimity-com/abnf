package encoding

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type CharacterSet struct {
	bytes map[rune]byte
	runes [256][]byte

	// byte value for substitution
	sub byte
}

func NewCharacterSet(m map[byte]rune, sub byte) *CharacterSet {
	set := new(CharacterSet)
	set.sub = sub
	set.bytes = make(map[rune]byte)
	for i := 0; i < 256; i++ {
		// check if i is in map
		r, ok := m[byte(i)]
		if !ok {
			r = rune(i)
		}
		// check for rune error
		if r != utf8.RuneError {
			set.bytes[r] = byte(i)
		}

		utf := make([]byte, utf8.RuneLen(r))
		utf8.EncodeRune(utf, r)
		set.runes[i] = utf
	}
	return set
}

func (set *CharacterSet) NewDecoder() *encoding.Decoder {
	return &encoding.Decoder{
		Transformer: &decoder{
			runes: set.runes,
		},
	}
}

type decoder struct {
	runes [256][]byte
}

func (d *decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for _, character := range src {
		b := d.runes[character]
		length := len(b)
		// check size destination buffer
		if len(dst) < nDst+length {
			err = transform.ErrShortDst
			break
		}
		for i := 0; i < length; i++ {
			dst[nDst] = b[i]
			nDst++
		}
	}
	return
}

func (d *decoder) Reset() {}

func (set *CharacterSet) NewEncoder() *encoding.Encoder {
	return &encoding.Encoder{
		Transformer: &encoder{
			bytes:   set.bytes,
			replace: set.sub,
		},
	}
}

type encoder struct {
	bytes   map[rune]byte
	replace byte
}

func (e *encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nSrc < len(src) {
		// check size destination buffer
		if len(dst) <= nDst {
			err = transform.ErrShortDst
		}
		r, size := utf8.DecodeRune(src[nSrc:])
		if r == utf8.RuneError && size == 1 {
			// insufficient data in source
			if !atEOF && !utf8.FullRune(src[nSrc:]) {
				err = transform.ErrShortSrc
				break
			}
		}

		c, ok := e.bytes[r]
		if ok {
			dst[nDst] = c
		} else {
			dst[nDst] = e.replace
		}

		nSrc += size
		nDst++
	}
	return
}

func (e *encoder) Reset() {}
