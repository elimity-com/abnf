package abnf

import (
	"bytes"
	"fmt"
	"io"
)

type ABNF struct {
	rules map[string]Operator
	s     *Scanner
}

func New(rules map[string]Operator, s string) *ABNF {
	return &ABNF{
		rules: rules,
		s:     NewScanner(s),
	}
}

func (abnf *ABNF) Get(key string) (string, error) {
	rule, ok := abnf.rules[key]
	if !ok {
		return "", fmt.Errorf("could not find rule with name \"%s\"", key)
	}
	runes := rule(abnf.s)
	if runes == nil {
		return "", fmt.Errorf("could not find string matching rule \"%s\"", key)
	}
	return string(runes), nil
}

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

func ParseString(s string, rule Operator) []rune {
	scanner := &Scanner{
		main:    make([]rune, 0, len(s)),
		r:       bytes.NewReader([]byte(s)),
		markers: make([]int, 0, len(s)),
	}

	if value := rule(scanner); value != nil {
		return value
	}
	return nil
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

type Operator func(s *Scanner) []rune

func Rune(r rune) Operator {
	return func(s *Scanner) []rune {
		if n := s.nextRune(); n != nil && n[0] == r {
			s.pointer++
			return n
		}
		return nil
	}
}

func Runes(rs ...rune) Operator {
	return func(s *Scanner) []rune {
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

func String(s string) Operator {
	rules := make([]Operator, len(s))
	for i, r := range s {
		rules[i] = Rune(r)
	}
	return Concat(rules...)
}

func Concat(r ...Operator) Operator {
	return func(s *Scanner) []rune {
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

func Alts(r ...Operator) Operator {
	return func(s *Scanner) []rune {
		for _, rule := range r {
			if alt := rule(s); alt != nil {
				return alt
			}
		}
		return nil
	}
}

func Range(l, h rune) Operator {
	return func(s *Scanner) []rune {
		if r := s.nextRune(); r != nil && l <= r[0] && r[0] <= h {
			s.pointer++
			return r
		}
		return nil
	}
}

func Repeat(min, max int, r Operator) Operator {
	return func(s *Scanner) []rune {
		return repeatRule(s, min, max, r)
	}
}

func RepeatN(n int, r Operator) Operator {
	return func(s *Scanner) []rune {
		return repeatRule(s, n, n, r)
	}
}

func Repeat0Inf(r Operator) Operator {
	return func(s *Scanner) []rune {
		return repeatRule(s, 0, -1, r)
	}
}

func Repeat1Inf(r Operator) Operator {
	return func(s *Scanner) []rune {
		return repeatRule(s, 1, -1, r)
	}
}

func Optional(r Operator) Operator {
	return func(s *Scanner) []rune {
		opt := r(s)
		if opt == nil {
			opt = make([]rune, 0)
		}
		return opt
	}
}

func repeatRule(s *Scanner, min, max int, r Operator) []rune {
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
