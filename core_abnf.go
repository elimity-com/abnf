package abnf

// RFC 5234: Appendix B. Core ABNF of ABNF

// ALPHA = %x41-5A / %x61-7A ; A-Z / a-z
func alpha() Operator {
	return Alts(Range('\x41', '\x5A'), Range('\x61', '\x7A'))
}

// BIT = "0" / "1"
func bit() Operator {
	return Runes('0', '1')
}

// CHAR = %x01-7F ; any 7-bit US-ASCII character, excluding NUL
func char() Operator {
	return Range('\x01', '\x7F')
}

// CR = %x0D ; carriage return
func cr() Operator {
	return Rune('\x0D')
}

// CRLF = CR LF ; Internet standard newline
func crlf() Operator {
	return Concat(cr(), lf())
}

// CTL = %x00-1F / %x7F ; controls
func ctl() Operator {
	return Alts(Range('\x00', '\x1F'), Rune('\x7F'))
}

// DIGIT = %x30-39 ; 0-9
func digit() Operator {
	return Range('\x30', '\x39')
}

// DQUOTE = %x22 ; " (Double Quote)
func dquote() Operator {
	return Rune('\x22')
}

// HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
func hexdig() Operator {
	return Alts(digit(), Rune('A'), Rune('B'), Rune('C'), Rune('D'), Rune('E'), Rune('F'))
}

// HTAB = %x09 ; horizontal tab
func htab() Operator {
	return Rune('\x09')
}

// LF = %x0A ; linefeed
func lf() Operator {
	return Rune('\x0A')
}

// LWSP = *(WSP / CRLF WSP) ; Use of this linear-white-space rule permits lines containing only white space that are
// no longer legal in mail headers and have caused interoperability problems in other contexts. Do not use when
// defining mail headers and use with caution in other contexts.
func lwsp() Operator {
	return Repeat0Inf(Alts(wsp(), Concat(crlf(), wsp())))
}

// OCTET = %x00-FF ; 8 bits of data
func octet() Operator {
	return Range('\x00', '\xFF')
}

// SP = %x20
func sp() Operator {
	return Rune('\x20')
}

// VCHAR = %x21-7E ; visible (printing) characters
func vchar() Operator {
	return Range('\x21', '\x7E')
}

// WSP = SP / HTAB ; white space
func wsp() Operator {
	return Alts(sp(), htab())
}
