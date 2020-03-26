package abnf

import (
	"fmt"
	"github.com/di-wu/regen"
	"testing"
)

func TestCore(t *testing.T) {
	for _, test := range []struct {
		name                     string
		validRegex, invalidRegex string
		rule                     Operator
	}{
		{
			name:         "ALPHA",
			validRegex:   `[a-zA-Z]`,
			invalidRegex: `[^a-zA-Z]`,
			rule:         alpha(),
		},
		{
			name:         "BIT",
			validRegex:   `[0-1]`,
			invalidRegex: `[^0-1]`,
			rule:         bit(),
		},
		{
			name:         "CHAR",
			validRegex:   `[\x01-\x7A]`,
			invalidRegex: `\x00`,
			rule:         char(),
		},
		{
			name:         "CR",
			validRegex:   `\x0D`,
			invalidRegex: `[^\x0D]`,
			rule:         cr(),
		},
		{
			name:         "CRLF",
			validRegex:   `\x0D\x0A`,
			invalidRegex: `[^\x0D\x0A]`,
			rule:         crlf(),
		},
		{
			name:         "CTL",
			validRegex:   `[\x00-\x1F]|\x7F`,
			invalidRegex: `[^\x00-\x1F]|[^\x7F]`,
			rule:         ctl(),
		},
		{
			name:         "DIGIT",
			validRegex:   `\d`,
			invalidRegex: `\D`,
			rule:         digit(),
		},
		{
			name:         "DQUOTE",
			validRegex:   `\x22`,
			invalidRegex: `[^\x22]`,
			rule:         dquote(),
		},
		{
			name:         "HEXDIG",
			validRegex:   `\d|[A-F]`,
			invalidRegex: `[^\d]|[^A-F]`,
			rule:         hexdig(),
		},
		{
			name:         "HTAB",
			validRegex:   `\x09`,
			invalidRegex: `[^\x09]`,
			rule:         htab(),
		},
		{
			name:         "LF",
			validRegex:   `\x0A`,
			invalidRegex: `[^\x0A]`,
			rule:         lf(),
		},
		{
			name:         "LWSP",
			validRegex:   `((\x0D\x0A)?(\x20|\x09))*`,
			invalidRegex: `[^\x20]&[^\x09]`, // difficult to specify
			rule:         lwsp(),
		},
		{
			name:         "OCTET",
			validRegex:   `[\x00-\xFF]`,
			invalidRegex: `[^\x00-\xFF]`,
			rule:         octet(),
		},
		{
			name:         "SP",
			validRegex:   `\x20`,
			invalidRegex: `[^\x20]`,
			rule:         sp(),
		},
		{
			name:         "VCHAR",
			validRegex:   `[\x21-\x7E]`,
			invalidRegex: `[^\x21-\x7E]`,
			rule:         vchar(),
		},
		{
			name:         "WSP",
			validRegex:   `\x20|\x09`,
			invalidRegex: `[^\x20]&[^\x09]`,
			rule:         wsp(),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			valid, _ := regen.New(test.validRegex)
			invalid, _ := regen.New(test.invalidRegex)

			for i := 0; i < 1000; i++ {
				validStr := valid.Generate()
				if ast := ParseString(validStr, test.rule); ast == nil {
					t.Errorf("no value found for: %s", validStr)
				} else {
					if !compareRunes(string(ast.Value), validStr) {
						t.Errorf("values do not match: %s %s", string(ast.Value), validStr)
					}
				}

				if invalidStr := invalid.Generate(); ParseString(invalidStr, alpha()) != nil {
					t.Errorf("value found for: %s", invalidStr)
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
			fmt.Println(a[i], b[i])
			return false
		}
	}
	return true
}
