package operators

import (
	"fmt"
	"testing"
)

var (
	a = Terminal(`a`, []byte("a"))
)

func TestRune(t *testing.T) {
	for i, s := range []string{
		"a",
		"aa",
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			if len(a([]byte(s))) == 0 {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if len(a([]byte("b"))) != 0 {
		t.Errorf("value found for \"b\"")
	}
}

func TestString(t *testing.T) {
	rule := StringCI(`abc`, "abc")
	for i, s := range []string{
		"abc",
		"aBc",
		"abc abc",
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			if len(rule([]byte(s))) == 0 {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if len(String(`abc`, "abc")([]byte("aBc"))) != 0 {
		t.Errorf("value found for \"aBc\"")
	}

	if rule([]byte("a bc")) != nil {
		t.Errorf("value found for \"a bc\"")
	}
}

func TestRange(t *testing.T) {
	rule := Range(`a-z`, []byte("a"), []byte("z"))
	for i, s := range []string{
		"a",
		"a&",
		"z",
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			if len(rule([]byte(s))) == 0 {
				t.Errorf("no value found for: %s", s)
			}
		})
	}

	if len(rule([]byte("&"))) != 0 {
		t.Errorf("value found for \"&\"")
	}

	if len(Range("%x5D-10FFFF", []byte{93}, []byte{16, 255, 255})([]byte("x"))) == 0 {
		t.Error("no value found for \"x\"")
	}
}
