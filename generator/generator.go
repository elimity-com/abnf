package generator

import (
	"encoding/hex"
	"fmt"
	"github.com/elimity-com/abnf"
	"github.com/elimity-com/abnf/operators"
	"strconv"

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

func GenerateABNFAsOperators(packageName, rawABNF string) *jen.File {
	g := generator{
		packageName: packageName,
		rawABNF:     rawABNF,
	}
	return g.Generate()
}

func GenerateABNFAsAlternatives(packageName, rawABNF string) *jen.File {
	g := generator{
		alts:        true,
		packageName: packageName,
		rawABNF:     rawABNF,
	}
	return g.Generate()
}

func (g generator) Generate() *jen.File {
	f := jen.NewFile(g.packageName)

	f.HeaderComment("This file is generated - do not edit.")
	f.Line()

	var returnParameter string
	if g.alts {
		returnParameter = "Alternatives"
	} else {
		returnParameter = "Operator"
	}

	alternatives := abnf.Rulelist([]rune(g.rawABNF))
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

type alts struct {
	key      string
	elements []generatorNode
}

func (a alts) toJen(k key) jen.Code {
	if len(a.elements) == 1 {
		return a.elements[0].toJen(k)
	}

	return jen.Qual(operatorsPkg, "Alts").CustomFunc(multiLineCall, func(g *jen.Group) {
		g.Lit(string(k))
		for _, e := range a.elements {
			g.Add(e.toJen(e.getKey()))
		}
	})
}

func (a alts) getKey() key {
	return key(a.key)
}

func (g generator) parseAlts(node *operators.Node) alts {
	elements := []generatorNode{g.parseConcat(node.GetSubNode("concatenation"))}
	for _, child := range node.GetSubNodesBefore("*c-wsp \"/\" *c-wsp concatenation", "\"(\"", "\"[\"") {
		if c := child.GetNode("concatenation"); c != nil {
			elements = append(elements, g.parseConcat(c))
		}
	}
	return alts{
		key:      node.String(),
		elements: elements,
	}
}

type concat struct {
	key      string
	elements []generatorNode
}

func (c concat) toJen(k key) jen.Code {
	if len(c.elements) == 1 {
		return c.elements[0].toJen(k)
	}

	return jen.Qual(operatorsPkg, "Concat").CustomFunc(multiLineCall, func(g *jen.Group) {
		g.Lit(string(k))
		for _, e := range c.elements {
			g.Add(e.toJen(e.getKey()))
		}
	})
}

func (c concat) getKey() key {
	return key(c.key)
}

func (g generator) parseConcat(node *operators.Node) concat {
	elements := []generatorNode{g.parseRep(node.GetSubNode("repetition"))}
	for _, child := range node.GetSubNodesBefore("1*c-wsp repetition", "\"(\"", "\"[\"") {
		if c := child.GetNode("repetition"); c != nil {
			elements = append(elements, g.parseRep(c))
		}
	}
	return concat{
		key:      node.String(),
		elements: elements,
	}
}

type rep struct {
	key      string
	min, max int
	repeat   bool
	element  generatorNode
}

func (r rep) toJen(k key) jen.Code {
	if !r.repeat {
		return r.element.toJen(k)
	}

	if r.min == r.max {
		return jen.Qual(operatorsPkg, "RepeatN").Call(
			jen.Lit(string(k)),
			jen.Lit(r.min),
			r.element.toJen(r.element.getKey()),
		)
	}

	if r.max == -1 {
		switch r.min {
		case 0:
			return jen.Qual(operatorsPkg, "Repeat0Inf").Call(
				jen.Lit(string(k)),
				r.element.toJen(r.element.getKey()),
			)
		case 1:
			return jen.Qual(operatorsPkg, "Repeat1Inf").Call(
				jen.Lit(string(k)),
				r.element.toJen(r.element.getKey()),
			)
		}
	}

	return jen.Qual(operatorsPkg, "Repeat").Call(
		jen.Lit(string(k)),
		jen.Lit(r.min),
		jen.Lit(r.max),
		r.element.toJen(r.getKey()),
	)
}

func (r rep) getKey() key {
	return key(r.key)
}

func (g generator) parseRep(node *operators.Node) rep {
	if node.Children[0].IsEmpty() {
		return rep{
			key:     node.String(),
			repeat:  false,
			element: g.parseElement(node.GetSubNode("element")),
		}
	}

	var min, max int
	repeat := node.GetSubNode("repeat")
	if repeat.Key == "1*DIGIT" {
		i, _ := strconv.Atoi(repeat.String())
		min = i
		max = i
	} else {
		max = -1
		var astrix bool
		for _, child := range repeat.Children[0].Children {
			if child.Key == "*DIGIT" {
				if child.IsEmpty() {
					continue
				}

				if !astrix {
					min, _ = strconv.Atoi(child.String())
				} else {
					max, _ = strconv.Atoi(child.String())
				}
			} else {
				astrix = true
			}
		}
	}

	return rep{
		key:     node.String(),
		min:     min,
		max:     max,
		repeat:  true,
		element: g.parseElement(node.GetSubNode("element")),
	}
}

func (g generator) parseElement(node *operators.Node) generatorNode {
	switch child := node.Children[0]; child.Key {
	case "rulename":
		return identifier{
			call:  !g.alts,
			value: formatRuleName(child.String()),
		}
	case "group":
		return g.parseAlts(child.GetSubNode("alternation"))
	case "option":
		value := child.GetSubNode("alternation")
		return optionalValue{
			key:     value.String(),
			element: g.parseAlts(value),
		}
	case "char-val":
		values := child.GetSubNodes("%x20-21 / %x23-7E")
		if len(values) == 1 {
			value := values[0]
			return runeValue{
				key:   value.String(),
				value: int(value.String()[0]),
			}
		} else {
			value := child.GetSubNode("*(%x20-21 / %x23-7E)")
			return stringValue{
				key:   value.String(),
				value: value.String(),
			}
		}
	case "num-val":
		return g.parseNumVal(child)
	case "prose-val":
		fmt.Println(child)
	default:
		fmt.Println(child.Key, child)
	}
	return nil
}

func (g generator) parseNumVal(node *operators.Node) generatorNode {
	switch child := node.Children[1].Children[0]; child.Key {
	default:
	case "bin-val":
		fmt.Println(child)
	case "dec-val":
		fmt.Println(child)
	case "hex-val":
		var first int
		var hyphen bool
		for _, v := range child.Children {
			if c := v.GetNode("1*HEXDIG"); c != nil {
				// TODO: "."

				if v.Contains("\"-\"") {
					hyphen = true
				}
				raw, _ := hex.DecodeString(c.String())
				if !hyphen {
					first = int(raw[0])
				} else {
					return rangeValue{
						key: child.String(),
						min: first,
						max: int(raw[0]),
					}
				}
			}
		}
		return runeValue{
			key:   child.String(),
			value: first,
		}
	}
	return nil
}

type optionalValue struct {
	key     string
	element generatorNode
}

func (v optionalValue) toJen(k key) jen.Code {
	return jen.Qual(operatorsPkg, "Optional").Call(
		jen.Lit(string(k)),
		v.element.toJen(v.getKey()),
	)
}

func (v optionalValue) getKey() key {
	return v.element.getKey()
}

type runeValue struct {
	key   string
	value int
}

func (v runeValue) toJen(k key) jen.Code {
	return jen.Qual(operatorsPkg, "Rune").Call(
		jen.Lit(string(k)),
		jen.Lit(v.value),
	)
}

func (v runeValue) getKey() key {
	return key(v.key)
}

type stringValue struct {
	key   string
	value string
}

func (v stringValue) toJen(k key) jen.Code {
	return jen.Qual(operatorsPkg, "String").Call(
		jen.Lit(string(k)),
		jen.Lit(v.value),
	)
}

func (v stringValue) getKey() key {
	return key(v.key)
}

type rangeValue struct {
	key      string
	min, max int
}

func (v rangeValue) toJen(k key) jen.Code {
	return jen.Qual(operatorsPkg, "Range").Call(
		jen.Lit(string(k)),
		jen.Lit(v.min),
		jen.Lit(v.max),
	)
}

func (v rangeValue) getKey() key {
	return key(v.key)
}

type identifier struct {
	call  bool
	value string
}

func (i identifier) toJen(_ key) jen.Code {
	if i.call {
		return jen.Id(i.value).Call()
	}
	return jen.Id(i.value)
}

func (i identifier) getKey() key {
	return key(i.value)
}
