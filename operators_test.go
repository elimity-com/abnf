package abnf

import (
	"strings"
	"testing"
)

func TestRune(t *testing.T) {
	for _, s := range []string{
		"a",
		"aa",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, Rune('a')) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("b", Rune('a')) != nil {
		t.Errorf("value found for \"b\"")
	}
}

func TestRunes(t *testing.T) {
	for _, s := range []string{
		"a",
		"b",
		"cba",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, Runes('a', 'b', 'c')) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("d", Runes('a', 'b', 'c')) != nil {
		t.Errorf("value found for \"d\"")
	}
}

func TestString(t *testing.T) {
	for _, s := range []string{
		"abc",
		"abc abc",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, String("abc")) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("a bc", String("abc")) != nil {
		t.Errorf("value found for \"a bc\"")
	}
}

func TestConcat(t *testing.T) {
	if ParseString("abc", Concat(Rune('a'), Rune('b'), Rune('c'))) == nil {
		t.Errorf("no value found for \"abc\"")
	}

	for _, s := range []string{
		"a",
		"b",
		"cba",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, Concat(Rune('a'), Rune('b'), Rune('c'))) != nil {
				t.Errorf("value found for: %s", s)
			}
		})
	}
}

func TestAlts(t *testing.T) {
	for _, s := range []string{
		"a",
		"b",
		"abc",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, Alts(Rune('a'), Rune('b'))) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("c", Alts(Rune('a'), Rune('b'))) != nil {
		t.Errorf("value found for \"c\"")
	}
}

func TestRange(t *testing.T) {
	for _, s := range []string{
		"a",
		"a&",
		"z",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, Range('a', 'z')) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("&", Range('a', 'z')) != nil {
		t.Errorf("value found for \"&\"")
	}
}

func TestRepeat(t *testing.T) {
	for _, s := range []string{
		"aa",
		"aaa",
		strings.Repeat("a", 99),
	} {
		t.Run("", func(t *testing.T) {
			r := ParseString(s, Repeat(2, 3, Rune('a')))
			if r == nil {
				t.Errorf("no value found for: %s", s)
				return
			}

			if len(r) > 3 {
				t.Errorf("value too long: %s", string(r))
			}
		})
	}

	if ParseString("a", Repeat(2, 3, Rune('a'))) != nil {
		t.Errorf("value found for \"a\"")
	}
}

func TestRepeatN(t *testing.T) {
	for _, s := range []string{
		"aaaaa",
		"aaaaaaaaaaa",
		strings.Repeat("a", 99),
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, RepeatN(5, Rune('a'))) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("aaaa", RepeatN(5, Rune('a'))) != nil {
		t.Errorf("value found for \"aaaa\"")

	}
}

func TestRepeat0Inf(t *testing.T) {
	for _, s := range []string{
		"",
		"b",
		"a",
		"aaa",
		strings.Repeat("a", 99),
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, Repeat0Inf(Rune('a'))) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}
}

func TestRepeat1Inf(t *testing.T) {
	for _, s := range []string{
		"a",
		"aaa",
		"aaaab",
		strings.Repeat("a", 99),
	} {
		t.Run("", func(t *testing.T) {
			r := ParseString(s, Repeat1Inf(Rune('a')))
			if r == nil {
				t.Errorf("no value found for: %s", s)
				return
			}

			for _, a := range string(r) {
				if a != 'a' {
					t.Errorf("value is not an \"a\": %s", string(a))
				}
			}
		})
	}

	for _, s := range []string{
		"",
		"b",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, Repeat1Inf(Rune('a'))) != nil {
				t.Errorf("value found for: %s", s)
			}
		})
	}
}

func TestOptional(t *testing.T) {
	for _, s := range []string{
		"a",
		"",
	} {
		t.Run(s, func(t *testing.T) {
			r := ParseString(s, Optional(Rune('a')))
			if r == nil {
				t.Errorf("no value found for: %s", s)
				return
			}

			str := string(r)
			if str != "a" && str != "" {
				t.Errorf("value does not match empty string or \"a\": %s", str)
			}
		})
	}
}
