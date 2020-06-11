package abnf

import (
	"strconv"

	"github.com/elimity-com/abnf/operators"

	"github.com/dave/jennifer/jen"
)

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

func (g *Generator) parseRep(node *operators.Node) rep {
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
