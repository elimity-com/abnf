package abnf

import (
	"fmt"
	"strings"
	"testing"
)

func TestAltsAndConcats(t *testing.T) {
	for _, test := range []struct {
		abnf     string
		contains string
	}{
		{
			abnf: "x = x x / x x\n",
			contains: `operators.Alts(
		"x",
		operators.Concat(
			"x x",
			X,
			X,
		),
		operators.Concat(
			"x x",
			X,
			X,
		),`,
		},
		{
			abnf: "x = x x / x ( x / x )\n",
			contains: `operators.Alts(
		"x",
		operators.Concat(
			"x x",
			X,
			X,
		),
		operators.Concat(
			"x ( x / x )",
			X,
			operators.Alts(
				"( x / x )",
				X,
				X,
			),
		),`,
		},
		{
			abnf: "x = x / ( x / x ) ( x x / x )\n",
			contains: `operators.Alts(
		"x",
		X,
		operators.Concat(
			"( x / x ) ( x x / x )",
			operators.Alts(
				"( x / x )",
				X,
				X,
			),
			operators.Alts(
				"( x x / x )",
				operators.Concat(
					"x x",
					X,
					X,
				),
				X,
			),
		),`,
		},
	} {
		t.Run("", func(t *testing.T) {
			g := Generator{
				PackageName:  "alts",
				RawABNF:      test.abnf,
			}
			f := g.GenerateABNFAsAlternatives()
			if !strings.Contains(fmt.Sprintf("%#v", f), test.contains) {
				t.Errorf("did not parse correctly")
			}
		})
	}
}
