package abnf

import (
	"bytes"
	"io"
)

type Scanner struct {
	main    []rune
	r       io.RuneReader
	markers []int
	pointer int
}

func NewScanner(s string) *Scanner {
	return &Scanner{
		main:    make([]rune, 0, len(s)),
		r:       bytes.NewReader([]byte(s)),
		markers: make([]int, 0, len(s)),
	}
}

func ParseString(s string, rule Operator) *AST {
	return rule(NewScanner(s))
}

func (s *Scanner) addMarker() {
	s.markers = append(s.markers, s.pointer) // add marker
}

func (s *Scanner) rollbackMarker() {
	s.pointer = s.markers[len(s.markers)-1]  // assign pointer
	s.markers = s.markers[:len(s.markers)-1] // pop marker
}

func (s *Scanner) commitValue() []rune {
	r := s.main[s.markers[len(s.markers)-1]:s.pointer] // load runes from marker until pointer
	s.markers = s.markers[:len(s.markers)-1]           // pop marker
	return r
}

func (s *Scanner) nextRune() []rune {
	if len(s.main) <= s.pointer {
		r, _, err := s.r.ReadRune()
		if err != nil {
			return nil
		}
		s.main = append(s.main, r)
	}
	return s.main[s.pointer : s.pointer+1]
}
