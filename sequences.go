package abnf

import (
	"github.com/elimity-com/abnf/operators"

	"github.com/dave/jennifer/jen"
)

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
