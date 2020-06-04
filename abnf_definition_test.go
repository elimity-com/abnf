package abnf

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/di-wu/abnf/operators"
	"github.com/di-wu/regen"
)

func TestDefinition(t *testing.T) {
	for _, test := range []struct {
		name     string
		rule     Operator
		examples []string
	}{
		{
			name: "ruleName",
			rule: ruleName,
			examples: []string{
				"name",
				`name123`,
				`name1-2-3`,
			},
		},
		{
			name: "definedAs",
			rule: definedAs,
			examples: []string{
				" = ",
				"=/",
			},
		},
		{
			name: `element`,
			rule: element,
			examples: []string{
				`rule-name`, // rule name
				`( %x01 )`,  // group
				`[ %x01 ]`,  // option
				`"charval"`, // CHAR value
				`%x01`,      // numerical value
				`<abc>`,     // prose value
			},
		},
	} {
		for _, s := range test.examples {
			t.Run(fmt.Sprintf("%s %s", test.name, s), func(t *testing.T) {
				if value := test.rule([]rune(s)); value == nil {
					t.Errorf("no value found for: %s", s)
				}
			})
		}
	}
}

func TestValues(t *testing.T) {
	for _, test := range []struct {
		name                     string
		validRegex, invalidRegex string
		rule                     Operator
	}{
		{
			name:         "CharVal",
			validRegex:   `"[a-zA-Z]"`,
			invalidRegex: `[a-zA-Z]`,
			rule:         charVal,
		},
		{
			name:       "NumVal",
			validRegex: `%((b[0-1]+(.[0-1]+|-[0-1]+)?)|(d\d+(.\d+|-\d+)?)|(x[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?))`,
			rule:       numVal,
		},
		{
			name:         "BinVal",
			validRegex:   `b[0-1]+(.[0-1]+|-[0-1]+)?`,
			invalidRegex: `[0-1]+(.[0-1]+|-[0-1]+)?`,
			rule:         binVal,
		},
		{
			name:         "DecVal",
			validRegex:   `d\d+(.\d+|-\d+)?`,
			invalidRegex: `\d+(.\d+|-\d+)?`,
			rule:         decVal,
		},
		{
			name:         "HexVal",
			validRegex:   `x[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?`,
			invalidRegex: `[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?`,
			rule:         hexVal,
		},
		{
			name:         "ProseVal",
			validRegex:   `<[a-zA-Z]*>`,
			invalidRegex: `[a-zA-Z]*`,
			rule:         proseVal,
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

				if invalidStr := invalid.Generate(); test.rule([]rune(invalidStr)) != nil {
					t.Errorf("tree found for: %s", invalidStr)
				}
			}
		})
	}
}

func TestABNF(t *testing.T) {
	raw, err := ioutil.ReadFile("testdata/core.abnf")
	if err != nil {
		t.Error(err)
	}
	strABNF := string(raw)
	list := ruleList([]rune(strABNF)).Best()

	if list.String() != strABNF {
		t.Error("parsed abnf does not match original")
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
