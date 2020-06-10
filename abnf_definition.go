// This file is generated - do not edit.

package abnf

import operators "github.com/elimity-com/abnf/operators"

// Rulelist = 1*( rule / (*WSP c-nl) )
func Rulelist(s []rune) operators.Alternatives {
	return operators.Repeat1Inf("rulelist", operators.Alts(
		"rule / (*WSP c-nl)",
		Rule,
		operators.Concat(
			"(*WSP c-nl)",
			operators.Repeat0Inf("*WSP", WSP),
			CNl,
		),
	))(s)
}

// Rule = rulename defined-as elements c-nl
func Rule(s []rune) operators.Alternatives {
	return operators.Concat(
		"rule",
		Rulename,
		DefinedAs,
		Elements,
		CNl,
	)(s)
}

// Rulename = ALPHA *(ALPHA / DIGIT / "-")
func Rulename(s []rune) operators.Alternatives {
	return operators.Concat(
		"rulename",
		ALPHA,
		operators.Repeat0Inf("*(ALPHA / DIGIT / \"-\")", operators.Alts(
			"ALPHA / DIGIT / \"-\"",
			ALPHA,
			DIGIT,
			operators.Rune("\"-\"", 45),
		)),
	)(s)
}

// DefinedAs = *c-wsp ("=" / "=/") *c-wsp
func DefinedAs(s []rune) operators.Alternatives {
	return operators.Concat(
		"defined-as",
		operators.Repeat0Inf("*c-wsp", CWsp),
		operators.Alts(
			"(\"=\" / \"=/\")",
			operators.Rune("\"=\"", 61),
			operators.String("\"=/\"", "=/"),
		),
		operators.Repeat0Inf("*c-wsp", CWsp),
	)(s)
}

// Elements = alternation *WSP
func Elements(s []rune) operators.Alternatives {
	return operators.Concat(
		"elements",
		Alternation,
		operators.Repeat0Inf("*WSP", WSP),
	)(s)
}

// CWsp = WSP / (c-nl WSP)
func CWsp(s []rune) operators.Alternatives {
	return operators.Alts(
		"c-wsp",
		WSP,
		operators.Concat(
			"(c-nl WSP)",
			CNl,
			WSP,
		),
	)(s)
}

// CNl = comment / CRLF
func CNl(s []rune) operators.Alternatives {
	return operators.Alts(
		"c-nl",
		Comment,
		CRLF,
	)(s)
}

// Comment = "
func Comment(s []rune) operators.Alternatives {
	return operators.Concat(
		"comment",
		operators.Rune("\";\"", 59),
		operators.Repeat0Inf("*(WSP / VCHAR)", operators.Alts(
			"WSP / VCHAR",
			WSP,
			VCHAR,
		)),
		CRLF,
	)(s)
}

// Alternation = concatenation *(*c-wsp "/" *c-wsp concatenation)
func Alternation(s []rune) operators.Alternatives {
	return operators.Concat(
		"alternation",
		Concatenation,
		operators.Repeat0Inf("*(*c-wsp \"/\" *c-wsp concatenation)", operators.Concat(
			"*c-wsp \"/\" *c-wsp concatenation",
			operators.Repeat0Inf("*c-wsp", CWsp),
			operators.Rune("\"/\"", 47),
			operators.Repeat0Inf("*c-wsp", CWsp),
			Concatenation,
		)),
	)(s)
}

// Concatenation = repetition *(1*c-wsp repetition)
func Concatenation(s []rune) operators.Alternatives {
	return operators.Concat(
		"concatenation",
		Repetition,
		operators.Repeat0Inf("*(1*c-wsp repetition)", operators.Concat(
			"1*c-wsp repetition",
			operators.Repeat1Inf("1*c-wsp", CWsp),
			Repetition,
		)),
	)(s)
}

// Repetition = [repeat] element
func Repetition(s []rune) operators.Alternatives {
	return operators.Concat(
		"repetition",
		operators.Optional("[repeat]", Repeat),
		Element,
	)(s)
}

// Repeat = 1*DIGIT / (*DIGIT "*" *DIGIT)
func Repeat(s []rune) operators.Alternatives {
	return operators.Alts(
		"repeat",
		operators.Repeat1Inf("1*DIGIT", DIGIT),
		operators.Concat(
			"(*DIGIT \"*\" *DIGIT)",
			operators.Repeat0Inf("*DIGIT", DIGIT),
			operators.Rune("\"*\"", 42),
			operators.Repeat0Inf("*DIGIT", DIGIT),
		),
	)(s)
}

// Element = rulename / group / option / char-val / num-val / prose-val
func Element(s []rune) operators.Alternatives {
	return operators.Alts(
		"element",
		Rulename,
		Group,
		Option,
		CharVal,
		NumVal,
		ProseVal,
	)(s)
}

// Group = "(" *c-wsp alternation *c-wsp ")"
func Group(s []rune) operators.Alternatives {
	return operators.Concat(
		"group",
		operators.Rune("\"(\"", 40),
		operators.Repeat0Inf("*c-wsp", CWsp),
		Alternation,
		operators.Repeat0Inf("*c-wsp", CWsp),
		operators.Rune("\")\"", 41),
	)(s)
}

// Option = "[" *c-wsp alternation *c-wsp "]"
func Option(s []rune) operators.Alternatives {
	return operators.Concat(
		"option",
		operators.Rune("\"[\"", 91),
		operators.Repeat0Inf("*c-wsp", CWsp),
		Alternation,
		operators.Repeat0Inf("*c-wsp", CWsp),
		operators.Rune("\"]\"", 93),
	)(s)
}

// CharVal = DQUOTE *(%x20-21 / %x23-7E) DQUOTE
func CharVal(s []rune) operators.Alternatives {
	return operators.Concat(
		"char-val",
		DQUOTE,
		operators.Repeat0Inf("*(%x20-21 / %x23-7E)", operators.Alts(
			"%x20-21 / %x23-7E",
			operators.Range("%x20-21", 32, 33),
			operators.Range("%x23-7E", 35, 126),
		)),
		DQUOTE,
	)(s)
}

// NumVal = "%" (bin-val / dec-val / hex-val)
func NumVal(s []rune) operators.Alternatives {
	return operators.Concat(
		"num-val",
		operators.Rune("\"%\"", 37),
		operators.Alts(
			"(bin-val / dec-val / hex-val)",
			BinVal,
			DecVal,
			HexVal,
		),
	)(s)
}

// BinVal = "b" 1*BIT [ 1*("." 1*BIT) / ("-" 1*BIT) ]
func BinVal(s []rune) operators.Alternatives {
	return operators.Concat(
		"bin-val",
		operators.Rune("\"b\"", 98),
		operators.Repeat1Inf("1*BIT", BIT),
		operators.Optional("[ 1*(\".\" 1*BIT) / (\"-\" 1*BIT) ]", operators.Alts(
			"1*(\".\" 1*BIT) / (\"-\" 1*BIT)",
			operators.Repeat1Inf("1*(\".\" 1*BIT)", operators.Concat(
				"\".\" 1*BIT",
				operators.Rune("\".\"", 46),
				operators.Repeat1Inf("1*BIT", BIT),
			)),
			operators.Concat(
				"(\"-\" 1*BIT)",
				operators.Rune("\"-\"", 45),
				operators.Repeat1Inf("1*BIT", BIT),
			),
		)),
	)(s)
}

// DecVal = "d" 1*DIGIT [ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]
func DecVal(s []rune) operators.Alternatives {
	return operators.Concat(
		"dec-val",
		operators.Rune("\"d\"", 100),
		operators.Repeat1Inf("1*DIGIT", DIGIT),
		operators.Optional("[ 1*(\".\" 1*DIGIT) / (\"-\" 1*DIGIT) ]", operators.Alts(
			"1*(\".\" 1*DIGIT) / (\"-\" 1*DIGIT)",
			operators.Repeat1Inf("1*(\".\" 1*DIGIT)", operators.Concat(
				"\".\" 1*DIGIT",
				operators.Rune("\".\"", 46),
				operators.Repeat1Inf("1*DIGIT", DIGIT),
			)),
			operators.Concat(
				"(\"-\" 1*DIGIT)",
				operators.Rune("\"-\"", 45),
				operators.Repeat1Inf("1*DIGIT", DIGIT),
			),
		)),
	)(s)
}

// HexVal = "x" 1*HEXDIG [ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]
func HexVal(s []rune) operators.Alternatives {
	return operators.Concat(
		"hex-val",
		operators.Rune("\"x\"", 120),
		operators.Repeat1Inf("1*HEXDIG", HEXDIG),
		operators.Optional("[ 1*(\".\" 1*HEXDIG) / (\"-\" 1*HEXDIG) ]", operators.Alts(
			"1*(\".\" 1*HEXDIG) / (\"-\" 1*HEXDIG)",
			operators.Repeat1Inf("1*(\".\" 1*HEXDIG)", operators.Concat(
				"\".\" 1*HEXDIG",
				operators.Rune("\".\"", 46),
				operators.Repeat1Inf("1*HEXDIG", HEXDIG),
			)),
			operators.Concat(
				"(\"-\" 1*HEXDIG)",
				operators.Rune("\"-\"", 45),
				operators.Repeat1Inf("1*HEXDIG", HEXDIG),
			),
		)),
	)(s)
}

// ProseVal = "<" *(%x20-3D / %x3F-7E) ">"
func ProseVal(s []rune) operators.Alternatives {
	return operators.Concat(
		"prose-val",
		operators.Rune("\"<\"", 60),
		operators.Repeat0Inf("*(%x20-3D / %x3F-7E)", operators.Alts(
			"%x20-3D / %x3F-7E",
			operators.Range("%x20-3D", 32, 61),
			operators.Range("%x3F-7E", 63, 126),
		)),
		operators.Rune("\">\"", 62),
	)(s)
}

// ALPHA = %x41-5A / %x61-7A
func ALPHA(s []rune) operators.Alternatives {
	return operators.Alts(
		"ALPHA",
		operators.Range("%x41-5A", 65, 90),
		operators.Range("%x61-7A", 97, 122),
	)(s)
}

// BIT = "0" / "1"
func BIT(s []rune) operators.Alternatives {
	return operators.Alts(
		"BIT",
		operators.Rune("\"0\"", 48),
		operators.Rune("\"1\"", 49),
	)(s)
}

// CHAR = %x01-7F
func CHAR(s []rune) operators.Alternatives {
	return operators.Range("CHAR", 1, 127)(s)
}

// CR = %x0D
func CR(s []rune) operators.Alternatives {
	return operators.Rune("CR", 13)(s)
}

// CRLF = CR LF / LF
func CRLF(s []rune) operators.Alternatives {
	return operators.Alts(
		"CRLF",
		operators.Concat(
			"CR LF",
			CR,
			LF,
		),
		LF,
	)(s)
}

// CTL = %x00-1F / %x7F
func CTL(s []rune) operators.Alternatives {
	return operators.Alts(
		"CTL",
		operators.Range("%x00-1F", 0, 31),
		operators.Rune("%x7F", 127),
	)(s)
}

// DIGIT = %x30-39
func DIGIT(s []rune) operators.Alternatives {
	return operators.Range("DIGIT", 48, 57)(s)
}

// DQUOTE = %x22
func DQUOTE(s []rune) operators.Alternatives {
	return operators.Rune("DQUOTE", 34)(s)
}

// HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
func HEXDIG(s []rune) operators.Alternatives {
	return operators.Alts(
		"HEXDIG",
		DIGIT,
		operators.Rune("\"A\"", 65),
		operators.Rune("\"B\"", 66),
		operators.Rune("\"C\"", 67),
		operators.Rune("\"D\"", 68),
		operators.Rune("\"E\"", 69),
		operators.Rune("\"F\"", 70),
	)(s)
}

// HTAB = %x09
func HTAB(s []rune) operators.Alternatives {
	return operators.Rune("HTAB", 9)(s)
}

// LF = %x0A
func LF(s []rune) operators.Alternatives {
	return operators.Rune("LF", 10)(s)
}

// LWSP = *(WSP / CRLF WSP)
func LWSP(s []rune) operators.Alternatives {
	return operators.Repeat0Inf("LWSP", operators.Alts(
		"WSP / CRLF WSP",
		WSP,
		operators.Concat(
			"CRLF WSP",
			CRLF,
			WSP,
		),
	))(s)
}

// OCTET = %x00-FF
func OCTET(s []rune) operators.Alternatives {
	return operators.Range("OCTET", 0, 255)(s)
}

// SP = %x20
func SP(s []rune) operators.Alternatives {
	return operators.Rune("SP", 32)(s)
}

// VCHAR = %x21-7E
func VCHAR(s []rune) operators.Alternatives {
	return operators.Range("VCHAR", 33, 126)(s)
}

// WSP = SP / HTAB
func WSP(s []rune) operators.Alternatives {
	return operators.Alts(
		"WSP",
		SP,
		HTAB,
	)(s)
}
