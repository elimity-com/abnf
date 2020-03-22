package abnf

// RFC 5234: 4. ABNF Definition of ABNF

func ruleList(s *Scanner) *AST {
	return Repeat1Inf(`rulelist`, Alts(
		`rule / (*c-wsp c-nl)`,
		rule,
		Concat(
			`*c-wsp c-nl`,
			repeatCWsp,
			cNl,
		),
	))(s)
}

func rule(s *Scanner) *AST {
	return Concat(
		`rule`,
		ruleName,
		definedAs,
		elements,
		cNl,
	)(s)
}

func ruleName(s *Scanner) *AST {
	return Concat(
		`rulename`,
		alpha(),
		Repeat0Inf(`*(ALPHA / DIGIT / "-")`, Alts(
			`ALPHA / DIGIT / "-"`,
			alpha(),
			digit(),
			Rune(`-`, '-'),
		)),
	)(s)
}

func definedAs(s *Scanner) *AST {
	return Concat(
		`defined-as`,
		repeatCWsp,
		Alts(
			`"=" / "=/"`,
			Rune(`=`, '='),
			Runes(`=/`, '=', '/'),
		),
		repeatCWsp,
	)(s)
}

func elements(s *Scanner) *AST {
	return Concat(
		`elements`,
		alternation,
		repeatCWsp,
	)(s)
}

func cWsp(s *Scanner) *AST {
	return Alts(
		`c-wsp`,
		wsp(),
		Concat(
			`c-nl WSP`,
			cNl,
			wsp(),
		),
	)(s)
}

func repeatCWsp(s *Scanner) *AST {
	return Repeat0Inf(`*c-wsp`, cWsp)(s)
}

func cNl(s *Scanner) *AST {
	return Alts(
		`c-nl`,
		comment,
		crlf(),
	)(s)
}

func comment(s *Scanner) *AST {
	return Concat(
		`comment`,
		Rune(`;`, ';'),
		Repeat0Inf(`*(WSP / VCHAR) CRLF`, Alts(
			`WSP / VCHAR`,
			wsp(),
			vchar(),
		)),
		crlf(),
	)(s)
}

func alternation(s *Scanner) *AST {
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

func concatenation(s *Scanner) *AST {
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

func repetition(s *Scanner) *AST {
	return Concat(
		`repetition`,
		Optional(`[repeat]`, repeat),
		element,
	)(s)
}

func repeat(s *Scanner) *AST {
	return Alts(
		`repeat`,
		Repeat1Inf(`1*DIGIT`, digit()),
		Concat(
			`*DIGIT "*" *DIGIT`,
			Repeat0Inf(`*DIGIT`, digit()),
			Rune(`*`, '*'),
			Repeat0Inf(`*DIGIT`, digit()),
		),
	)(s)
}

func element(s *Scanner) *AST {
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

func group(s *Scanner) *AST {
	return Concat(
		`group`,
		Rune(`(`, '('),
		repeatCWsp,
		alternation,
		repeatCWsp,
		Rune(`(`, ')'),
	)(s)
}

func option(s *Scanner) *AST {
	return Concat(
		`option`,
		Rune(`[`, '['),
		repeatCWsp,
		alternation,
		repeatCWsp,
		Rune(`]`, ']'),
	)(s)
}

func charVal(s *Scanner) *AST {
	return Concat(
		`char-val`,
		dquote(),
		Repeat0Inf(`*(%x20-21 / %x23-7E)`, Alts(
			`%x20-21 / %x23-7E`,
			Range(`%x20-21`, '\x20', '\x21'),
			Range(`%x23-7E`, '\x23', '\x7E'),
		)),
		dquote(),
	)(s)
}

func numVal(s *Scanner) *AST {
	return Concat(
		`num-val`,
		Rune(`%`, '%'),
		Alts(`bin-val / dec-val / hex-val`, binVal, decVal, hexVal))(s)
}

func binVal(s *Scanner) *AST {
	return Concat(`bin-val`,
		Rune(`b`, 'b'),
		Repeat1Inf(`1*BIT`, bit()),
		Optional(`[1*("." 1*BIT) / ("-" 1*BIT)]`, Alts(
			`1*("." 1*BIT) / ("-" 1*BIT)`,
			Repeat1Inf(`1*("." 1*BIT)`, Concat(`"." 1*BIT`, Rune(`.`, '.'), Repeat1Inf(`1*BIT`, bit()))),
			Rune(`-`, '-'), Repeat1Inf(`1*BIT`, bit()),
		)),
	)(s)
}

func decVal(s *Scanner) *AST {
	return Concat(`dec-val`,
		Rune(`d`, 'd'),
		Repeat1Inf(`1*DIGIT`, digit()),
		Optional(`[1*("." 1*DIGIT) / ("-" 1*DIGIT)]`, Alts(
			`1*("." 1*DIGIT) / ("-" 1*DIGIT)`,
			Repeat1Inf(`1*("." 1*DIGIT)`, Concat(`"." 1*DIGIT`, Rune(`.`, '.'), Repeat1Inf(`1*DIGIT`, digit()))),
			Rune(`-`, '-'), Repeat1Inf(`1*DIGIT`, digit()),
		)),
	)(s)
}

func hexVal(s *Scanner) *AST {
	return Concat(`hex-val`,
		Rune(`x`, 'x'),
		Repeat1Inf(`1*HEXDIG`, hexdig()),
		Optional(`[1*("." 1*HEXDIG) / ("-" 1*HEXDIG)]`, Alts(
			`1*("." 1*HEXDIG) / ("-" 1*HEXDIG)`,
			Repeat1Inf(`1*("." 1*HEXDIG)`, Concat(`"." 1*HEXDIG`, Rune(`.`, '.'), Repeat1Inf(`1*HEXDIG`, hexdig()))),
			Rune(`-`, '-'), Repeat1Inf(`1*HEXDIG`, hexdig()),
		)),
	)(s)
}

func proseVal(s *Scanner) *AST {
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
