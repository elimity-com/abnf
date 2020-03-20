package abnf

// RFC 5234: 4. ABNF Definition of ABNF

func ruleList() RuleFunc {
	return Repeat1Inf(Alts(rule(), Concat(repeatCWsp(), cNl())))
}

func rule() RuleFunc {
	return Concat(ruleName(), definedAs(), elements(), cNl())
}

func ruleName() RuleFunc {
	return Concat(alpha(), Repeat0Inf(Alts(alpha(), digit(), Rune('-'))))
}

func definedAs() RuleFunc {
	return Concat(repeatCWsp(), Alts(Rune('='), Runes('=', '/')), repeatCWsp())
}

func elements() RuleFunc {
	return Concat(alternation(), repeatCWsp())
}

func cWsp() RuleFunc {
	return Alts(wsp(), Concat(cNl(), wsp()))
}

func repeatCWsp() RuleFunc {
	return Repeat0Inf(cWsp())
}

func cNl() RuleFunc {
	return Alts(comment(), crlf())
}

func comment() RuleFunc {
	return Concat(Rune(';'), Repeat0Inf(Alts(wsp(), vchar())), crlf())
}

func alternation() RuleFunc {
	return Concat(concatenation(), Repeat0Inf(Concat(repeatCWsp(), Rune('/'), repeatCWsp(), concatenation())))
}

func concatenation() RuleFunc {
	return Concat(repetition(), Repeat0Inf(Concat(Repeat1Inf(cWsp()), repetition())))
}

func repetition() RuleFunc {
	return Concat(Optional(repeat()), element())
}

func repeat() RuleFunc {
	return Alts(Repeat1Inf(digit()), Concat(Repeat0Inf(digit()), Rune('*'), Repeat0Inf(digit())))
}

func element() RuleFunc {
	return Alts(ruleName(), group(), option(), charVal(), numVal(), proseVal())
}

func group() RuleFunc {
	return Concat(Rune('('), repeatCWsp(), alternation(), repeatCWsp(), Rune(')'))
}

func option() RuleFunc {
	return Concat(Rune('['), repeatCWsp(), alternation(), repeatCWsp(), Rune(']'))
}

func charVal() RuleFunc {
	return Concat(dquote(), Repeat0Inf(Alts(Range('\x20', '\x21'), Range('\x23', '\x7E'))), dquote())
}

func numVal() RuleFunc {
	return Concat(Rune('%'), Alts(binVal(), decVal(), hexVal()))
}

func val(r rune, rule RuleFunc) RuleFunc {
	return Concat(Rune(r), Repeat1Inf(rule),
		Optional(Repeat1Inf(Alts(
			Rune('.'), Repeat1Inf(rule),
			Rune('-'), Repeat1Inf(rule),
		))),
	)
}

func binVal() RuleFunc {
	return val('b', bit())
}

func decVal() RuleFunc {
	return val('d', digit())
}

func hexVal() RuleFunc {
	return val('x', hexdig())
}

func proseVal() RuleFunc {
	return Concat(Rune('<'), Repeat0Inf(Alts(Range('\x20', '\x3D'), Range('\x3F', '\x7E'))), Rune('>'))
}
