package abnf

import (
	"github.com/elimity-com/abnf/core"
	. "github.com/elimity-com/abnf/operators"
)

// RFC 5234: 4. ABNF Definition of ABNF

func ruleList(s []rune) Alternatives {
	return Repeat1Inf(`rulelist`, Alts(
		`rule / (*c-wsp c-nl)`,
		rule,
		Concat(
			`*WSP c-nl`,
			repeatWSP,
			cNl,
		),
	))(s)
}

func rule(s []rune) Alternatives {
	return Concat(
		`rule`,
		ruleName,
		definedAs,
		elements,
		cNl,
	)(s)
}

func ruleName(s []rune) Alternatives {
	return Concat(
		`rulename`,
		core.ALPHA(),
		Repeat0Inf(`*(ALPHA / DIGIT / "-")`, Alts(
			`ALPHA / DIGIT / "-"`,
			core.ALPHA(),
			core.DIGIT(),
			Rune(`-`, '-'),
		)),
	)(s)
}

func definedAs(s []rune) Alternatives {
	return Concat(
		`defined-as`,
		repeatCWsp,
		Alts(
			`"=" / "=/"`,
			Rune(`=`, '='),
			Concat(`=/`, Rune(`=`, '='), Rune(`/`, '/')),
		),
		repeatCWsp,
	)(s)
}

func elements(s []rune) Alternatives {
	return Concat(
		`elements`,
		alternation,
		repeatWSP,
	)(s)
}

func cWsp(s []rune) Alternatives {
	return Alts(
		`c-wsp`,
		core.WSP(),
		Concat(
			`c-nl WSP`,
			cNl,
			core.WSP(),
		),
	)(s)
}

func repeatWSP(s []rune) Alternatives {
	return Repeat0Inf(`*WSP`, core.WSP())(s)
}

func repeatCWsp(s []rune) Alternatives {
	return Repeat0Inf(`*c-wsp`, cWsp)(s)
}

func cNl(s []rune) Alternatives {
	return Alts(
		`c-nl`,
		comment,
		core.CRLF(),
	)(s)
}

func comment(s []rune) Alternatives {
	return Concat(
		`comment`,
		Rune(`;`, ';'),
		Repeat0Inf(`*(WSP / VCHAR) CRLF`, Alts(
			`WSP / VCHAR`,
			core.WSP(),
			core.VCHAR(),
		)),
		core.CRLF(),
	)(s)
}

func alternation(s []rune) Alternatives {
	return Concat(
		`alternation`,
		concatenation,
		Repeat0Inf(`*(*c-wsp "/" *c-wsp concatenation)`, Concat(
			`*c-wsp "/" *c-wsp concatenation`,
			repeatCWsp,
			Rune(`/`, '/'),
			repeatCWsp,
			concatenation,
		)),
	)(s)
}

func concatenation(s []rune) Alternatives {
	return Concat(
		`concatenation`,
		repetition,
		Repeat0Inf(`*(1*c-wsp repetition)`, Concat(
			`1*c-wsp repetition`,
			Repeat1Inf(`1*c-wsp`, cWsp),
			repetition,
		)),
	)(s)
}

func repetition(s []rune) Alternatives {
	return Concat(
		`repetition`,
		Optional(`[repeat]`, repeat),
		element,
	)(s)
}

func repeat(s []rune) Alternatives {
	return Alts(
		`repeat`,
		Repeat1Inf(`1*DIGIT`, core.DIGIT()),
		Concat(
			`*DIGIT "*" *DIGIT`,
			Repeat0Inf(`*DIGIT`, core.DIGIT()),
			Rune(`*`, '*'),
			Repeat0Inf(`*DIGIT`, core.DIGIT()),
		),
	)(s)
}

func element(s []rune) Alternatives {
	return Alts(
		`element`,
		ruleName,
		group,
		option,
		charVal,
		numVal,
		proseVal,
	)(s)
}

func group(s []rune) Alternatives {
	return Concat(
		`group`,
		Rune(`(`, '('),
		repeatCWsp,
		alternation,
		repeatCWsp,
		Rune(`(`, ')'),
	)(s)
}

func option(s []rune) Alternatives {
	return Concat(
		`option`,
		Rune(`[`, '['),
		repeatCWsp,
		alternation,
		repeatCWsp,
		Rune(`]`, ']'),
	)(s)
}

func charVal(s []rune) Alternatives {
	return Concat(
		`CHAR-val`,
		core.DQUOTE(),
		Repeat0Inf(`*(%x20-21 / %x23-7E)`, Alts(
			`%x20-21 / %x23-7E`,
			Range(`%x20-21`, '\x20', '\x21'),
			Range(`%x23-7E`, '\x23', '\x7E'),
		)),
		core.DQUOTE(),
	)(s)
}

func numVal(s []rune) Alternatives {
	return Concat(
		`num-val`,
		Rune(`%`, '%'),
		Alts(`bin-val / dec-val / hex-val`, binVal, decVal, hexVal))(s)
}

func binVal(s []rune) Alternatives {
	return Concat(`bin-val`,
		Rune(`b`, 'b'),
		Repeat1Inf(`1*BIT`, core.BIT()),
		Optional(`[1*("." 1*BIT) / ("-" 1*BIT)]`, Alts(
			`1*("." 1*BIT) / ("-" 1*BIT)`,
			Repeat1Inf(`1*("." 1*BIT)`, Concat(`"." 1*BIT`, Rune(`.`, '.'), Repeat1Inf(`1*BIT`, core.BIT()))),
			Concat(`"-" 1*BIT`, Rune(`-`, '-'), Repeat1Inf(`1*BIT`, core.BIT())),
		)),
	)(s)
}

func decVal(s []rune) Alternatives {
	return Concat(`dec-val`,
		Rune(`d`, 'd'),
		Repeat1Inf(`1*DIGIT`, core.DIGIT()),
		Optional(`[1*("." 1*DIGIT) / ("-" 1*DIGIT)]`, Alts(
			`1*("." 1*DIGIT) / ("-" 1*DIGIT)`,
			Repeat1Inf(`1*("." 1*DIGIT)`, Concat(`"." 1*DIGIT`, Rune(`.`, '.'), Repeat1Inf(`1*DIGIT`, core.DIGIT()))),
			Concat(`"-" 1*DIGIT`, Rune(`-`, '-'), Repeat1Inf(`1*DIGIT`, core.DIGIT())),
		)),
	)(s)
}

func hexVal(s []rune) Alternatives {
	return Concat(`hex-val`,
		Rune(`x`, 'x'),
		Repeat1Inf(`1*HEXDIG`, core.HEXDIG()),
		Optional(`[1*("." 1*HEXDIG) / ("-" 1*HEXDIG)]`, Alts(
			`1*("." 1*HEXDIG) / ("-" 1*HEXDIG)`,
			Repeat1Inf(`1*("." 1*HEXDIG)`, Concat(`"." 1*HEXDIG`, Rune(`.`, '.'), Repeat1Inf(`1*HEXDIG`, core.HEXDIG()))),
			Concat(`"-" 1*HEXDIG`, Rune(`-`, '-'), Repeat1Inf(`1*HEXDIG`, core.HEXDIG())),
		)),
	)(s)
}

func proseVal(s []rune) Alternatives {
	return Concat(
		`prose-val`,
		Rune(`<`, '<'),
		Repeat0Inf(`*(%x20-3D / %x3F-7E)`, Alts(
			`%x20-3D / %x3F-7E`,
			Range(`%x20-3D`, '\x20', '\x3D'),
			Range(`%x3F-7E`, '\x3F', '\x7E'),
		)),
		Rune(`>`, '>'),
	)(s)
}
