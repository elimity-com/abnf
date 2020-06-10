// This file is generated - do not edit.

package core

import operators "github.com/elimity-com/abnf/operators"

// ALPHA = %x41-5A / %x61-7A
func ALPHA() operators.Operator {
	return operators.Alts("ALPHA", operators.Range("%x41-5A", 65, 90), operators.Range("%x61-7A", 97, 122))
}

// BIT = "0" / "1"
func BIT() operators.Operator {
	return operators.Alts("BIT", operators.Rune("\"0\"", 48), operators.Rune("\"1\"", 49))
}

// CHAR = %x01-7F
func CHAR() operators.Operator {
	return operators.Range("CHAR", 1, 127)
}

// CR = %x0D
func CR() operators.Operator {
	return operators.Rune("CR", 13)
}

// CRLF = CR LF / LF
func CRLF() operators.Operator {
	return operators.Alts("CRLF", operators.Concat("CR LF", CR(), LF()), LF())
}

// CTL = %x00-1F / %x7F
func CTL() operators.Operator {
	return operators.Alts("CTL", operators.Range("%x00-1F", 0, 31), operators.Rune("%x7F", 127))
}

// DIGIT = %x30-39
func DIGIT() operators.Operator {
	return operators.Range("DIGIT", 48, 57)
}

// DQUOTE = %x22
func DQUOTE() operators.Operator {
	return operators.Rune("DQUOTE", 34)
}

// HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
func HEXDIG() operators.Operator {
	return operators.Alts("HEXDIG", DIGIT(), operators.Rune("\"A\"", 65), operators.Rune("\"B\"", 66), operators.Rune("\"C\"", 67), operators.Rune("\"D\"", 68), operators.Rune("\"E\"", 69), operators.Rune("\"F\"", 70))
}

// HTAB = %x09
func HTAB() operators.Operator {
	return operators.Rune("HTAB", 9)
}

// LF = %x0A
func LF() operators.Operator {
	return operators.Rune("LF", 10)
}

// LWSP = *(WSP / CRLF WSP)
func LWSP() operators.Operator {
	return operators.Repeat0Inf("LWSP", operators.Alts("WSP / CRLF WSP", WSP(), operators.Concat("CRLF WSP", CRLF(), WSP())))
}

// OCTET = %x00-FF
func OCTET() operators.Operator {
	return operators.Range("OCTET", 0, 255)
}

// SP = %x20
func SP() operators.Operator {
	return operators.Rune("SP", 32)
}

// VCHAR = %x21-7E
func VCHAR() operators.Operator {
	return operators.Range("VCHAR", 33, 126)
}

// WSP = SP / HTAB
func WSP() operators.Operator {
	return operators.Alts("WSP", SP(), HTAB())
}
