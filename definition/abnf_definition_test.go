package definition

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/di-wu/regen"
	"github.com/elimity-com/abnf/operators"
)

func TestDefinition(t *testing.T) {
	for _, test := range []struct {
		name     string
		rule     operators.Operator
		examples []string
	}{
		{
			name: "ruleName",
			rule: Rulename,
			examples: []string{
				"name",
				`name123`,
				`name1-2-3`,
			},
		},
		{
			name: "definedAs",
			rule: DefinedAs,
			examples: []string{
				" = ",
				"=/",
			},
		},
		{
			name: `element`,
			rule: Element,
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
				if value := test.rule([]byte(s)); value == nil {
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
		rule                     operators.Operator
	}{
		{
			name:         "CharVal",
			validRegex:   `"[a-zA-Z]"`,
			invalidRegex: `[a-zA-Z]`,
			rule:         CharVal,
		},
		{
			name:       "NumVal",
			validRegex: `%((b[0-1]+(.[0-1]+|-[0-1]+)?)|(d\d+(.\d+|-\d+)?)|(x[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?))`,
			rule:       NumVal,
		},
		{
			name:         "BinVal",
			validRegex:   `b[0-1]+(.[0-1]+|-[0-1]+)?`,
			invalidRegex: `[0-1]+(.[0-1]+|-[0-1]+)?`,
			rule:         BinVal,
		},
		{
			name:         "DecVal",
			validRegex:   `d\d+(.\d+|-\d+)?`,
			invalidRegex: `\d+(.\d+|-\d+)?`,
			rule:         DecVal,
		},
		{
			name:         "HexVal",
			validRegex:   `x[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?`,
			invalidRegex: `[0-9A-F]+(.[0-9A-F]+|-[0-9A-F]+)?`,
			rule:         HexVal,
		},
		{
			name:         "ProseVal",
			validRegex:   `<[a-zA-Z]*>`,
			invalidRegex: `[a-zA-Z]*`,
			rule:         ProseVal,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			valid, _ := regen.New(test.validRegex)
			invalid, _ := regen.New(test.invalidRegex)

			for i := 0; i < 1000; i++ {
				validStr := valid.Generate()
				if nodes := test.rule([]byte(validStr)); nodes == nil {
					t.Errorf("no value found for: %s", validStr)
				} else {
					if best := nodes.Best(); !compareRunes(string(best.Value), validStr) {
						t.Errorf("values do not match: %s %s", string(best.Value), validStr)
					}
				}

				if invalidStr := invalid.Generate(); test.rule([]byte(invalidStr)) != nil {
					t.Errorf("tree found for: %s", invalidStr)
				}
			}
		})
	}
}

func TestABNF(t *testing.T) {
	raw, err := ioutil.ReadFile("../testdata/core.abnf")
	if err != nil {
		t.Error(err)
	}
	strABNF := string(raw)
	list := Rulelist([]byte(strABNF)).Best()

	if list.String() != regexp.MustCompile(`\s+`).ReplaceAllString(strABNF, " ") {
		t.Error("parsed abnf does not match original")
	}

	if l := len(list.GetSubNodes("rule")); l != 16 {
		t.Errorf("should have 16 rules, got %d", l)
	}

	if l := len(list.GetSubNodes("=")); l != 16 {
		t.Errorf("should have 16 =, got %d", l)
	}

	if l := len(list.GetSubNodes("comment")); l != 22 {
		t.Errorf("should have 22 comments, got %d", l)
	}

	if l := len(list.GetSubNodes("CRLF")); l != 35 {
		t.Errorf("should have 35 EOLs, got %d", l)
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
