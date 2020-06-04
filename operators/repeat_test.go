package operators

import (
	"fmt"
	"strings"
	"testing"
)

func TestRepeat0to1(t *testing.T) {
	rule := Repeat(`*1( a )`, 0, 1, a)
	for i, s := range []string{
		"",
		"a",
		strings.Repeat("a", 99),
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			nodes := rule([]rune(s))
			if len(nodes) == 0 {
				t.Errorf("no value found for: %s", s)
				return
			}

			for _, node := range nodes {
				l := len(node.Value)
				if 1 < l {
					t.Error("invalid length")
				}
			}
		})
	}
}

func TestRepeat2to3(t *testing.T) {
	rule := Repeat(`2*3( a )`, 2, 3, a)
	for i, s := range []string{
		"aa",
		"aaa",
		strings.Repeat("a", 99),
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			nodes := rule([]rune(s))
			if len(nodes) == 0 {
				t.Errorf("no value found for: %s", s)
				return
			}

			for _, node := range nodes {
				l := len(node.Value)
				if l < 2 || 3 < l {
					t.Error("invalid length")
				}
			}
		})
	}

	if rule([]rune("a")) != nil {
		t.Errorf("value found for \"a\"")
	}
}

func TestRepeat0toInf(t *testing.T) {
	rule := Repeat(`*( a )`, 0, -1, a)
	for i, s := range []string{
		"aaa",
		"aaaaa",
		strings.Repeat("a", 99),
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			nodes := rule([]rune(s))
			if len(nodes) == 0 {
				t.Errorf("no value found for: %s", s)
				return
			}

			if len(nodes) != len(s)+1 {
				t.Errorf("not enough permutations found: %d", len(nodes))
			}

			var maxLen bool
			for _, node := range nodes {
				if len(node.Value) == len(s) {
					maxLen = true
				}
			}
			if !maxLen {
				t.Error("could not find node with max possible length")
			}
		})
	}

	t.Run(`Complex`, func(t *testing.T) {
		nodes := Repeat0Inf(`*( [ a ] )`, a)([]rune(""))
		if len(nodes) > 1 {
			t.Error("too much nodes found")
		}
	})
}

func TestRepeat1toInf(t *testing.T) {
	rule := Repeat(`*( a )`, 1, -1, a)
	for i, s := range []string{
		"aaa",
		"aaaaa",
		strings.Repeat("a", 99),
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			nodes := rule([]rune(s))
			if len(nodes) == 0 {
				t.Errorf("no value found for: %s", s)
				return
			}

			if len(nodes) != len(s) {
				t.Errorf("not enough permutations found: %d", len(nodes))
			}

			var maxLen bool
			for _, node := range nodes {
				if len(node.Value) == len(s) {
					maxLen = true
				}
			}
			if !maxLen {
				t.Error("could not find node with max possible length")
			}
		})
	}
}

func TestRepeat5(t *testing.T) {
	nodes := Repeat(`5( a )`, 5, 5, a)([]rune("aaaaa"))
	if len(nodes) == 0 {
		t.Errorf("no value found for aaaaa")
		return
	}

	for _, node := range nodes {
		if len(node.Value) != 5 {
			t.Error("invalid length")
		}
	}
}

func TestRepeatChildren(t *testing.T) {
	for i, rule := range []Operator{
		Repeat(`*2( a )`, 0, 2, a),
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			nodes := rule([]rune("aa"))
			for _, node := range nodes {
				if len(node.Children) != len(string(node.Value)) {
					t.Error("not enough children")
				}
			}
		})
	}
}
