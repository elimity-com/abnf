package core

import (
	"testing"

	. "github.com/di-wu/abnf/operators"
	"github.com/di-wu/regen"
)

func TestCore(t *testing.T) {
	for _, test := range []struct {
		name                     string
		validRegex, invalidRegex string
		rule                     Operator
		allowsEmpty              bool
	}{
		{
			name:         "ALPHA",
			validRegex:   `[a-zA-Z]`,
			invalidRegex: `[^a-zA-Z]`,
			rule:         ALPHA(),
		},
		{
			name:         "BIT",
			validRegex:   `[0-1]`,
			invalidRegex: `[^0-1]`,
			rule:         BIT(),
		},
		{
			name:         "CHAR",
			validRegex:   `[\x01-\x7A]`,
			invalidRegex: `\x00`,
			rule:         CHAR(),
		},
		{
			name:         "CR",
			validRegex:   `\x0D`,
			invalidRegex: `[^\x0D]`,
			rule:         CR(),
		},
		{
			name:         "CRLF",
			validRegex:   `\x0D\x0A`,
			invalidRegex: `[^\x0D\x0A]`,
			rule:         CRLF(),
		},
		{
			name:         "CTL",
			validRegex:   `[\x00-\x1F]|\x7F`,
			invalidRegex: `[^\x00-\x1F]|[^\x7F]`,
			rule:         CTL(),
		},
		{
			name:         "DIGIT",
			validRegex:   `\d`,
			invalidRegex: `\D`,
			rule:         DIGIT(),
		},
		{
			name:         "DQUOTE",
			validRegex:   `\x22`,
			invalidRegex: `[^\x22]`,
			rule:         DQUOTE(),
		},
		{
			name:         "HEXDIG",
			validRegex:   `\d|[A-F]`,
			invalidRegex: `[^\d]|[^A-F]`,
			rule:         HEXDIG(),
		},
		{
			name:         "HTAB",
			validRegex:   `\x09`,
			invalidRegex: `[^\x09]`,
			rule:         HTAB(),
		},
		{
			name:         "LF",
			validRegex:   `\x0A`,
			invalidRegex: `[^\x0A]`,
			rule:         LF(),
		},
		{
			name:         "LWSP",
			validRegex:   `((\x0D\x0A)?(\x20|\x09))*`,
			invalidRegex: `[a-zA-Z]`, // difficult to specify
			rule:         LWSP(),
			allowsEmpty:  true,
		},
		{
			name:         "OCTET",
			validRegex:   `[\x00-\xFF]`,
			invalidRegex: `[^\x00-\xFF]`,
			rule:         OCTET(),
		},
		{
			name:         "SP",
			validRegex:   `\x20`,
			invalidRegex: `[^\x20]`,
			rule:         SP(),
		},
		{
			name:         "VCHAR",
			validRegex:   `[\x21-\x7E]`,
			invalidRegex: `[^\x21-\x7E]`,
			rule:         VCHAR(),
		},
		{
			name:       "WSP",
			validRegex: `\x20|\x09`,
			rule:       WSP(),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			valid, _ := regen.New(test.validRegex)
			invalid, _ := regen.New(test.invalidRegex)

			for i := 0; i < 1000; i++ {
				validStr := valid.Generate()
				if nodes := test.rule([]rune(validStr)); nodes == nil {
					t.Errorf("no value found for: %s", validStr)
				} else {
					if best := nodes.Best(); !compareRunes(string(best.Value), validStr) {
						t.Errorf("values do not match: %s %s", string(best.Value), validStr)
					}
				}

				invalidStr := invalid.Generate()
				if nodes := test.rule([]rune(invalidStr)); len(nodes) != 0 {
					if test.allowsEmpty {
						for _, node := range nodes {
							if node.String() != "" {
								t.Errorf("value found for %s: %s", invalidStr, node)
							}
						}
					}
				}
			}
		})
	}
}

func compareRunes(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
