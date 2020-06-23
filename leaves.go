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

type terminalValue struct {
	key   string
	values []int
}

func (v terminalValue) toJen(k key) jen.Code {
	values := make([]jen.Code, len(v.values))
	for i, v := range v.values {
		values[i] = jen.Lit(v)
	}

	return jen.Qual(operatorsPkg, "Terminal").Call(
		jen.Lit(string(k)),
		jen.Index().Byte().Values(values...),
	)
}

func (v terminalValue) getKey() key {
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
	min, max []int
}

func (v rangeValue) toJen(k key) jen.Code {
	minValues := make([]jen.Code, len(v.min))
	for i, v := range v.min {
		minValues[i] = jen.Lit(v)
	}
	maxValues := make([]jen.Code, len(v.max))
	for i, v := range v.max {
		maxValues[i] = jen.Lit(v)
	}

	return jen.Qual(operatorsPkg, "Range").Call(
		jen.Lit(string(k)),
		jen.Index().Byte().Values(minValues...),
		jen.Index().Byte().Values(maxValues...),
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

type externalIdentifier struct {
	call       bool
	pkg, value string
}

func (i externalIdentifier) toJen(_ key) jen.Code {
	if i.call {
		return jen.Qual(i.pkg, i.value).Call()
	}
	return jen.Qual(i.pkg, i.value)
}

func (i externalIdentifier) getKey() key {
	return key(i.value)
}
