package abnf

import (
	"strings"
	"testing"
)

var (
	a = Rune(`a`, 'a')
	b = Rune(`b`, 'b')
	c = Rune(`c`, 'c')
)

func TestRune(t *testing.T) {
	for _, s := range []string{
		"a",
		"aa",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, a) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("b", a) != nil {
		t.Errorf("value found for \"b\"")
	}
}

func TestRunes(t *testing.T) {
	rule := Runes(`a / b / c`, 'a', 'b', 'c')
	for _, s := range []string{
		"a",
		"b",
		"cba",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, rule) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("d", rule) != nil {
		t.Errorf("value found for \"d\"")
	}
}

func TestString(t *testing.T) {
	rule := String(`abc`, "abc")
	for _, s := range []string{
		"abc",
		"abc abc",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, rule) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("a bc", rule) != nil {
		t.Errorf("value found for \"a bc\"")
	}
}

func TestConcat(t *testing.T) {
	rule := Concat(`abc`, a, b, c)
	if ParseString("abc", rule) == nil {
		t.Errorf("no value found for \"abc\"")
	}

	for _, s := range []string{
		"a",
		"b",
		"cba",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, rule) != nil {
				t.Errorf("value found for: %s", s)
			}
		})
	}
}

func TestAlts(t *testing.T) {
	rule := Alts(`a / b`, a, b)
	for _, s := range []string{
		"a",
		"b",
		"abc",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, rule) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("c", rule) != nil {
		t.Errorf("value found for \"c\"")
	}
}

func TestRange(t *testing.T) {
	rule := Range(`a-z`, 'a', 'z')
	for _, s := range []string{
		"a",
		"a&",
		"z",
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, rule) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("&", rule) != nil {
		t.Errorf("value found for \"&\"")
	}
}

func TestRepeat(t *testing.T) {
	rule := Repeat(`2*3a`, 2, 3, a)
	for _, s := range []string{
		"aa",
		"aaa",
		strings.Repeat("a", 99),
	} {
		t.Run("", func(t *testing.T) {
			r := ParseString(s, rule)
			if r == nil {
				t.Errorf("no value found for: %s", s)
				return
			}

			if len(r.Raw) > 3 {
				t.Errorf("value too long: %s", string(r.Raw))
			}
		})
	}

	if ParseString("a", rule) != nil {
		t.Errorf("value found for \"a\"")
	}
}

func TestRepeatN(t *testing.T) {
	rule := RepeatN(`5a`, 5, a)
	for _, s := range []string{
		"aaaaa",
		"aaaaaaaaaaa",
		strings.Repeat("a", 99),
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, rule) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if ParseString("aaaa", rule) != nil {
		t.Errorf("value found for \"aaaa\"")

	}
}

func TestRepeat0Inf(t *testing.T) {
	rule := Repeat0Inf(`*a`, a)
	for _, s := range []string{
		"",
		"b",
		"a",
		"aaa",
		strings.Repeat("a", 99),
	} {
		t.Run("", func(t *testing.T) {
			if ParseString(s, rule) == nil {
				t.Errorf("no value found for: %s", s)
			}
		})
	}
}

func TestRepeat1Inf(t *testing.T) {
	rule := Repeat1Inf(`1*a`, a)
	for _, s := range []string{
		"a",
		"aaa",
		"aaaab",
		strings.Repeat("a", 99),
	} {
		t.Run("", func(t *testing.T) {
			r := ParseString(s, rule)
			if r == nil {
				t.Errorf("no value found for: %s", s)
				return
			}

			for _, a := range string(r.Raw) {
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
			if ParseString(s, rule) != nil {
				t.Errorf("value found for: %s", s)
			}
		})
	}
}

func TestOptional(t *testing.T) {
	rule := Optional(`[a]`, a)
	for _, s := range []string{
		"a",
		"",
	} {
		t.Run(s, func(t *testing.T) {
			r := ParseString(s, rule)
			if r == nil {
				t.Errorf("no value found for: %s", s)
				return
			}

			str := string(r.Raw)
			if str != "a" && str != "" {
				t.Errorf("value does not match empty string or \"a\": %s", str)
			}
		})
	}
}
