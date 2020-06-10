package core

import (
	"fmt"
	"testing"

	"github.com/di-wu/regen"
	. "github.com/elimity-com/abnf/operators"
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

func TestNode(t *testing.T) {
	for _, test := range []struct {
		name    string
		rule    Operator
		str     string
		correct Alternatives
	}{
		{
			name: "Alpha Lower",
			rule: ALPHA(),
			str:  "a",
			correct: Alternatives{
				{
					Key:   "ALPHA",
					Value: []rune("a"),
					Children: Children{
						{
							Key:   "%x61-7A",
							Value: []rune("a"),
						},
					},
				},
			},
		},
		{
			name: "Alpha Upper",
			rule: ALPHA(),
			str:  "Z",
			correct: Alternatives{
				{
					Key:   "ALPHA",
					Value: []rune("Z"),
					Children: Children{
						{
							Key:   "%x41-5A",
							Value: []rune("Z"),
						},
					},
				},
			},
		},
		{
			name: "Bit 0",
			rule: BIT(),
			str:  "0",
			correct: Alternatives{
				{
					Key:   "BIT",
					Value: []rune("0"),
					Children: Children{
						{
							Key:   "\"0\"",
							Value: []rune("0"),
						},
					},
				},
			},
		},
		{
			name: "Bit 1",
			rule: BIT(),
			str:  "1",
			correct: Alternatives{
				{
					Key:   "BIT",
					Value: []rune("1"),
					Children: Children{
						{
							Key:   "\"1\"",
							Value: []rune("1"),
						},
					},
				},
			},
		},
		{
			name: "Character",
			rule: CHAR(),
			str:  "~",
			correct: Alternatives{
				{
					Key:   "CHAR",
					Value: []rune("~"),
				},
			},
		},
		{
			name: "NewLine",
			rule: CRLF(),
			str:  "\r\n",
			correct: Alternatives{
				{
					Key:   "CRLF",
					Value: []rune("\r\n"),
					Children: Children{
						{
							Key:   "CR LF",
							Value: []rune("\r\n"),
							Children: Children{
								{
									Key:   "CR",
									Value: []rune("\r"),
								},
								{
									Key:   "LF",
									Value: []rune("\n"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NewLine Unix",
			rule: CRLF(),
			str:  "\n",
			correct: Alternatives{
				{
					Key:   "CRLF",
					Value: []rune("\n"),
					Children: Children{
						{
							Key:   "LF",
							Value: []rune("\n"),
						},
					},
				},
			},
		},
		{
			name: "Control",
			rule: CTL(),
			str:  "\u001B", // escape
			correct: Alternatives{
				{
					Key:   "CTL",
					Value: []rune("\u001B"),
					Children: Children{
						{
							Key:   "%x00-1F",
							Value: []rune("\u001B"),
						},
					},
				},
			},
		},
		{
			name: "Digit",
			rule: DIGIT(),
			str:  "7",
			correct: Alternatives{
				{
					Key:   "DIGIT",
					Value: []rune("7"),
				},
			},
		},
		{
			name: "DoubleQuote",
			rule: DQUOTE(),
			str:  "\"",
			correct: Alternatives{
				{
					Key:   "DQUOTE",
					Value: []rune("\""),
				},
			},
		},
		{
			name: "HexDigit Digit",
			rule: HEXDIG(),
			str:  "7",
			correct: Alternatives{
				{
					Key:   "HEXDIG",
					Value: []rune("7"),
					Children: Children{
						{
							Key:   "DIGIT",
							Value: []rune("7"),
						},
					},
				},
			},
		},
		{
			name: "HexDigit Hex",
			rule: HEXDIG(),
			str:  "A",
			correct: Alternatives{
				{
					Key:   "HEXDIG",
					Value: []rune("A"),
					Children: Children{
						{
							Key:   "\"A\"",
							Value: []rune("A"),
						},
					},
				},
			},
		},
		{
			name: "HorizontalTab",
			rule: HTAB(),
			str:  "\t",
			correct: Alternatives{
				{
					Key:   "HTAB",
					Value: []rune("\t"),
				},
			},
		},
		{
			name: "Linefeed",
			rule: LF(),
			str:  "\n",
			correct: Alternatives{
				{
					Key:   "LF",
					Value: []rune("\n"),
				},
			},
		},
		{
			name: "LinearWhiteSpace Space",
			rule: LWSP(),
			str:  " ",
			correct: Alternatives{
				{
					Key:   "LWSP",
					Value: []rune(" "),
					Children: Children{
						{
							Key:   "WSP / CRLF WSP",
							Value: []rune(" "),
							Children: Children{
								{
									Key:   "WSP",
									Value: []rune(" "),
									Children: Children{
										{
											Key:   "SP",
											Value: []rune(" "),
										},
									},
								},
							},
						},
					},
				},
				{
					Key:   "LWSP",
					Value: []rune(""),
				},
			},
		},
		{
			name: "LinearWhiteSpace EmptyLine",
			rule: LWSP(),
			str:  "\n ",
			correct: Alternatives{
				{
					Key:   "LWSP",
					Value: []rune("\n "),
					Children: Children{
						{
							Key:   "WSP / CRLF WSP",
							Value: []rune("\n "),
							Children: Children{
								{
									Key:   "CRLF WSP",
									Value: []rune("\n "),
									Children: Children{
										{
											Key:   "CRLF",
											Value: []rune("\n"),
											Children: Children{
												{
													Key:   "LF",
													Value: []rune("\n"),
												},
											},
										},
										{
											Key:   "WSP",
											Value: []rune(" "),
											Children: Children{
												{
													Key:   "SP",
													Value: []rune(" "),
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Key:   "LWSP",
					Value: []rune(""),
				},
			},
		},
		{
			name: "Octet",
			rule: OCTET(),
			str:  "o",
			correct: Alternatives{
				{
					Key:   "OCTET",
					Value: []rune("o"),
				},
			},
		},
		{
			name: "Space",
			rule: SP(),
			str:  " ",
			correct: Alternatives{
				{
					Key:   "SP",
					Value: []rune(" "),
				},
			},
		},
		{
			name: "VisibleCharacters",
			rule: VCHAR(),
			str:  "~",
			correct: Alternatives{
				{
					Key:   "VCHAR",
					Value: []rune("~"),
				},
			},
		},
		{
			name: "WhiteSpace Space",
			rule: WSP(),
			str:  " ",
			correct: Alternatives{
				{
					Key:   "WSP",
					Value: []rune(" "),
					Children: Children{
						{
							Key:   "SP",
							Value: []rune(" "),
						},
					},
				},
			},
		},
		{
			name: "WhiteSpace Tab",
			rule: WSP(),
			str:  "\t",
			correct: Alternatives{
				{
					Key:   "WSP",
					Value: []rune("\t"),
					Children: Children{
						{
							Key:   "HTAB",
							Value: []rune("\t"),
						},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			nodes := test.rule([]rune(test.str))
			if err := nodes.Equals(test.correct); err != nil {
				for _, node := range nodes {
					fmt.Print(node.StringRecursive())
				}
				t.Error(err)
			}
		})
	}
}
