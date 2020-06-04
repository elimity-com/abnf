package operators

import (
	"fmt"
	"testing"
)

var (
	a = Rune(`a`, 'a')
)

func TestRune(t *testing.T) {
	for i, s := range []string{
		"a",
		"aa",
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			if len(a([]rune(s))) == 0 {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if len(a([]rune("b"))) != 0 {
		t.Errorf("value found for \"b\"")
	}
}

func TestString(t *testing.T) {
	rule := String(`abc`, "abc")
	for i, s := range []string{
		"abc",
		"aBc",
		"abc abc",
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			if len(rule([]rune(s))) == 0 {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if len(StringCS(`abc`, "abc")([]rune("aBc"))) != 0 {
		t.Errorf("value found for \"aBc\"")
	}

	if rule([]rune("a bc")) != nil {
		t.Errorf("value found for \"a bc\"")
	}
}

func TestRange(t *testing.T) {
	rule := Range(`a-z`, 'a', 'z')
	for i, s := range []string{
		"a",
		"a&",
		"z",
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			if len(rule([]rune(s))) == 0 {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if len(rule([]rune("&"))) != 0 {
		t.Errorf("value found for \"&\"")
	}
}
