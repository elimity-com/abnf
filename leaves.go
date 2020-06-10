package abnf

import "github.com/dave/jennifer/jen"

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
