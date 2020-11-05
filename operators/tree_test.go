package operators

import (
	"testing"
)

func TestNode(t *testing.T) {
	for _, test := range []struct {
		name    string
		rule    Operator
		str     string
		correct Alternatives
	}{
		{
			name: "Terminal",
			rule: a,
			str:  "a",
			correct: Alternatives{
				{
					Key:   "a",
					Value: []byte("a"),
				},
			},
		},
		{
			name: "RuneEmpty",
			rule: a,
			str:  "b",
		},
		{
			name: "String",
			rule: String(`cs`, "str"),
			str:  "str",
			correct: Alternatives{
				{
					Key:   "cs",
					Value: []byte("str"),
				},
			},
		},
		{
			name: "StringEmpty",
			rule: String(`cs`, "str"),
			str:  "rts",
		},
		{
			name: "String",
			rule: StringCI(`s`, "str"),
			str:  "Str",
			correct: Alternatives{
				{
					Key:   "s",
					Value: []byte("Str"),
				},
			},
		},
		{
			name: "StringCIEmpty",
			rule: String(`s`, "str"),
			str:  "rts",
		},
		{
			name: "Range",
			rule: Range(`r`, []byte("a"), []byte("z")),
			str:  "x",
			correct: Alternatives{
				{
					Key:   "r",
					Value: []byte("x"),
				},
			},
		},
		{
			name: "RangeEmpty",
			rule: Range(`r`, []byte("a"), []byte("z")),
			str:  "0",
		},
		{
			name: "Optional",
			rule: Optional(`o`, a),
			str:  "a",
			correct: Alternatives{
				{
					Key:   "o",
					Value: []byte("a"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						},
					},
				},
				{
					Key:   "o",
					Value: []byte(""),
				},
			},
		},
		{
			name: "OptionalEmpty",
			rule: Optional(`o`, a),
			correct: Alternatives{
				{
					Key:   "o",
					Value: []byte(""),
				},
			},
		},
		{
			name: "Repeat0",
			rule: RepeatN(`r0`, 0, a),
			correct: Alternatives{
				{
					Key:   "r0",
					Value: []byte(""),
				},
			},
		},
		{
			name: "Repeat1",
			rule: RepeatN(`r1`, 1, a),
			str:  "aaa",
			correct: Alternatives{
				{
					Key:   "r1",
					Value: []byte("a"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						},
					},
				},
			},
		},
		{
			name: "Repeat1Empty",
			rule: RepeatN(`r1`, 1, a),
			str:  "bbb",
		},
		{
			name: "Repeat0Inf",
			rule: Repeat0Inf(`r0i`, a),
			str:  "aa",
			correct: Alternatives{
				{
					Key:   "r0i",
					Value: []byte("aa"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						}, {
							Key:   "a",
							Value: []byte("a"),
						},
					},
				},
				{
					Key:   "r0i",
					Value: []byte("a"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						},
					},
				},
				{
					Key:   "r0i",
					Value: []byte(""),
				},
			},
		},
		{
			name: "Repeat0InfEmpty",
			rule: Repeat0Inf(`r0i`, a),
			str:  "",
			correct: Alternatives{
				{
					Key:   "r0i",
					Value: []byte(""),
				},
			},
		},
		{
			name: "Repeat1Inf",
			rule: Repeat1Inf(`r1i`, a),
			str:  "aa",
			correct: Alternatives{
				{
					Key:   "r1i",
					Value: []byte("aa"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						},
						{
							Key:   "a",
							Value: []byte("a"),
						},
					},
				},
				{
					Key:   "r1i",
					Value: []byte("a"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						},
					},
				},
			},
		},
		{
			name: "Repeat1InfEmpty",
			rule: Repeat1Inf(`r1i`, a),
			str:  "bbb",
		},
		{
			name: "RepeatOptional",
			rule: RepeatOptional(`ro`, a),
			str:  "a",
			correct: Alternatives{
				{
					Key:   "ro",
					Value: []byte("a"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						},
					},
				},
				{
					Key:   "ro",
					Value: []byte(""),
				},
			},
		},
		{
			name: "RepeatOptionalEmpty",
			rule: RepeatOptional(`ro`, a),
			str:  "bbb",
			correct: Alternatives{
				{
					Key:   "ro",
					Value: []byte(""),
				},
			},
		},
		{
			name: "Concat",
			rule: Concat(`c`, a, b, c),
			str:  "abc",
			correct: Alternatives{
				{
					Key:   "c",
					Value: []byte("abc"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						},
						{
							Key:   "b",
							Value: []byte("b"),
						},
						{
							Key:   "c",
							Value: []byte("c"),
						},
					},
				},
			},
		},
		{
			name: "ConcatEmpty",
			rule: Concat(`c`, a, b, c),
			str:  "cba",
		},
		{
			name: "ConcatRepeat",
			rule: Concat(`cr`,
				Repeat0Inf(`r0i1`, a),
				Repeat0Inf(`r0i2`, a),
			),
			str: "a",
			correct: Alternatives{
				{
					Key:   "cr",
					Value: []byte("a"),
					Children: Children{
						{
							Key:   "r0i1",
							Value: []byte("a"),
							Children: Children{
								{
									Key:   "a",
									Value: []byte("a"),
								},
							},
						},
						{
							Key:   "r0i2",
							Value: []byte(""),
						},
					},
				},
				{
					Key:   "cr",
					Value: []byte("a"),
					Children: Children{
						{
							Key:   "r0i1",
							Value: []byte(""),
						},
						{
							Key:   "r0i2",
							Value: []byte("a"),
							Children: Children{
								{
									Key:   "a",
									Value: []byte("a"),
								},
							},
						},
					},
				},
				{
					Key:   "cr",
					Value: []byte(""),
					Children: Children{
						{
							Key:   "r0i1",
							Value: []byte(""),
						},
						{
							Key:   "r0i2",
							Value: []byte(""),
						},
					},
				},
			},
		},
		{
			name: "AltsOption1",
			rule: Alts(`a`, a, b),
			str:  "a",
			correct: Alternatives{
				{
					Key:   "a",
					Value: []byte("a"),
					Children: Children{
						{
							Key:   "a",
							Value: []byte("a"),
						},
					},
				},
			},
		},
		{
			name: "AltsOption2",
			rule: Alts(`a`, a, b),
			str:  "b",
			correct: Alternatives{
				{
					Key:   "a",
					Value: []byte("b"),
					Children: Children{
						{
							Key:   "b",
							Value: []byte("b"),
						},
					},
				},
			},
		},
		{
			name: "AltsEmpty",
			rule: Alts(`a`, a, b),
			str:  "c",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			nodes := test.rule([]byte(test.str))
			if err := nodes.Equals(test.correct); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestBestAlternative(t *testing.T) {
	t.Run(`Simple`, func(t *testing.T) {
		str := "aaa"
		nodes := Repeat0Inf(``, a)([]byte(str))
		if best := nodes.Best(); best.String() != str {
			t.Error("did not get best alternative")
		}
	})
}
