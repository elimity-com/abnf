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
				key: `DIGIT / "A" / "B" / "C" / "D" / "E" / "F"`,
				subOperators: []Operator{
					RuleNameOperator{"DIGIT"},
					CharacterValueOperator{"A"},
					CharacterValueOperator{"B"},
					CharacterValueOperator{"C"},
					CharacterValueOperator{"D"},
					CharacterValueOperator{"E"},
					CharacterValueOperator{"F"},
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
				key:          "SP / HTAB",
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
