package abnf

import (
	"fmt"

	"github.com/elimity-com/abnf/operators"
)

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
	}
	return nil
}
