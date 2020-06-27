package abnf

import (
	"github.com/elimity-com/abnf/encoding"
	"github.com/elimity-com/abnf/operators"
	"sync"
)

type ParserGenerator struct {
	// RawABNF syntax to parse
	RawABNF []byte
	// ExternalABNF reference to abnf syntax
	// e.g. ALPHA from github.com/elimity-com/abnf/core
	ExternalABNF map[string]operators.Operator

	sync.WaitGroup
	internalABNFMutex sync.RWMutex
	internalABNF      map[string]operators.Operator
}

func (g *ParserGenerator) GenerateABNFAsOperators() map[string]operators.Operator {
	ruleSet := NewRuleSet(g.RawABNF)

	g.internalABNFMutex = sync.RWMutex{}
	g.internalABNF = make(map[string]operators.Operator)
	for name, _ := range ruleSet {
		g.internalABNFMutex.Lock()
		g.internalABNF[name] = nil
		g.internalABNFMutex.Unlock()
	}

	for name, rule := range ruleSet {
		g.Add(1)
		go func(name string, rule Rule) {
			function := rule.toFunc(g)
			g.internalABNFMutex.Lock()
			g.internalABNF[name] = function
			g.internalABNFMutex.Unlock()
			defer g.Done()
		}(name, rule)
	}
	g.Wait()
	return g.internalABNF
}

type parserGeneratorNode interface {
	toFunc(g *ParserGenerator) operators.Operator
}

func (r Rule) toFunc(g *ParserGenerator) operators.Operator {
	return r.operator.toFunc(g)
}

func (alt AlternationOperator) toFunc(g *ParserGenerator) operators.Operator {
	var rules []operators.Operator
	for _, subOperator := range alt.subOperators {
		rules = append(rules, subOperator.toFunc(g))
	}
	return operators.Alts(alt.key, rules...)
}

func (concat ConcatenationOperator) toFunc(g *ParserGenerator) operators.Operator {
	var rules []operators.Operator
	for _, subOperator := range concat.subOperators {
		rules = append(rules, subOperator.toFunc(g))
	}
	return operators.Concat(concat.key, rules...)
}

func (rep RepetitionOperator) toFunc(g *ParserGenerator) operators.Operator {
	if rep.min == rep.max {
		return operators.RepeatN(rep.key, rep.min, rep.subOperator.toFunc(g))
	}

	if rep.max == -1 {
		switch rep.min {
		case 0:
			return operators.Repeat0Inf(rep.key, rep.subOperator.toFunc(g))
		case 1:
			return operators.Repeat1Inf(rep.key, rep.subOperator.toFunc(g))
		}
	}

	return operators.Repeat(rep.key, rep.min, rep.max, rep.subOperator.toFunc(g))
}

func (name RuleNameOperator) toFunc(g *ParserGenerator) operators.Operator {
	if external, ok := g.ExternalABNF[name.key]; ok {
		return external
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			g.internalABNFMutex.RLock()
			v := g.internalABNF[name.key]
			g.internalABNFMutex.RUnlock()
			if v != nil {
				break
			}
		}
	}()
	wg.Wait()
	return g.internalABNF[name.key]
}

func (opt OptionOperator) toFunc(g *ParserGenerator) operators.Operator {
	return operators.Optional(opt.key, opt.subOperator.toFunc(g))
}

func (value CharacterValueOperator) toFunc(_ *ParserGenerator) operators.Operator {
	return operators.String(value.value, value.value)
}

func (value NumericValueOperator) toFunc(g *ParserGenerator) operators.Operator {
	values := value.toIntegers()

	if value.hyphen {
		min, max := values[0], values[1]

		minValues := make([]byte, len(min))
		for i, v := range min {
			minValues[i] = byte(v)
		}
		maxValues := make([]byte, len(max))
		for i, v := range max {
			maxValues[i] = byte(v)
		}

		return operators.Range(value.key, minValues, maxValues)
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
		return operators.String(value.key, str)
	}

	bytes := make([]byte, len(values[0]))
	for i, v := range values[0] {
		bytes[i] = byte(v)
	}
	bytes, _ = encoding.ASCII.NewEncoder().Bytes(bytes)
	return operators.Terminal(value.key, bytes)
}
