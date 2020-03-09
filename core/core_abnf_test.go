package core

import (
	"fmt"
	"github.com/di-wu/regen"
	"testing"

	. "github.com/di-wu/abnf"
)

func TestCore(t *testing.T) {
	for _, test := range []struct {
		name                     string
		validRegex, invalidRegex string
		rule                     Rule
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
			invalidRegex: `[^\x20]&[^\x09]`, // difficult to specify
			rule:         LWSP(),
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
			name:         "WSP",
			validRegex:   `\x20|\x09`,
			invalidRegex: `[^\x20]&[^\x09]`,
			rule:         WSP(),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			valid, _ := regen.New(test.validRegex)
			invalid, _ := regen.New(test.invalidRegex)

			for i := 0; i < 1000; i++ {
				validStr := valid.Generate()
				if value := ParseString(validStr, test.rule); value == nil {
					t.Errorf("no tree found for: %s", validStr)
				} else {
					if !compareRunes(string(value), validStr) {
						t.Errorf("values do not match: %s %s", string(value), validStr)
					}
				}

				if invalidStr := invalid.Generate(); ParseString(invalidStr, ALPHA()) != nil {
					t.Errorf("tree fround for: %s", invalidStr)
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
