package abnf

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"sort"
)

const operatorsPkg = "github.com/elimity-com/abnf/operators"

var multiLineCall = jen.Options{
	Close:     ")",
	Multi:     true,
	Open:      "(",
	Separator: ",",
}

type CodeGenerator struct {
	// PackageName of the generated file
	PackageName string
	// RawABNF syntax to parse
	RawABNF []byte
	// ExternalABNF reference to abnf syntax
	// e.g. ALPHA from github.com/elimity-com/abnf/core
	ExternalABNF map[string]ExternalABNF

	isOperator bool
	synonyms   map[string]string
}

type ExternalABNF struct {
	// IsOperator: operator / alternatives
	IsOperator bool
	// PackageName: e.g. github.com/elimity-com/abnf/core
	PackageName string
}

func (g *CodeGenerator) syn(key string) string {
	if syn, ok := g.synonyms[key]; ok {
		delete(g.synonyms, key) // consume
		return syn
	}
	return key
}

// GenerateABNFAsOperators returns a *jen.File containing the given ABNF syntax as Go Operator functions.
func (g *CodeGenerator) GenerateABNFAsOperators() *jen.File {
	g.isOperator = true
	return g.generate()
}

// GenerateABNFAsAlternatives returns a *jen.File containing the given ABNF syntax as Go functions that return Alternatives.
func (g *CodeGenerator) GenerateABNFAsAlternatives() *jen.File {
	return g.generate()
}

func (g *CodeGenerator) generate() *jen.File {
	g.synonyms = make(map[string]string) // synonyms

	f := jen.NewFile(g.PackageName)

	f.HeaderComment("This file is generated - do not edit.")
	f.Line()

	f.ImportName(operatorsPkg, "operators")

	returnParameter := func() (string, string) {
		if g.isOperator {
			return operatorsPkg, "Operator"
		}
		return operatorsPkg, "Alternatives"
	}

	returnValue := func(rule Rule) jen.Code {
		if g.isOperator {
			return jen.Return(rule.toJen(g))
		}
		return jen.Return(rule.toJen(g)).Call(jen.Id("s"))
	}

	ruleSet := NewRuleSet(g.RawABNF)

	keys := make([]string, 0)
	for k, _ := range ruleSet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		rule := ruleSet[k]

		f.Comment(fmt.Sprintf("%s = %s", rule.name, rule.operator.Key()))

		var params []jen.Code
		if !g.isOperator {
			params = []jen.Code{
				jen.Id("s").Op("[]").Id("byte"),
			}
		}

		f.Func().Id(formatRuleName(rule.name)).Call(params...).Qual(returnParameter()).Block(
			returnValue(rule),
		)

		f.Line()
	}

	return f
}

type codeGeneratorNode interface {
	toJen(g *CodeGenerator) jen.Code
}

func (r Rule) toJen(g *CodeGenerator) jen.Code {
	g.synonyms[r.operator.Key()] = r.name
	return r.operator.toJen(g)
}

func (alt AlternationOperator) toJen(g *CodeGenerator) jen.Code {
	return jen.Qual(operatorsPkg, "Alts").CustomFunc(multiLineCall, func(group *jen.Group) {
		group.Lit(g.syn(alt.key))
		for _, operator := range alt.subOperators {
			group.Add(operator.toJen(g))
		}
	})
}

func (concat ConcatenationOperator) toJen(g *CodeGenerator) jen.Code {
	return jen.Qual(operatorsPkg, "Concat").CustomFunc(multiLineCall, func(group *jen.Group) {
		group.Lit(g.syn(concat.key))
		for _, operator := range concat.subOperators {
			group.Add(operator.toJen(g))
		}
	})
}

func (rep RepetitionOperator) toJen(g *CodeGenerator) jen.Code {
	if rep.min == rep.max {
		return jen.Qual(operatorsPkg, "RepeatN").Call(
			jen.Lit(g.syn(rep.key)),
			jen.Lit(rep.min),
			rep.subOperator.toJen(g),
		)
	}

	if rep.max == -1 {
		switch rep.min {
		case 0:
			return jen.Qual(operatorsPkg, "Repeat0Inf").Call(
				jen.Lit(g.syn(rep.key)),
				rep.subOperator.toJen(g),
			)
		case 1:
			return jen.Qual(operatorsPkg, "Repeat1Inf").Call(
				jen.Lit(g.syn(rep.key)),
				rep.subOperator.toJen(g),
			)
		}
	}

	return jen.Qual(operatorsPkg, "Repeat").Call(
		jen.Lit(g.syn(rep.key)),
		jen.Lit(rep.min),
		jen.Lit(rep.max),
		rep.subOperator.toJen(g),
	)
}

func (name RuleNameOperator) toJen(g *CodeGenerator) jen.Code {
	if external, ok := g.ExternalABNF[name.key]; ok {
		if external.IsOperator {
			return jen.Qual(external.PackageName, name.key).Call()
		}
		return jen.Qual(external.PackageName, name.key)
	}
	if g.isOperator {
		return jen.Id(formatRuleName(g.syn(name.key))).Call()
	}
	return jen.Id(formatRuleName(g.syn(name.key)))
}

func (opt OptionOperator) toJen(g *CodeGenerator) jen.Code {
	return jen.Qual(operatorsPkg, "Optional").Call(
		jen.Lit(g.syn(opt.key)),
		opt.subOperator.toJen(g),
	)
}

func (value CharacterValueOperator) toJen(g *CodeGenerator) jen.Code {
	return jen.Qual(operatorsPkg, "String").Call(
		jen.Lit(g.syn(value.value)),
		jen.Lit(value.value),
	)
}

func (value NumericValueOperator) toJen(g *CodeGenerator) jen.Code {
	values := value.toIntegers()

	if value.hyphen {
		min, max := values[0], values[1]

		minValues := make([]jen.Code, len(min))
		for i, v := range min {
			minValues[i] = jen.Lit(v)
		}
		maxValues := make([]jen.Code, len(max))
		for i, v := range max {
			maxValues[i] = jen.Lit(v)
		}

		return jen.Qual(operatorsPkg, "Range").Call(
			jen.Lit(g.syn(value.key)),
			jen.Index().Byte().Values(minValues...),
			jen.Index().Byte().Values(maxValues...),
		)
	}

	if value.points {
		var str string
		for _, part := range values {
			var bytes []byte
			for _, i := range part {
				bytes = append(bytes, byte(i))
			}
			str += string(bytes)
		}
		return jen.Qual(operatorsPkg, "String").Call(
			jen.Lit(g.syn(value.key)),
			jen.Lit(str),
		)
	}

	bytes := make([]jen.Code, len(values[0]))
	for i, v := range values[0] {
		bytes[i] = jen.Lit(v)
	}
	return jen.Qual(operatorsPkg, "Terminal").Call(
		jen.Lit(g.syn(value.key)),
		jen.Index().Byte().Values(bytes...),
	)
}
