package operators

import (
	"bytes"
	"strings"
)

// Terminal defines a single character.
func Terminal(key string, value []byte) Operator {
	return func(s []byte) Alternatives {
		if len(s) < len(value) || bytes.Compare(s[:len(value)], value) != 0 {
			return nil
		}
		return []*Node{
			{
				Key:   key,
				Value: s[:len(value)],
			},
		}
	}
}

// String defines a certain sequence of case sensitive characters.
func String(key string, str string) Operator {
	return func(s []byte) Alternatives {
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

// StringCS defines a certain sequence of case insensitive character.
func StringCI(key string, str string) Operator {
	return func(s []byte) Alternatives {
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

// Range defines the range of alternative numeric values compactly.
func Range(key string, l, h []byte) Operator {
	return func(s []byte) Alternatives {
		if len(s) == 0 || len(s) < len(l) || bytes.Compare(s[:len(l)], l) < 0 {
			return nil
		}

		var l int
		for i := range h {
			if i+1 <= i+1 && bytes.Compare(s[:i+1], h) <= 0 {
				l++
			} else {
				break
			}
		}

		if l == 0 {
			return nil
		}

		return []*Node{
			{
				Key:   key,
				Value: s[:l],
			},
		}
	}
}
