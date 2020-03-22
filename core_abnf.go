package abnf

// RFC 5234: Appendix B. Core ABNF of ABNF

// ALPHA = %x41-5A / %x61-7A ; A-Z / a-z
func alpha() Operator {
	return Alts(`ALPHA`, Range(`%x41-5A`, '\x41', '\x5A'), Range(`%x61-7A`, '\x61', '\x7A'))
}

// BIT = "0" / "1"
func bit() Operator {
	return Runes(`BIT`, '0', '1')
}

// CHAR = %x01-7F ; any 7-bit US-ASCII character, excluding NUL
func char() Operator {
	return Range(`CHAR`, '\x01', '\x7F')
}

// CR = %x0D ; carriage return
func cr() Operator {
	return Rune(`CR`, '\x0D')
}

// CRLF = CR LF ; Internet standard newline
func crlf() Operator {
	// TODO: abnf only allows \r\n, yet this not not practical because unix only uses \n.
	return Alts(`CRLF`, Concat(`CR LF`, cr(), lf()), lf())
}

// CTL = %x00-1F / %x7F ; controls
func ctl() Operator {
	return Alts(`CTL`, Range(`%x00-1F`, '\x00', '\x1F'), Rune(`%x7F`, '\x7F'))
}

// DIGIT = %x30-39 ; 0-9
func digit() Operator {
	return Range(`DIGIT`, '\x30', '\x39')
}

// DQUOTE = %x22 ; " (Double Quote)
func dquote() Operator {
	return Rune(`DQUOTE`, '\x22')
}

// HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
func hexdig() Operator {
	return Alts(`HEXDIG`, digit(),
		Rune(`A`, 'A'), Rune(`B`, 'B'), Rune(`C`, 'C'),
		Rune(`D`, 'D'), Rune(`E`, 'E'), Rune(`F`, 'F'),
	)
}

// HTAB = %x09 ; horizontal tab
func htab() Operator {
	return Rune(`HTAB`, '\x09')
}

// LF = %x0A ; linefeed
func lf() Operator {
	return Rune(`LF`, '\x0A')
}

// LWSP = *(WSP / CRLF WSP) ; Use of this linear-white-space rule permits lines containing only white space that are
// no longer legal in mail headers and have caused interoperability problems in other contexts. Do not use when
// defining mail headers and use with caution in other contexts.
func lwsp() Operator {
	return Repeat0Inf(`LWSP`, Alts(`WSP / CRLF WSP`, wsp(), Concat(`CRLF WSP`, crlf(), wsp())))
}

// OCTET = %x00-FF ; 8 bits of data
func octet() Operator {
	return Range(`OCTET`, '\x00', '\xFF')
}

// SP = %x20
func sp() Operator {
	return Rune(`SP`, '\x20')
}

// VCHAR = %x21-7E ; visible (printing) characters
func vchar() Operator {
	return Range(`VCHAR`, '\x21', '\x7E')
}

// WSP = SP / HTAB ; white space
func wsp() Operator {
	return Alts(`WSP`, sp(), htab())
}
