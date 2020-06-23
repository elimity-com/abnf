package operators

import (
	"fmt"
	"testing"
)

func TestOptional(t *testing.T) {
	rule := Optional(`[ a ]`, a)
	for i, s := range []string{
		"a",
		"",
	} {
		t.Run(fmt.Sprintf("Simple %d", i), func(t *testing.T) {
			nodes := rule([]byte(s))
			if len(nodes) == 0 {
				t.Errorf("no value found for: %s", s)
				return
			}

			for _, node := range nodes {
				str := string(node.Value)
				if str != "a" && str != "" {
					t.Errorf("value does not match empty string or \"a\": %s", str)
				}
			}
		})
	}

	t.Run("Complex", func(t *testing.T) {
		rule := Concat(`[ *1( a ":" ) a ] "::"`,
			Optional(`[ *1( a ":" ) a ]`,
				Concat(`*1( a ":" ) a`,
					Repeat(`*1( a ":" )`, 0, 1,
						Concat(`a ":"`,
							Terminal(`a`, []byte("a")),
							Terminal(`:`, []byte(":")),
						),
					),
					Terminal(`a`, []byte("a")),
				),
			),
			String(`::`, "::"),
		)

		for _, s := range []string{
			"::",
			"a::",
			"a:a::",
		} {
			nodes := rule([]byte(s))
			if len(nodes) == 0 {
				t.Error("no value found")
			}
		}
	})
}
