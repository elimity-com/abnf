package core

import (
	"fmt"
	"github.com/elimity-com/abnf/encoding"
	"github.com/elimity-com/abnf/operators"
	"testing"

	"github.com/di-wu/regen"
)

func TestCore(t *testing.T) {
	for _, test := range []struct {
		name                     string
		validRegex, invalidRegex string
		rule                     operators.Operator
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

				// these rules are only valid for ABNF encoded in 7-bit ASCII
				e := encoding.ASCII.NewEncoder()
				validStr, _ = e.String(validStr)

				if nodes := test.rule([]byte(validStr)); nodes == nil {
					t.Errorf("no value found for: %s", validStr)
				} else {
					if best := nodes.Best(); !compareRunes(string(best.Value), validStr) {
						t.Errorf("values do not match: %s %s", string(best.Value), validStr)
					}
				}

				invalidStr := invalid.Generate()
				if nodes := test.rule([]byte(invalidStr)); len(nodes) != 0 {
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
		rule    operators.Operator
		str     string
		correct operators.Alternatives
	}{
		{
			name: "Alpha Lower",
			rule: ALPHA(),
			str:  "a",
			correct: operators.Alternatives{
				{
					Key:   "ALPHA",
					Value: []byte("a"),
					Children: operators.Children{
						{
							Key:   "%x61-7A",
							Value: []byte("a"),
						},
					},
				},
			},
		},
		{
			name: "Alpha Upper",
			rule: ALPHA(),
			str:  "Z",
			correct: operators.Alternatives{
				{
					Key:   "ALPHA",
					Value: []byte("Z"),
					Children: operators.Children{
						{
							Key:   "%x41-5A",
							Value: []byte("Z"),
						},
					},
				},
			},
		},
		{
			name: "Bit 0",
			rule: BIT(),
			str:  "0",
			correct: operators.Alternatives{
				{
					Key:   "BIT",
					Value: []byte("0"),
					Children: operators.Children{
						{
							Key:   "\"0\"",
							Value: []byte("0"),
						},
					},
				},
			},
		},
		{
			name: "Bit 1",
			rule: BIT(),
			str:  "1",
			correct: operators.Alternatives{
				{
					Key:   "BIT",
					Value: []byte("1"),
					Children: operators.Children{
						{
							Key:   "\"1\"",
							Value: []byte("1"),
						},
					},
				},
			},
		},
		{
			name: "Character",
			rule: CHAR(),
			str:  "~",
			correct: operators.Alternatives{
				{
					Key:   "CHAR",
					Value: []byte("~"),
				},
			},
		},
		{
			name: "NewLine",
			rule: CRLF(),
			str:  "\r\n",
			correct: operators.Alternatives{
				{
					Key:   "CRLF",
					Value: []byte("\r\n"),
					Children: operators.Children{
						{
							Key:   "CR LF",
							Value: []byte("\r\n"),
							Children: operators.Children{
								{
									Key:   "CR",
									Value: []byte("\r"),
								},
								{
									Key:   "LF",
									Value: []byte("\n"),
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
			correct: operators.Alternatives{
				{
					Key:   "CRLF",
					Value: []byte("\n"),
					Children: operators.Children{
						{
							Key:   "LF",
							Value: []byte("\n"),
						},
					},
				},
			},
		},
		{
			name: "Control",
			rule: CTL(),
			str:  "\u001B", // escape
			correct: operators.Alternatives{
				{
					Key:   "CTL",
					Value: []byte("\u001B"),
					Children: operators.Children{
						{
							Key:   "%x00-1F",
							Value: []byte("\u001B"),
						},
					},
				},
			},
		},
		{
			name: "Digit",
			rule: DIGIT(),
			str:  "7",
			correct: operators.Alternatives{
				{
					Key:   "DIGIT",
					Value: []byte("7"),
				},
			},
		},
		{
			name: "DoubleQuote",
			rule: DQUOTE(),
			str:  "\"",
			correct: operators.Alternatives{
				{
					Key:   "DQUOTE",
					Value: []byte("\""),
				},
			},
		},
		{
			name: "HexDigit Digit",
			rule: HEXDIG(),
			str:  "7",
			correct: operators.Alternatives{
				{
					Key:   "HEXDIG",
					Value: []byte("7"),
					Children: operators.Children{
						{
							Key:   "DIGIT",
							Value: []byte("7"),
						},
					},
				},
			},
		},
		{
			name: "HexDigit Hex",
			rule: HEXDIG(),
			str:  "A",
			correct: operators.Alternatives{
				{
					Key:   "HEXDIG",
					Value: []byte("A"),
					Children: operators.Children{
						{
							Key:   "\"A\"",
							Value: []byte("A"),
						},
					},
				},
			},
		},
		{
			name: "HorizontalTab",
			rule: HTAB(),
			str:  "\t",
			correct: operators.Alternatives{
				{
					Key:   "HTAB",
					Value: []byte("\t"),
				},
			},
		},
		{
			name: "Linefeed",
			rule: LF(),
			str:  "\n",
			correct: operators.Alternatives{
				{
					Key:   "LF",
					Value: []byte("\n"),
				},
			},
		},
		{
			name: "LinearWhiteSpace Space",
			rule: LWSP(),
			str:  " ",
			correct: operators.Alternatives{
				{
					Key:   "LWSP",
					Value: []byte(" "),
					Children: operators.Children{
						{
							Key:   "WSP / CRLF WSP",
							Value: []byte(" "),
							Children: operators.Children{
								{
									Key:   "WSP",
									Value: []byte(" "),
									Children: operators.Children{
										{
											Key:   "SP",
											Value: []byte(" "),
										},
									},
								},
							},
						},
					},
				},
				{
					Key:   "LWSP",
					Value: []byte(""),
				},
			},
		},
		{
			name: "LinearWhiteSpace EmptyLine",
			rule: LWSP(),
			str:  "\n ",
			correct: operators.Alternatives{
				{
					Key:   "LWSP",
					Value: []byte("\n "),
					Children: operators.Children{
						{
							Key:   "WSP / CRLF WSP",
							Value: []byte("\n "),
							Children: operators.Children{
								{
									Key:   "CRLF WSP",
									Value: []byte("\n "),
									Children: operators.Children{
										{
											Key:   "CRLF",
											Value: []byte("\n"),
											Children: operators.Children{
												{
													Key:   "LF",
													Value: []byte("\n"),
												},
											},
										},
										{
											Key:   "WSP",
											Value: []byte(" "),
											Children: operators.Children{
												{
													Key:   "SP",
													Value: []byte(" "),
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
					Value: []byte(""),
				},
			},
		},
		{
			name: "Octet",
			rule: OCTET(),
			str:  "o",
			correct: operators.Alternatives{
				{
					Key:   "OCTET",
					Value: []byte("o"),
				},
			},
		},
		{
			name: "Space",
			rule: SP(),
			str:  " ",
			correct: operators.Alternatives{
				{
					Key:   "SP",
					Value: []byte(" "),
				},
			},
		},
		{
			name: "VisibleCharacters",
			rule: VCHAR(),
			str:  "~",
			correct: operators.Alternatives{
				{
					Key:   "VCHAR",
					Value: []byte("~"),
				},
			},
		},
		{
			name: "WhiteSpace Space",
			rule: WSP(),
			str:  " ",
			correct: operators.Alternatives{
				{
					Key:   "WSP",
					Value: []byte(" "),
					Children: operators.Children{
						{
							Key:   "SP",
							Value: []byte(" "),
						},
					},
				},
			},
		},
		{
			name: "WhiteSpace Tab",
			rule: WSP(),
			str:  "\t",
			correct: operators.Alternatives{
				{
					Key:   "WSP",
					Value: []byte("\t"),
					Children: operators.Children{
						{
							Key:   "HTAB",
							Value: []byte("\t"),
						},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			nodes := test.rule([]byte(test.str))
			if err := nodes.Equals(test.correct); err != nil {
				for _, node := range nodes {
					fmt.Print(node.StringRecursive())
				}
				t.Error(err)
			}
		})
	}
}
