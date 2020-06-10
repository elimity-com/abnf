package abnf

import (
	"fmt"
	"github.com/elimity-com/abnf/definition"

	"github.com/dave/jennifer/jen"
)

const operatorsPkg = "github.com/elimity-com/abnf/operators"

var multiLineCall = jen.Options{
	Close:     ")",
	Multi:     true,
	Open:      "(",
	Separator: ",",
}

type generator struct {
	alts        bool
	packageName string
	rawABNF     string
}

// GenerateABNFAsOperators returns a *jen.File containing the given ABNF syntax as Go Operator functions.
func GenerateABNFAsOperators(packageName, rawABNF string) *jen.File {
	g := generator{
		packageName: packageName,
		rawABNF:     rawABNF,
	}
	return g.generate()
}

// GenerateABNFAsAlternatives returns a *jen.File containing the given ABNF syntax as Go functions that return Alternatives.
func GenerateABNFAsAlternatives(packageName, rawABNF string) *jen.File {
	g := generator{
		alts:        true,
		packageName: packageName,
		rawABNF:     rawABNF,
	}
	return g.generate()
}

func (g generator) generate() *jen.File {
	f := jen.NewFile(g.packageName)

	f.HeaderComment("This file is generated - do not edit.")
	f.Line()

	var returnParameter string
	if g.alts {
		returnParameter = "Alternatives"
	} else {
		returnParameter = "Operator"
	}

	alternatives := definition.Rulelist([]rune(g.rawABNF))
	for _, line := range alternatives.Best().Children {
		if line.Contains("rule") {
			f.Comment(fmt.Sprintf("%s", formatFuncComment(line.GetSubNode("rule").String())))

			name := line.GetSubNode("rulename").String()
			node := g.parseAlts(line.GetSubNode("alternation"))
			var returnValue jen.Code
			if g.alts {
				returnValue = jen.Return(node.toJen(key(name))).Call(jen.Id("s"))
			} else {
				returnValue = jen.Return(node.toJen(key(name)))
			}

			var params []jen.Code
			if g.alts {
				params = append(params, jen.Id("s").Op("[]").Id("rune"))
			}
			f.Func().Id(formatRuleName(name)).Call(params...).Qual(operatorsPkg, returnParameter).Block(
				returnValue,
			)
		}
	}
	return f
}

type key string

type generatorNode interface {
	toJen(k key) jen.Code
	getKey() key
}