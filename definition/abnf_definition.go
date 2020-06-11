// This file is generated - do not edit.

package definition

import (
	core "github.com/elimity-com/abnf/core"
	operators "github.com/elimity-com/abnf/operators"
)

// Rulelist = 1*( rule / (*WSP c-nl) )
func Rulelist(s []rune) operators.Alternatives {
	return operators.Repeat1Inf("rulelist", operators.Alts(
		"rule / (*WSP c-nl)",
		Rule,
		operators.Concat(
			"(*WSP c-nl)",
			operators.Repeat0Inf("*WSP", core.WSP()),
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
		core.ALPHA(),
		operators.Repeat0Inf("*(ALPHA / DIGIT / \"-\")", operators.Alts(
			"ALPHA / DIGIT / \"-\"",
			core.ALPHA(),
			core.DIGIT(),
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
		operators.Repeat0Inf("*WSP", core.WSP()),
	)(s)
}

// CWsp = WSP / (c-nl WSP)
func CWsp(s []rune) operators.Alternatives {
	return operators.Alts(
		"c-wsp",
		core.WSP(),
		operators.Concat(
			"(c-nl WSP)",
			CNl,
			core.WSP(),
		),
	)(s)
}

// CNl = comment / CRLF
func CNl(s []rune) operators.Alternatives {
	return operators.Alts(
		"c-nl",
		Comment,
		core.CRLF(),
	)(s)
}

// Comment = "
func Comment(s []rune) operators.Alternatives {
	return operators.Concat(
		"comment",
		operators.Rune("\";\"", 59),
		operators.Repeat0Inf("*(WSP / VCHAR)", operators.Alts(
			"WSP / VCHAR",
			core.WSP(),
			core.VCHAR(),
		)),
		core.CRLF(),
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
		operators.Repeat1Inf("1*DIGIT", core.DIGIT()),
		operators.Concat(
			"(*DIGIT \"*\" *DIGIT)",
			operators.Repeat0Inf("*DIGIT", core.DIGIT()),
			operators.Rune("\"*\"", 42),
			operators.Repeat0Inf("*DIGIT", core.DIGIT()),
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
		core.DQUOTE(),
		operators.Repeat0Inf("*(%x20-21 / %x23-7E)", operators.Alts(
			"%x20-21 / %x23-7E",
			operators.Range("%x20-21", 32, 33),
			operators.Range("%x23-7E", 35, 126),
		)),
		core.DQUOTE(),
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
		operators.Repeat1Inf("1*BIT", core.BIT()),
		operators.Optional("[ 1*(\".\" 1*BIT) / (\"-\" 1*BIT) ]", operators.Alts(
			"1*(\".\" 1*BIT) / (\"-\" 1*BIT)",
			operators.Repeat1Inf("1*(\".\" 1*BIT)", operators.Concat(
				"\".\" 1*BIT",
				operators.Rune("\".\"", 46),
				operators.Repeat1Inf("1*BIT", core.BIT()),
			)),
			operators.Concat(
				"(\"-\" 1*BIT)",
				operators.Rune("\"-\"", 45),
				operators.Repeat1Inf("1*BIT", core.BIT()),
			),
		)),
	)(s)
}

// DecVal = "d" 1*DIGIT [ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]
func DecVal(s []rune) operators.Alternatives {
	return operators.Concat(
		"dec-val",
		operators.Rune("\"d\"", 100),
		operators.Repeat1Inf("1*DIGIT", core.DIGIT()),
		operators.Optional("[ 1*(\".\" 1*DIGIT) / (\"-\" 1*DIGIT) ]", operators.Alts(
			"1*(\".\" 1*DIGIT) / (\"-\" 1*DIGIT)",
			operators.Repeat1Inf("1*(\".\" 1*DIGIT)", operators.Concat(
				"\".\" 1*DIGIT",
				operators.Rune("\".\"", 46),
				operators.Repeat1Inf("1*DIGIT", core.DIGIT()),
			)),
			operators.Concat(
				"(\"-\" 1*DIGIT)",
				operators.Rune("\"-\"", 45),
				operators.Repeat1Inf("1*DIGIT", core.DIGIT()),
			),
		)),
	)(s)
}

// HexVal = "x" 1*HEXDIG [ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]
func HexVal(s []rune) operators.Alternatives {
	return operators.Concat(
		"hex-val",
		operators.Rune("\"x\"", 120),
		operators.Repeat1Inf("1*HEXDIG", core.HEXDIG()),
		operators.Optional("[ 1*(\".\" 1*HEXDIG) / (\"-\" 1*HEXDIG) ]", operators.Alts(
			"1*(\".\" 1*HEXDIG) / (\"-\" 1*HEXDIG)",
			operators.Repeat1Inf("1*(\".\" 1*HEXDIG)", operators.Concat(
				"\".\" 1*HEXDIG",
				operators.Rune("\".\"", 46),
				operators.Repeat1Inf("1*HEXDIG", core.HEXDIG()),
			)),
			operators.Concat(
				"(\"-\" 1*HEXDIG)",
				operators.Rune("\"-\"", 45),
				operators.Repeat1Inf("1*HEXDIG", core.HEXDIG()),
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
