package generator

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/elimity-com/abnf"
	"github.com/elimity-com/abnf/operators"

	"github.com/dave/jennifer/jen"
)

const operatorsPkg = "github.com/elimity-com/abnf/operators"

var space = regexp.MustCompile(`\s+`)

var multiLineCall = jen.Options{
	Close:     ")",
	Multi:     true,
	Open:      "(",
	Separator: ",",
}

func GenerateABNF(packageName, rawABNF string) *jen.File {
	f := jen.NewFile(packageName)

	f.HeaderComment("This file is generated - do not edit.")
	f.Line()

	alternatives := abnf.RuleList([]rune(rawABNF))
	for _, line := range alternatives.Best().Children {
		if line.Contains("rule") {
			comment := space.ReplaceAllString(strings.Split(line.GetSubNode("rule").String(), ";")[0], " ")
			f.Comment(fmt.Sprintf("%s", comment))

			name := line.GetSubNode("rulename").String()
			node := parseAlts(line.GetSubNode("alternation"))
			f.Func().Id(name).Call().Qual(operatorsPkg, "Operator").Block(
				jen.Return(node.toJen(key(name))),
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

func parseAlts(node *operators.Node) alts {
	elements := []generatorNode{parseConcat(node.GetSubNode("concatenation"))}
	for _, child := range node.GetSubNodesBefore("*c-wsp \"/\" *c-wsp concatenation", "(") {
		if c := child.GetNode("concatenation"); c != nil {
			elements = append(elements, parseConcat(c))
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

func parseConcat(node *operators.Node) concat {
	elements := []generatorNode{parseRep(node.GetSubNode("repetition"))}
	for _, child := range node.GetSubNodesBefore("1*c-wsp repetition", "(") {
		if c := child.GetNode("repetition"); c != nil {
			elements = append(elements, parseRep(c))
		}
	}
	return concat{
		key:      node.String(),
		elements: elements,
	}
}

type rep struct {
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
			r.element.toJen(r.getKey()),
		)
	}

	if r.max == -1 {
		switch r.min {
		case 0:
			return jen.Qual(operatorsPkg, "Repeat0Inf").Call(
				jen.Lit(string(k)),
				r.element.toJen(r.getKey()),
			)
		case 1:
			return jen.Qual(operatorsPkg, "Repeat1Inf").Call(
				jen.Lit(string(k)),
				r.element.toJen(r.getKey()),
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
	return r.element.getKey()
}

func parseRep(node *operators.Node) rep {
	if node.Children[0].IsEmpty() {
		return rep{
			repeat:  false,
			element: parseElement(node.GetSubNode("element")),
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
		min:     min,
		max:     max,
		repeat:  true,
		element: parseElement(node.GetSubNode("element")),
	}
}

func parseElement(node *operators.Node) generatorNode {
	switch child := node.Children[0]; child.Key {
	case "rulename":
		return identifier{
			value: child.String(),
		}
	case "group":
		return parseAlts(child.GetSubNode("alternation"))
	case "option":
		fmt.Println("option")
	case "char-val":
		values := child.GetSubNodes("%x20-21 / %x23-7E")
		if len(values) == 1 {
			value := values[0]
			return runeValue{
				key:   value.String(),
				value: int(value.String()[0]),
			}
		} else {
			fmt.Println("oops")
		}
	case "num-val":
		return parseNumVal(child)
	case "prose-val":
		fmt.Println(child)
	default:
		fmt.Println(child.Key, child)
	}
	return nil
}

func parseNumVal(node *operators.Node) generatorNode {
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

				if v.Contains("-") {
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
	value string
}

func (i identifier) toJen(_ key) jen.Code {
	return jen.Id(i.value).Call()
}

func (i identifier) getKey() key {
	return key(i.value)
}
