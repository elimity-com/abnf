package abnf

import (
	"io/ioutil"
	"testing"
)

func TestNewRuleList(t *testing.T) {
	rawABNF, err := ioutil.ReadFile("./testdata/core.abnf")
	if err != nil {
		t.Error(err)
	}

	ruleSet := NewRuleSet(rawABNF)
	for _, rule := range []Rule{
		{
			name: "ALPHA",
			operator: AlternationOperator{
				key: "%x41-5A / %x61-7A",
				subOperators: []Operator{
					NumericValueOperator{
						key:         "%x41-5A",
						hyphen:      true,
						numericType: hexadecimal,
						value: []string{
							"41", "5A",
						},
					},
					NumericValueOperator{
						key:         "%x61-7A",
						hyphen:      true,
						numericType: hexadecimal,
						value: []string{
							"61", "7A",
						},
					},
				},
			},
		},
		{
			name: "BIT",
			operator: AlternationOperator{
				key: `"0" / "1"`,
				subOperators: []Operator{
					CharacterValueOperator{"0"},
					CharacterValueOperator{"1"},
				},
			},
		},
		{
			name: "CHAR",
			operator: NumericValueOperator{
				key:         "%x01-7F",
				hyphen:      true,
				numericType: hexadecimal,
				value:       []string{"01", "7F"},
			},
		},
		{
			name: "CR",
			operator: NumericValueOperator{
				key:         "%x0D",
				numericType: hexadecimal,
				value:       []string{"0D"},
			},
		},
		{
			name: "CRLF",
			operator: AlternationOperator{
				key: "CR LF / LF",
				subOperators: []Operator{
					ConcatenationOperator{
						key: "CR LF",
						subOperators: []Operator{
							RuleNameOperator{"CR"},
							RuleNameOperator{"LF"},
						},
					},
					RuleNameOperator{"LF"},
				},
			},
		},
		{
			name: "CTL",
			operator: AlternationOperator{
				key: "%x00-1F / %x7F",
				subOperators: []Operator{
					NumericValueOperator{
						key:         "%x00-1F",
						hyphen:      true,
						numericType: hexadecimal,
						value:       []string{"00", "1F"},
					},
					NumericValueOperator{
						key:         "%x7F",
						numericType: hexadecimal,
						value:       []string{"7F"},
					},
				},
			},
		},
		{
			name: "DIGIT",
			operator: NumericValueOperator{
				key:         "%x30-39",
				hyphen:      true,
				numericType: hexadecimal,
				value:       []string{"30", "39"},
			},
		},
		{
			name: "DQUOTE",
			operator: NumericValueOperator{
				key:         "%x22",
				numericType: hexadecimal,
				value:       []string{"22"},
			},
		},
		{
			name: "HEXDIG",
			operator: AlternationOperator{
				key: `DIGIT / "A" / "B" / "C" / "D" / "E" / "F" / "a" / "b" / "c" / "d" / "e" / "f"`,
				subOperators: []Operator{
					RuleNameOperator{"DIGIT"},
					CharacterValueOperator{"A"},
					CharacterValueOperator{"B"},
					CharacterValueOperator{"C"},
					CharacterValueOperator{"D"},
					CharacterValueOperator{"E"},
					CharacterValueOperator{"F"},
					CharacterValueOperator{"a"},
					CharacterValueOperator{"b"},
					CharacterValueOperator{"c"},
					CharacterValueOperator{"d"},
					CharacterValueOperator{"e"},
					CharacterValueOperator{"f"},
				},
			},
		},
		{
			name: "HTAB",
			operator: NumericValueOperator{
				key:         "%x09",
				numericType: hexadecimal,
				value:       []string{"09"},
			},
		},
		{
			name: "LF",
			operator: NumericValueOperator{
				key:         "%x0A",
				numericType: hexadecimal,
				value:       []string{"0A"},
			},
		},
		{
			name: "LWSP",
			operator: RepetitionOperator{
				key: "*(WSP / CRLF WSP)",
				min: 0,
				max: -1,
				subOperator: AlternationOperator{
					key: "WSP / CRLF WSP",
					subOperators: []Operator{
						RuleNameOperator{"WSP"},
						ConcatenationOperator{
							key: "CRLF WSP",
							subOperators: []Operator{
								RuleNameOperator{"CRLF"},
								RuleNameOperator{"WSP"},
							},
						},
					},
				},
			},
		},
		{
			name: "OCTET",
			operator: NumericValueOperator{
				key:         "%x00-FF",
				hyphen:      true,
				numericType: hexadecimal,
				value:       []string{"00", "FF"},
			},
		},
		{
			name: "SP",
			operator: NumericValueOperator{
				key:         "%x20",
				numericType: hexadecimal,
				value:       []string{"20"},
			},
		},
		{
			name: "VCHAR",
			operator: NumericValueOperator{
				key:         "%x21-7E",
				hyphen:      true,
				numericType: hexadecimal,
				value:       []string{"21", "7E"},
			},
		},
		{
			name: "WSP",
			operator: AlternationOperator{
				key: "SP / HTAB",
				subOperators: []Operator{
					RuleNameOperator{"SP"},
					RuleNameOperator{"HTAB"},
				},
			},
		},
	} {
		if err := rule.Equals(ruleSet[rule.name]); err != nil {
			t.Errorf("%s: %s", rule.name, err)
		}
	}
}

func TestNumericValue(t *testing.T) {
	abnf := "false = %x66.61.6c.73.65 ; false\nnull = %x6e.75.6c.6c ; null\ntrue = %x74.72.75.65 ; true\n"
	set := NewRuleSet([]byte(abnf))
	for _, rule := range []Rule{
		{
			name: "false",
			operator: NumericValueOperator{
				key:         "%x66.61.6c.73.65",
				points:      true,
				numericType: hexadecimal,
				value:       []string{"66", "61", "6c", "73", "65"},
			},
		},
		{
			name: "null",
			operator: NumericValueOperator{
				key:         "%x6e.75.6c.6c",
				points:      true,
				numericType: hexadecimal,
				value:       []string{"6e", "75", "6c", "6c"},
			},
		},
		{
			name: "true",
			operator: NumericValueOperator{
				key:         "%x74.72.75.65",
				points:      true,
				numericType: hexadecimal,
				value:       []string{"74", "72", "75", "65"},
			},
		},
	} {
		if err := rule.Equals(set[rule.name]); err != nil {
			t.Errorf("%s: %s", rule.name, err)
		}
	}
}

func TestRepetition(t *testing.T) {
	abnf := "rep04 = *4X\nrep14 = 1*4X\nrep44 = 4X\nrep0i = *X\nrep4i = 4*X\n"
	set := NewRuleSet([]byte(abnf))
	for _, rule := range []Rule{
		{
			name: "rep04",
			operator: RepetitionOperator{
				key: "*4X",
				min: 0, max: 4,
				subOperator: RuleNameOperator{"X"},
			},
		},
		{
			name: "rep14",
			operator: RepetitionOperator{
				key: "1*4X",
				min: 1, max: 4,
				subOperator: RuleNameOperator{"X"},
			},
		},
		{
			name: "rep44",
			operator: RepetitionOperator{
				key: "4X",
				min: 4, max: 4,
				subOperator: RuleNameOperator{"X"},
			},
		},
		{
			name: "rep0i",
			operator: RepetitionOperator{
				key: "*X",
				min: 0, max: -1,
				subOperator: RuleNameOperator{"X"},
			},
		},
		{
			name: "rep4i",
			operator: RepetitionOperator{
				key: "4*X",
				min: 4, max: -1,
				subOperator: RuleNameOperator{"X"},
			},
		},
	} {
		if err := rule.Equals(set[rule.name]); err != nil {
			t.Errorf("%s: %s", rule.name, err)
		}
	}
}
