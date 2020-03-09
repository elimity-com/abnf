package core

import . "github.com/di-wu/abnf"

// RFC 5234 - Appendix B: Core ABNF of ABNF
// Certain basic rules are in uppercase, such as SP, HTAB, CRLF, DIGIT, ALPHA, etc.

// ALPHA = %x41-5A / %x61-7A ; A-Z / a-z
func ALPHA() Rule {
	return Alts(Range('\x41', '\x5A'), Range('\x61', '\x7A'))
}

// BIT = "0" / "1"
func BIT() Rule {
	return Runes('0', '1')
}

// CHAR = %x01-7F ; any 7-bit US-ASCII character, excluding NUL
func CHAR() Rule {
	return Range('\x01', '\x7F')
}

// CR = %x0D ; carriage return
func CR() Rule {
	return Rune('\x0D')
}

// CRLF = CR LF ; Internet standard newline
func CRLF() Rule {
	return Concat(CR(), LF())
}

// CTL = %x00-1F / %x7F ; controls
func CTL() Rule {
	return Alts(Range('\x00', '\x1F'), Rune('\x7F'))
}

// DIGIT = %x30-39 ; 0-9
func DIGIT() Rule {
	return Range('\x30', '\x39')
}

// DQUOTE = %x22 ; " (Double Quote)
func DQUOTE() Rule {
	return Rune('\x22')
}

// HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
func HEXDIG() Rule {
	return Alts(DIGIT(), Rune('A'), Rune('B'), Rune('C'), Rune('D'), Rune('E'), Rune('F'))
}

// HTAB = %x09 ; horizontal tab
func HTAB() Rule {
	return Rune('\x09')
}

// LF = %x0A ; linefeed
func LF() Rule {
	return Rune('\x0A')
}

// LWSP = *(WSP / CRLF WSP) ; Use of this linear-white-space rule permits lines containing only white space that are
// no longer legal in mail headers and have caused interoperability problems in other contexts. Do not use when
// defining mail headers and use with caution in other contexts.
func LWSP() Rule {
	return DefaultRepeat(Alts(WSP(), Concat(CRLF(), WSP())))
}

// OCTET = %x00-FF ; 8 bits of data
func OCTET() Rule {
	return Range('\x00', '\xFF')
}

// SP = %x20
func SP() Rule {
	return Rune('\x20')
}

// VCHAR = %x21-7E ; visible (printing) characters
func VCHAR() Rule {
	return Range('\x21', '\x7E')
}

// WSP = SP / HTAB ; white space
func WSP() Rule {
	return Alts(SP(), HTAB())
}
