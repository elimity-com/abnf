package abnf

import (
	"fmt"
	"strings"
)

type Operator func(s *Scanner) *AST

func Rune(name string, r rune) Operator {
	return func(s *Scanner) *AST {
		if n := s.nextRune(); n != nil && n[0] == r {
			s.pointer++
			return &AST{
				Key:      name,
				Value:    n,
				Children: nil,
			}
		}
		return nil
	}
}

func String(name, str string, caseSensitive bool) Operator {
	if caseSensitive {
		return func(s *Scanner) *AST {
			rules := make([]Operator, len(str))
			for i, r := range str {
				rules[i] = Rune(string(r), r)
			}
			return Concat(name, rules...)(s)
		}
	}
	return func(s *Scanner) *AST {
		str = strings.ToLower(str)
		rules := make([]Operator, len(str))
		for i, r := range str {
			tmp := r
			if '\x61' <= tmp && tmp <= '\x7A' {
				tmp -= '\x20'
				rules[i] = Alts(
					fmt.Sprintf("%s / %s", string(tmp), string(r)),
					Rune(string(tmp), tmp),
					Rune(string(r), r),
				)
			} else {
				rules[i] = Rune(string(r), r)
			}
		}
		return Concat(name, rules...)(s)
	}
}

func Concat(name string, r ...Operator) Operator {
	return func(s *Scanner) *AST {
		s.addMarker()
		children := make([]AST, len(r))
		for i, rule := range r {
			n := rule(s)
			if n == nil {
				s.rollbackMarker()
				return nil
			}
			children[i] = *n
		}
		return &AST{
			Key:      name,
			Value:    s.commitValue(),
			Children: children,
		}
	}
}

func Alts(name string, r ...Operator) Operator {
	return func(s *Scanner) *AST {
		var alt *AST // the (longest) best match
		var size int // the length of the raw result in alt
		for _, rule := range r {
			n := rule(s)
			if n == nil {
				continue
			}
			if s := len(n.Value); s > size {
				alt = n
				size = s
			}
		}
		if alt != nil {
			return &AST{
				Key:      name,
				Value:    alt.Value,
				Children: []AST{*alt},
			}
		}
		return nil
	}
}

func Range(name string, l, h rune) Operator {
	return func(s *Scanner) *AST {
		if r := s.nextRune(); r != nil && l <= r[0] && r[0] <= h {
			s.pointer++
			return &AST{
				Key:      name,
				Value:    r,
				Children: nil,
			}
		}
		return nil
	}
}

func Repeat(name string, min, max int, r Operator) Operator {
	return func(s *Scanner) *AST {
		return repeatRule(name, s, min, max, r)
	}
}

func RepeatN(name string, n int, r Operator) Operator {
	return func(s *Scanner) *AST {
		return repeatRule(name, s, n, n, r)
	}
}

func Repeat0Inf(name string, r Operator) Operator {
	return func(s *Scanner) *AST {
		return repeatRule(name, s, 0, -1, r)
	}
}

func Repeat1Inf(name string, r Operator) Operator {
	return func(s *Scanner) *AST {
		return repeatRule(name, s, 1, -1, r)
	}
}

func Optional(name string, r Operator) Operator {
	return func(s *Scanner) *AST {
		n := r(s)
		if n == nil {
			return &AST{
				Key:      name,
				Value:    nil,
				Children: nil,
			}
		}
		return &AST{
			Key:      name,
			Value:    n.Value,
			Children: []AST{*n},
		}
	}
}

func repeatRule(name string, s *Scanner, min, max int, r Operator) *AST {
	s.addMarker()
	children := make([]AST, 0)
	var i int
	for max < 0 || i < max {
		n := r(s)
		if n == nil {
			break
		}
		children = append(children, *n)
		i++
	}
	if i < min {
		s.rollbackMarker()
		return nil
	}
	return &AST{
		Key:      name,
		Value:    s.commitValue(),
		Children: children,
	}
}
