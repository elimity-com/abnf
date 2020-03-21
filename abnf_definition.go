package abnf

// RFC 5234: 4. ABNF Definition of ABNF

func ruleList(s *Scanner) []rune {
	return Repeat1Inf(Alts(
		rule,
		Concat(
			repeatCWsp,
			cNl,
		),
	))(s)
}

func rule(s *Scanner) []rune {
	return Concat(
		ruleName,
		definedAs,
		elements,
		cNl,
	)(s)
}

func ruleName(s *Scanner) []rune {
	return Concat(
		alpha(),
		Repeat0Inf(Alts(
			alpha(),
			digit(),
			Rune('-'),
		)),
	)(s)
}

func definedAs(s *Scanner) []rune {
	return Concat(
		repeatCWsp,
		Alts(
			Runes('=', '/'),
			Rune('='),
		),
		repeatCWsp,
	)(s)
}

func elements(s *Scanner) []rune {
	return Concat(
		alternation,
		repeatCWsp,
	)(s)
}

func cWsp(s *Scanner) []rune {
	return Alts(
		wsp(),
		Concat(
			cNl,
			wsp(),
		),
	)(s)
}

func repeatCWsp(s *Scanner) []rune {
	return Repeat0Inf(cWsp)(s)
}

func cNl(s *Scanner) []rune {
	return Alts(
		comment,
		crlf(),
	)(s)
}

func comment(s *Scanner) []rune {
	return Concat(
		Rune(';'),
		Repeat0Inf(Alts(
			wsp(),
			vchar(),
		)),
		crlf(),
	)(s)
}

func alternation(s *Scanner) []rune {
	return Concat(
		concatenation,
		Repeat0Inf(Concat(
			repeatCWsp,
			Rune('/'),
			repeatCWsp,
			concatenation,
		)),
	)(s)
}

func concatenation(s *Scanner) []rune {
	return Concat(
		repetition,
		Repeat0Inf(Concat(
			Repeat1Inf(cWsp),
			repetition,
		)),
	)(s)
}

func repetition(s *Scanner) []rune {
	return Concat(
		Optional(repeat),
		element,
	)(s)
}

func repeat(s *Scanner) []rune {
	return Alts(
		Repeat1Inf(digit()),
		Concat(
			Repeat0Inf(digit()), Rune('*'),
			Repeat0Inf(digit()),
		),
	)(s)
}

func element(s *Scanner) []rune {
	return Alts(
		ruleName,
		group,
		option,
		charVal,
		numVal,
		proseVal,
	)(s)
}

func group(s *Scanner) []rune {
	return Concat(
		Rune('('),
		repeatCWsp,
		alternation,
		repeatCWsp,
		Rune(')'),
	)(s)
}

func option(s *Scanner) []rune {
	return Concat(
		Rune('['),
		repeatCWsp,
		alternation,
		repeatCWsp,
		Rune(']'),
	)(s)
}

func charVal(s *Scanner) []rune {
	return Concat(
		dquote(),
		Repeat0Inf(Alts(
			Range('\x20', '\x21'),
			Range('\x23', '\x7E'),
		)),
		dquote(),
	)(s)
}

func numVal(s *Scanner) []rune {
	return Concat(Rune('%'), Alts(binVal, decVal, hexVal))(s)
}

func val(r rune, rule Operator) Operator {
	return Concat(Rune(r), Repeat1Inf(rule),
		Optional(Repeat1Inf(Alts(
			Rune('.'), Repeat1Inf(rule),
			Rune('-'), Repeat1Inf(rule),
		))),
	)
}

func binVal(s *Scanner) []rune {
	return val('b', bit())(s)
}

func decVal(s *Scanner) []rune {
	return val('d', digit())(s)
}

func hexVal(s *Scanner) []rune {
	return val('x', hexdig())(s)
}

func proseVal(s *Scanner) []rune {
	return Concat(
		Rune('<'),
		Repeat0Inf(Alts(
			Range('\x20', '\x3D'),
			Range('\x3F', '\x7E'),
		)),
		Rune('>'),
	)(s)
}
