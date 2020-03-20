package abnf

import (
	"testing"

	"github.com/di-wu/regen"
)

func TestDefinition(t *testing.T) {
	for _, test := range []struct {
		name                     string
		validRegex, invalidRegex string
		rule                     RuleFunc
	}{
		{
			name:         "CharVal",
			validRegex:   `"[a-zA-Z]"`,
			invalidRegex: `[a-zA-Z]`,
			rule:         charVal(),
		},
		{
			name:       "NumVal",
			validRegex: `%((b[0-1]+(.[0-1]+|-[0-1]+)?)|(d\d+(.\d+|-\d+)?)|(x[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?))`,
			rule:       numVal(),
		},
		{
			name:         "BinVal",
			validRegex:   `b[0-1]+(.[0-1]+|-[0-1]+)?`,
			invalidRegex: `[0-1]+(.[0-1]+|-[0-1]+)?`,
			rule:         binVal(),
		},
		{
			name:         "DecVal",
			validRegex:   `d\d+(.\d+|-\d+)?`,
			invalidRegex: `\d+(.\d+|-\d+)?`,
			rule:         decVal(),
		},
		{
			name:         "HexVal",
			validRegex:   `x[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?`,
			invalidRegex: `[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?`,
			rule:         hexVal(),
		},
		{
			name:         "ProseVal",
			validRegex:   `<[a-zA-Z]*>`,
			invalidRegex: `[a-zA-Z]*`,
			rule:         proseVal(),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			valid, _ := regen.New(test.validRegex)
			invalid, _ := regen.New(test.invalidRegex)

			for i := 0; i < 1000; i++ {
				validStr := valid.Generate()
				if value := ParseString(validStr, test.rule); value == nil {
					t.Errorf("no value found for: %s", validStr)
				} else {
					if !compareRunes(string(value), validStr) {
						t.Errorf("values do not match: %s %s", string(value), validStr)
					}
				}

				if invalidStr := invalid.Generate(); ParseString(invalidStr, test.rule) != nil {
					t.Errorf("tree found for: %s", invalidStr)
				}
			}
		})
	}
}
