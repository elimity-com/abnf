package abnf

import (
	"bytes"
	"io"
)

func ParseString(s string, rule Rule) []rune {
	scanner := &scanner{
		main:    make([]rune, 0, len(s)),
		r:       bytes.NewReader([]byte(s)),
		markers: make([]int, 0, len(s)),
	}

	if value := rule(scanner); value != nil {
		return value
	}
	return nil
}

type scanner struct {
	main    []rune
	r       io.RuneReader
	markers []int
	pointer int
}

func (s *scanner) addMarker() {
	s.markers = append(s.markers, s.pointer) // add marker
}

func (s *scanner) rollbackMarker() {
	s.pointer = s.markers[len(s.markers)-1]  // assign pointer
	s.markers = s.markers[:len(s.markers)-1] // pop marker
}

func (s *scanner) commitValue() []rune {
	r := s.main[s.markers[len(s.markers)-1]:s.pointer] // load runes from marker until pointer
	s.markers = s.markers[:len(s.markers)-1]           // pop marker
	return r
}

func (s *scanner) nextRune() []rune {
	if len(s.main) <= s.pointer {
		r, _, err := s.r.ReadRune()
		if err != nil {
			return nil
		}
		s.main = append(s.main, r)
	}
	return s.main[s.pointer : s.pointer+1]
}

type Rule func(s *scanner) []rune

func Rune(r rune) Rule {
	return func(s *scanner) []rune {
		if n := s.nextRune(); n != nil && n[0] == r {
			s.pointer++
			return n
		}
		return nil
	}
}

func Runes(rs ...rune) Rule {
	return func(s *scanner) []rune {
		n := s.nextRune()
		if n == nil {
			return nil
		}
		for _, r := range rs {
			if n[0] == r {
				s.pointer++
				return n
			}
		}
		return nil
	}
}

func String(s string) Rule {
	rules := make([]Rule, len(s))
	for i, r := range s {
		rules[i] = Rune(r)
	}
	return Concat(rules...)
}

func Concat(r ...Rule) Rule {
	return func(s *scanner) []rune {
		s.addMarker()
		for _, rule := range r {
			if rule(s) == nil {
				s.rollbackMarker()
				return nil
			}
		}
		return s.commitValue()
	}
}

func Alts(r ...Rule) Rule {
	return func(s *scanner) []rune {
		for _, rule := range r {
			if alt := rule(s); alt != nil {
				return alt
			}
		}
		return nil
	}
}

func Range(l, h rune) Rule {
	return func(s *scanner) []rune {
		if r := s.nextRune(); r != nil && l <= r[0] && r[0] <= h {
			s.pointer++
			return r
		}
		return nil
	}
}

func DefaultRepeat(r Rule) Rule {
	return func(s *scanner) []rune {
		return repeat(s, 0, -1, r)
	}
}

func VarRepeat(min, max int, r Rule) Rule {
	return func(s *scanner) []rune {
		return repeat(s, min, max, r)
	}
}

func Repeat(n int, r Rule) Rule {
	return func(s *scanner) []rune {
		return repeat(s, n, n, r)
	}
}

func Optional(r Rule) Rule {
	return func(s *scanner) []rune {
		opt := r(s)
		if opt == nil {
			opt = make([]rune, 0)
		}
		return opt
	}
}

func repeat(s *scanner, min, max int, r Rule) []rune {
	s.addMarker()
	var i int
	for max < 0 || i < max {
		if r(s) == nil {
			break
		}
		i++
	}
	if i < min {
		s.rollbackMarker()
		return nil
	}
	return s.commitValue()
}
