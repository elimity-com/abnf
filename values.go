package abnf

import (
	"encoding/hex"
	"github.com/elimity-com/abnf/operators"
	"log"
	"strconv"
)

func (g *Generator) parseNumVal(node *operators.Node) generatorNode {
	child := node.Children[1].Children[0]
	var numValue string
	switch child.Key {
	case "bin-val":
		numValue = "1*BIT"
	case "dec-val":
		numValue = "1*DIGIT"
	case "hex-val":
		numValue = "1*HEXDIG"
	}

	var first string
	second := make([]string, 0)
	var hyphen, point bool
	for _, v := range child.Children {
		if c := v.GetNode(numValue); c != nil {
			if v.Contains("\"-\"") {
				hyphen = true
			}
			if v.Contains("\".\"") {
				point = true
			}

			if hyphen {
				second = []string{c.String()}
			} else if point {
				second = []string{first}
				for _, s := range v.GetSubNodes(numValue) {
					second = append(second, s.String())
				}
			} else {
				first = c.String()
			}
		}
	}

	var min, max []int
	var values string
	switch child.Key {
	case "bin-val":
		raw, _ := strconv.ParseInt(first, 2, 64)
		min = []int{int(raw)}
		if hyphen {
			raw, _ := strconv.ParseInt(second[0], 2, 64)
			max = []int{int(raw)}
		}
		if point {
			for _, v := range second {
				raw, _ := strconv.ParseInt(v, 2, 64)
				values += string(raw)
			}
		}
	case "dec-val":
		d, _ := strconv.Atoi(first)
		min = []int{d}
		if hyphen {
			d, _ = strconv.Atoi(second[0])
			max = []int{d}
		}
		if point {
			for _, v := range second {
				raw, _ := strconv.Atoi(v)
				values += string(raw)
			}
		}
	case "hex-val":
		min = g.hexStringToBytes(first)
		if hyphen {
			max = g.hexStringToBytes(second[0])
		}
		if point {
			for _, v := range second {
				raw, _ := hex.DecodeString(v)
				values += string(raw[0])
			}
		}
	}

	if hyphen {
		return rangeValue{
			key: child.String(),
			min: min,
			max: max,
		}
	}

	if point {
		return stringValue{
			key:   child.String(),
			value: values,
		}
	}

	return terminalValue{
		key:    child.String(),
		values: min,
	}
}

func (g *Generator) parseProseVal(node *operators.Node) generatorNode {
	log.Fatal("not implemented, to be used as last resort")
	return nil
}

func (g *Generator) hexStringToBytes(hexStr string) []int {
	n, _ := strconv.ParseInt(hexStr, 16, 64)
	b := make([]int, (len(hexStr)+1)/2)
	for i := range b {
		b[i] = int(byte(n >> uint64(8*(len(b)-i-1))))
	}
	return b
}
