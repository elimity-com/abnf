// This file is generated - do not edit.

package core

import "github.com/elimity-com/abnf/operators"

// ALPHA = %x41-5A / %x61-7A
func ALPHA() operators.Operator {
	return operators.Alts(
		"ALPHA",
		operators.Range("%x41-5A", []byte{65}, []byte{90}),
		operators.Range("%x61-7A", []byte{97}, []byte{122}),
	)
}

// BIT = "0" / "1"
func BIT() operators.Operator {
	return operators.Alts(
		"BIT",
		operators.String("0", "0"),
		operators.String("1", "1"),
	)
}

// CHAR = %x01-7F
func CHAR() operators.Operator {
	return operators.Range("CHAR", []byte{1}, []byte{127})
}

// CR = %x0D
func CR() operators.Operator {
	return operators.Terminal("CR", []byte{13})
}

// CRLF = CR LF / LF
func CRLF() operators.Operator {
	return operators.Alts(
		"CRLF",
		operators.Concat(
			"CR LF",
			CR(),
			LF(),
		),
		LF(),
	)
}

// CTL = %x00-1F / %x7F
func CTL() operators.Operator {
	return operators.Alts(
		"CTL",
		operators.Range("%x00-1F", []byte{0}, []byte{31}),
		operators.Terminal("%x7F", []byte{127}),
	)
}

// DIGIT = %x30-39
func DIGIT() operators.Operator {
	return operators.Range("DIGIT", []byte{48}, []byte{57})
}

// DQUOTE = %x22
func DQUOTE() operators.Operator {
	return operators.Terminal("DQUOTE", []byte{34})
}

// HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
func HEXDIG() operators.Operator {
	return operators.Alts(
		"HEXDIG",
		DIGIT(),
		operators.String("A", "A"),
		operators.String("B", "B"),
		operators.String("C", "C"),
		operators.String("D", "D"),
		operators.String("E", "E"),
		operators.String("F", "F"),
	)
}

// HTAB = %x09
func HTAB() operators.Operator {
	return operators.Terminal("HTAB", []byte{9})
}

// LF = %x0A
func LF() operators.Operator {
	return operators.Terminal("LF", []byte{10})
}

// LWSP = *(WSP / CRLF WSP)
func LWSP() operators.Operator {
	return operators.Repeat0Inf("LWSP", operators.Alts(
		"WSP / CRLF WSP",
		WSP(),
		operators.Concat(
			"CRLF WSP",
			CRLF(),
			WSP(),
		),
	))
}

// OCTET = %x00-FF
func OCTET() operators.Operator {
	return operators.Range("OCTET", []byte{0}, []byte{255})
}

// SP = %x20
func SP() operators.Operator {
	return operators.Terminal("SP", []byte{32})
}

// VCHAR = %x21-7E
func VCHAR() operators.Operator {
	return operators.Range("VCHAR", []byte{33}, []byte{126})
}

// WSP = SP / HTAB
func WSP() operators.Operator {
	return operators.Alts(
		"WSP",
		SP(),
		HTAB(),
	)
}
