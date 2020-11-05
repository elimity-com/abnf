package abnf

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

const operatorsPkg = "github.com/elimity-com/abnf/operators"

type CodeGenerator struct {
	writer io.Writer
	last   rune
	prefix string

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

func (g *CodeGenerator) c(format string, args ...interface{}) error {
	if err := g.w("// "); err != nil {
		return err
	}
	return g.wlnf(format, args...)
}

func (g *CodeGenerator) w(p string) error {
	if g.last == '\n' && p != "\n" {
		p = g.prefix + p
	}
	g.last = rune(p[len(p)-1])
	_, err := g.writer.Write([]byte(p))
	return err
}

func (g *CodeGenerator) wf(format string, args ...interface{}) error {
	return g.w(fmt.Sprintf(format, args...))
}

func (g *CodeGenerator) wln(p string) error {
	return g.w(p + "\n")
}

func (g *CodeGenerator) wlnf(format string, args ...interface{}) error {
	return g.wln(fmt.Sprintf(format, args...))
}

func (g *CodeGenerator) ln() error {
	return g.w("\n")
}

func (g *CodeGenerator) in(f func()) {
	tmp := g.prefix
	g.prefix += "\t"
	defer func() {
		g.prefix = tmp
	}()

	f()
}

type ExternalABNF struct {
	// IsOperator: operator / alternatives
	IsOperator bool
	// PackagePath: e.g. github.com/elimity-com/abnf/core
	PackagePath string
	// PackageName: e.g. core
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
func (g *CodeGenerator) GenerateABNFAsOperators(w io.Writer) {
	g.writer = w
	g.isOperator = true
	g.generate()
}

// GenerateABNFAsAlternatives returns a *jen.File containing the given ABNF syntax as Go functions that return Alternatives.
func (g *CodeGenerator) GenerateABNFAsAlternatives(w io.Writer) {
	g.writer = w
	g.generate()
}

func (g *CodeGenerator) generate() {
	g.synonyms = make(map[string]string) // synonyms

	g.c("This file is generated - do not edit.")
	g.ln()
	g.wlnf("package %s", g.PackageName)
	g.ln()
	g.w("import ")
	if len(g.ExternalABNF) != 0 {
		imports := make(map[string]struct{})
		for _, i := range g.ExternalABNF {
			imports[i.PackagePath] = struct{}{}
		}
		g.wln("(")
		g.in(func() {
			for i := range imports {
				g.wlnf("%q", i)
			}
			g.ln()
			g.wlnf("%q", operatorsPkg)
		})
		g.wln(")")
	} else {
		g.wlnf("%q", operatorsPkg)
	}

	ruleSet := NewRuleSet(g.RawABNF)

	keys := make([]string, 0)
	for k, _ := range ruleSet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		rule := ruleSet[k]

		g.ln()
		g.c("%s = %s", rule.name, rule.operator.Key())
		g.wf("func %s(", formatRuleName(rule.name))
		if g.isOperator {
			g.wln(") operators.Operator {")
		} else {
			g.wln("s []byte) operators.Alternatives {")
		}
		g.in(func() {
			g.w("return ")
			rule.generate(g)
			if !g.isOperator {
				g.w("(s)")
			}
			g.ln()
		})
		g.wln("}")
	}
}

type codeGeneratorNode interface {
	generate(g *CodeGenerator)
}

func (r Rule) generate(g *CodeGenerator) {
	g.synonyms[r.operator.Key()] = r.name
	r.operator.generate(g)
}

func (alt AlternationOperator) generate(g *CodeGenerator) {
	g.wln("operators.Alts(")
	g.in(func() {
		g.wlnf("%q,", g.syn(alt.key))
		for _, operator := range alt.subOperators {
			operator.generate(g)
			g.wln(",")
		}
	})
	g.w(")")
}

func (concat ConcatenationOperator) generate(g *CodeGenerator) {
	g.wln("operators.Concat(")
	g.in(func() {
		g.wlnf("%q,", g.syn(concat.key))
		for _, operator := range concat.subOperators {
			operator.generate(g)
			g.wln(",")
		}
	})
	g.w(")")
}

func (rep RepetitionOperator) generate(g *CodeGenerator) {
	if rep.min == rep.max {
		g.wf("operators.RepeatN(%q, %q, ", g.syn(rep.key), rep.min)
	} else if rep.max == -1 {
		switch rep.min {
		case 0:
			g.wf("operators.Repeat0Inf(%q, ", g.syn(rep.key))
		case 1:
			g.wf("operators.Repeat1Inf(%q, ", g.syn(rep.key))
		}
	} else {
		g.wf("operators.Repeat1Inf(%q, %q, %q, ", g.syn(rep.key), rep.min, rep.max)
	}
	rep.subOperator.generate(g)
	g.w(")")
}

func (name RuleNameOperator) generate(g *CodeGenerator) {
	if external, ok := g.ExternalABNF[name.key]; ok {
		g.wf("%s.%s", external.PackageName, name.key)
		if external.IsOperator {
			g.w("()")
		}
	} else {
		g.w(formatRuleName(g.syn(name.key)))
		if g.isOperator {
			g.w("()")
		}
	}
}

func (opt OptionOperator) generate(g *CodeGenerator) {
	g.wf("operators.Optional(%q, ", g.syn(opt.key))
	opt.subOperator.generate(g)
	g.w(")")
}

func (value CharacterValueOperator) generate(g *CodeGenerator) {
	g.wf("operators.String(%q, %q)", g.syn(value.value), value.value)
}

func (value NumericValueOperator) generate(g *CodeGenerator) {
	values := value.toIntegers()

	if value.hyphen {
		min, max := values[0], values[1]
		var minValues string
		for _, v := range min {
			minValues = strconv.Itoa(v)
		}
		var maxValues string
		for _, v := range max {
			maxValues = strconv.Itoa(v)
		}
		g.wf("operators.Range(%q, []byte{%s}, []byte{%s})", g.syn(value.key), minValues, maxValues)
	} else if value.points {
		var str string
		for _, part := range values {
			var bytes []byte
			for _, i := range part {
				bytes = append(bytes, byte(i))
			}
			str += string(bytes)
		}
		g.wf("operators.String(%q, %q)", g.syn(value.key), str)
	} else {
		var bytes string
		for _, v := range values[0] {
			bytes = strconv.Itoa(v)
		}
		g.wf("operators.Terminal(%q, []byte{%s})", g.syn(value.key), bytes)
	}
}
