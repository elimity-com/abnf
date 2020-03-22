package abnf

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
