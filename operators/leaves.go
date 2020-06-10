package operators

import "strings"

// Rune defines a single character.
func Rune(key string, r rune) Operator {
	return func(s []rune) Alternatives {
		if len(s) == 0 || s[0] != r {
			return nil
		}
		return []*Node{
			{
				Key:   key,
				Value: s[:1],
			},
		}
	}
}

// String defines a certain sequence of case insensitive characters.
func String(key string, str string) Operator {
	return func(s []rune) Alternatives {
		if len(str) > len(s) ||
			strings.ToLower(string(s[:len(str)])) != strings.ToLower(str) {
			return nil
		}
		return []*Node{
			{
				Key:   key,
				Value: s[:len(str)],
			},
		}
	}
}

// StringCS defines a certain sequence of case sensitive character.
func StringCS(key string, str string) Operator {
	return func(s []rune) Alternatives {
		if len(str) > len(s) || string(s[:len(str)]) != str {
			return nil
		}
		return []*Node{
			{
				Key:   key,
				Value: s[:len(str)],
			},
		}
	}
}

// Range defines the range of alternative numeric values compactly.
func Range(key string, l, h rune) Operator {
	return func(s []rune) Alternatives {
		if len(s) == 0 || s[0] < l || h < s[0] {
			return nil
		}
		return []*Node{
			{
				Key:   key,
				Value: s[:1],
			},
		}
	}
}
