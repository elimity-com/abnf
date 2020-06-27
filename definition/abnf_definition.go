// This file is generated - do not edit.

package definition

import (
	core "github.com/elimity-com/abnf/core"
	"github.com/elimity-com/abnf/operators"
)

// alternation = concatenation *(*c-wsp "/" *c-wsp concatenation)
func Alternation(s []byte) operators.Alternatives {
	return operators.Concat(
		"alternation",
		Concatenation,
		operators.Repeat0Inf("*(*c-wsp \"/\" *c-wsp concatenation)", operators.Concat(
			"*c-wsp \"/\" *c-wsp concatenation",
			operators.Repeat0Inf("*c-wsp", CWsp),
			operators.String("/", "/"),
			operators.Repeat0Inf("*c-wsp", CWsp),
			Concatenation,
		)),
	)(s)
}

// bin-val = "b" 1*BIT [ 1*("." 1*BIT) / ("-" 1*BIT) ]
func BinVal(s []byte) operators.Alternatives {
	return operators.Concat(
		"bin-val",
		operators.String("b", "b"),
		operators.Repeat1Inf("1*BIT", core.BIT()),
		operators.Optional("[ 1*(\".\" 1*BIT) / (\"-\" 1*BIT) ]", operators.Alts(
			"1*(\".\" 1*BIT) / (\"-\" 1*BIT)",
			operators.Repeat1Inf("1*(\".\" 1*BIT)", operators.Concat(
				"\".\" 1*BIT",
				operators.String(".", "."),
				operators.Repeat1Inf("1*BIT", core.BIT()),
			)),
			operators.Concat(
				"\"-\" 1*BIT",
				operators.String("-", "-"),
				operators.Repeat1Inf("1*BIT", core.BIT()),
			),
		)),
	)(s)
}

// c-nl = comment / CRLF
func CNl(s []byte) operators.Alternatives {
	return operators.Alts(
		"c-nl",
		Comment,
		core.CRLF(),
	)(s)
}

// c-wsp = WSP / (c-nl WSP)
func CWsp(s []byte) operators.Alternatives {
	return operators.Alts(
		"c-wsp",
		core.WSP(),
		operators.Concat(
			"c-nl WSP",
			CNl,
			core.WSP(),
		),
	)(s)
}

// char-val = DQUOTE *(%x20-21 / %x23-7E) DQUOTE
func CharVal(s []byte) operators.Alternatives {
	return operators.Concat(
		"char-val",
		core.DQUOTE(),
		operators.Repeat0Inf("*(%x20-21 / %x23-7E)", operators.Alts(
			"%x20-21 / %x23-7E",
			operators.Range("%x20-21", []byte{32}, []byte{33}),
			operators.Range("%x23-7E", []byte{35}, []byte{126}),
		)),
		core.DQUOTE(),
	)(s)
}

// comment = ";" *(WSP / VCHAR) CRLF
func Comment(s []byte) operators.Alternatives {
	return operators.Concat(
		"comment",
		operators.String(";", ";"),
		operators.Repeat0Inf("*(WSP / VCHAR)", operators.Alts(
			"WSP / VCHAR",
			core.WSP(),
			core.VCHAR(),
		)),
		core.CRLF(),
	)(s)
}

// concatenation = repetition *(1*c-wsp repetition)
func Concatenation(s []byte) operators.Alternatives {
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

// dec-val = "d" 1*DIGIT [ 1*("." 1*DIGIT) / ("-" 1*DIGIT) ]
func DecVal(s []byte) operators.Alternatives {
	return operators.Concat(
		"dec-val",
		operators.String("d", "d"),
		operators.Repeat1Inf("1*DIGIT", core.DIGIT()),
		operators.Optional("[ 1*(\".\" 1*DIGIT) / (\"-\" 1*DIGIT) ]", operators.Alts(
			"1*(\".\" 1*DIGIT) / (\"-\" 1*DIGIT)",
			operators.Repeat1Inf("1*(\".\" 1*DIGIT)", operators.Concat(
				"\".\" 1*DIGIT",
				operators.String(".", "."),
				operators.Repeat1Inf("1*DIGIT", core.DIGIT()),
			)),
			operators.Concat(
				"\"-\" 1*DIGIT",
				operators.String("-", "-"),
				operators.Repeat1Inf("1*DIGIT", core.DIGIT()),
			),
		)),
	)(s)
}

// defined-as = *c-wsp ("=" / "=/") *c-wsp
func DefinedAs(s []byte) operators.Alternatives {
	return operators.Concat(
		"defined-as",
		operators.Repeat0Inf("*c-wsp", CWsp),
		operators.Alts(
			"\"=\" / \"=/\"",
			operators.String("=", "="),
			operators.String("=/", "=/"),
		),
		operators.Repeat0Inf("*c-wsp", CWsp),
	)(s)
}

// element = rulename / group / option / char-val / num-val / prose-val
func Element(s []byte) operators.Alternatives {
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

// elements = alternation *WSP
func Elements(s []byte) operators.Alternatives {
	return operators.Concat(
		"elements",
		Alternation,
		operators.Repeat0Inf("*WSP", core.WSP()),
	)(s)
}

// group = "(" *c-wsp alternation *c-wsp ")"
func Group(s []byte) operators.Alternatives {
	return operators.Concat(
		"group",
		operators.String("(", "("),
		operators.Repeat0Inf("*c-wsp", CWsp),
		Alternation,
		operators.Repeat0Inf("*c-wsp", CWsp),
		operators.String(")", ")"),
	)(s)
}

// hex-val = "x" 1*HEXDIG [ 1*("." 1*HEXDIG) / ("-" 1*HEXDIG) ]
func HexVal(s []byte) operators.Alternatives {
	return operators.Concat(
		"hex-val",
		operators.String("x", "x"),
		operators.Repeat1Inf("1*HEXDIG", core.HEXDIG()),
		operators.Optional("[ 1*(\".\" 1*HEXDIG) / (\"-\" 1*HEXDIG) ]", operators.Alts(
			"1*(\".\" 1*HEXDIG) / (\"-\" 1*HEXDIG)",
			operators.Repeat1Inf("1*(\".\" 1*HEXDIG)", operators.Concat(
				"\".\" 1*HEXDIG",
				operators.String(".", "."),
				operators.Repeat1Inf("1*HEXDIG", core.HEXDIG()),
			)),
			operators.Concat(
				"\"-\" 1*HEXDIG",
				operators.String("-", "-"),
				operators.Repeat1Inf("1*HEXDIG", core.HEXDIG()),
			),
		)),
	)(s)
}

// num-val = "%" (bin-val / dec-val / hex-val)
func NumVal(s []byte) operators.Alternatives {
	return operators.Concat(
		"num-val",
		operators.String("%", "%"),
		operators.Alts(
			"bin-val / dec-val / hex-val",
			BinVal,
			DecVal,
			HexVal,
		),
	)(s)
}

// option = "[" *c-wsp alternation *c-wsp "]"
func Option(s []byte) operators.Alternatives {
	return operators.Concat(
		"option",
		operators.String("[", "["),
		operators.Repeat0Inf("*c-wsp", CWsp),
		Alternation,
		operators.Repeat0Inf("*c-wsp", CWsp),
		operators.String("]", "]"),
	)(s)
}

// prose-val = "<" *(%x20-3D / %x3F-7E) ">"
func ProseVal(s []byte) operators.Alternatives {
	return operators.Concat(
		"prose-val",
		operators.String("<", "<"),
		operators.Repeat0Inf("*(%x20-3D / %x3F-7E)", operators.Alts(
			"%x20-3D / %x3F-7E",
			operators.Range("%x20-3D", []byte{32}, []byte{61}),
			operators.Range("%x3F-7E", []byte{63}, []byte{126}),
		)),
		operators.String(">", ">"),
	)(s)
}

// repeat = 1*DIGIT / (*DIGIT "*" *DIGIT)
func Repeat(s []byte) operators.Alternatives {
	return operators.Alts(
		"repeat",
		operators.Repeat1Inf("1*DIGIT", core.DIGIT()),
		operators.Concat(
			"*DIGIT \"*\" *DIGIT",
			operators.Repeat0Inf("*DIGIT", core.DIGIT()),
			operators.String("*", "*"),
			operators.Repeat0Inf("*DIGIT", core.DIGIT()),
		),
	)(s)
}

// repetition = [repeat] element
func Repetition(s []byte) operators.Alternatives {
	return operators.Concat(
		"repetition",
		operators.Optional("[repeat]", Repeat),
		Element,
	)(s)
}

// rule = rulename defined-as elements c-nl
func Rule(s []byte) operators.Alternatives {
	return operators.Concat(
		"rule",
		Rulename,
		DefinedAs,
		Elements,
		CNl,
	)(s)
}

// rulelist = 1*( rule / (*WSP c-nl) )
func Rulelist(s []byte) operators.Alternatives {
	return operators.Repeat1Inf("rulelist", operators.Alts(
		"rule / (*WSP c-nl)",
		Rule,
		operators.Concat(
			"*WSP c-nl",
			operators.Repeat0Inf("*WSP", core.WSP()),
			CNl,
		),
	))(s)
}

// rulename = ALPHA *(ALPHA / DIGIT / "-")
func Rulename(s []byte) operators.Alternatives {
	return operators.Concat(
		"rulename",
		core.ALPHA(),
		operators.Repeat0Inf("*(ALPHA / DIGIT / \"-\")", operators.Alts(
			"ALPHA / DIGIT / \"-\"",
			core.ALPHA(),
			core.DIGIT(),
			operators.String("-", "-"),
		)),
	)(s)
}
