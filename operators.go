package abnf

import (
	"strings"
)

type Operator func([]rune) *AST

func Rune(name string, r rune) Operator {
	return func(s []rune) *AST {
		if len(s) == 0 {
			return nil
		}
		if s[0] == r {
			return &AST{
				Key:   name,
				Value: s[:1],
			}
		}
		return nil
	}
}

func String(name, str string, caseSensitive bool) Operator {
	if caseSensitive {
		return func(s []rune) *AST {
			if len(s) < len(str) {
				return nil
			}
			if string(s[:len(str)]) == str {
				return &AST{
					Key:   name,
					Value: s[:len(str)],
				}
			}
			return nil
		}
	}
	return func(s []rune) *AST {
		if len(str) > len(s) {
			return nil
		}
		if strings.ToLower(string(s[:len(str)])) == strings.ToLower(str) {
			return &AST{
				Key:   name,
				Value: s[:len(str)],
			}
		}
		return nil
	}
}

func Concat(name string, r ...Operator) Operator {
	return func(s []rune) *AST {
		children := make([]AST, len(r))
		tmp, l := s, 0
		for i, rule := range r {
			child := rule(tmp)
			if child == nil {
				return nil
			}
			children[i] = *child
			l += len(child.Value)
			tmp = s[l:]
		}
		return &AST{
			Key:      name,
			Value:    s[:l],
			Children: children,
		}
	}
}

func Alts(name string, r ...Operator) Operator {
	return func(s []rune) *AST {
		var alt *AST // the (longest) best match
		var size int // the length of the raw result in alt
		for _, rule := range r {
			n := rule(s)
			if n == nil {
				continue
			}
			if l := len(n.Value); size < l {
				alt = n
				size = l
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
	return func(s []rune) *AST {
		if len(s) == 0 {
			return nil
		}
		if l <= s[0] && s[0] <= h {
			return &AST{
				Key:   name,
				Value: s[:1],
			}
		}
		return nil
	}
}

func Repeat(name string, min, max int, r Operator) Operator {
	return func(s []rune) *AST {
		return repeatRule(name, s, min, max, r)
	}
}

func RepeatN(name string, n int, r Operator) Operator {
	return func(s []rune) *AST {
		return repeatRule(name, s, n, n, r)
	}
}

func Repeat0Inf(name string, r Operator) Operator {
	return func(s []rune) *AST {
		return repeatRule(name, s, 0, -1, r)
	}
}

func Repeat1Inf(name string, r Operator) Operator {
	return func(s []rune) *AST {
		return repeatRule(name, s, 1, -1, r)
	}
}

func Optional(name string, r Operator) Operator {
	return func(s []rune) *AST {
		n := r(s)
		if n == nil {
			return &AST{
				Key:   name,
				Value: s[:0],
			}
		}
		return &AST{
			Key:      name,
			Value:    n.Value,
			Children: []AST{*n},
		}
	}
}

func repeatRule(name string, s []rune, min, max int, r Operator) *AST {
	children := make([]AST, 0)
	tmp, l, i := s, 0, 0
	for max < 0 || i < max {
		n := r(tmp)
		if n == nil {
			break
		}
		children = append(children, *n)
		l += len(n.Value)
		tmp = s[l:]
		i++
	}
	if i < min {
		return nil
	}
	return &AST{
		Key:      name,
		Value:    s[:l],
		Children: children,
	}
}
