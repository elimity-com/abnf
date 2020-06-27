package abnf

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/elimity-com/abnf/definition"
	"github.com/elimity-com/abnf/operators"
)

// NewRuleSet converts given raw data to a set of ABNF rules.
func NewRuleSet(rawABNF []byte) RuleSet {
	rawRuleList := definition.Rulelist(rawABNF).Best()

	ruleSet := make(RuleSet)
	for _, line := range rawRuleList.Children {
		if line.Contains("rule") {
			rule := parseRule(line)
			ruleSet[rule.name] = rule
		}
	}
	return ruleSet
}

// RuleList is a list of ABNF rules.
type RuleList []Rule

// RuleSet converts a rule list to a set.
func (list RuleList) RuleSet() RuleSet {
	ruleSet := make(RuleSet)
	for _, rule := range list {
		ruleSet[rule.Name()] = rule
	}
	return ruleSet
}

// RuleSet is a set of ABNF rules.
type RuleSet map[string]Rule

// RuleList converts a rule set to a list.
func (set RuleSet) RuleList() RuleList {
	var ruleList RuleList
	for _, rule := range set {
		ruleList = append(ruleList, rule)
	}
	return ruleList
}

// Rule represents an ABNF rule.
type Rule struct {
	name     string
	operator Operator
}

// Equals checks whether both rule trees are equal to each other.
func (r Rule) Equals(other Rule) error {
	if r.name != other.name {
		return fmt.Errorf("names do not match: expected %s, got %s", r.name, other.name)
	}
	return r.operator.equals(other.operator)
}

// Name returns the name of the rule.
func (r Rule) Name() string {
	return r.name
}

// parseRule converts a raw rule node to a (more) readable one.
// ABNF: rule = rulename defined-as elements c-nl
func parseRule(rawNode *operators.Node) Rule {
	return Rule{
		name:     rawNode.GetSubNode("rulename").String(),
		operator: parseAlternation(rawNode.GetSubNode("alternation")),
	}
}

// Operator represents a node of a rule.
type Operator interface {
	// Key returns that key (name) of the operator.s
	Key() string
	// equals checks whether this operator equals given other.
	equals(other Operator) error

	codeGeneratorNode   // code generator
	parserGeneratorNode // parser generator
}

// AlternationOperator represents an alternation node of a rule.
type AlternationOperator struct {
	key          string
	subOperators []Operator
}

func (alt AlternationOperator) Key() string {
	return alt.key
}

func (alt AlternationOperator) equals(other Operator) error {
	otherAlt, ok := other.(AlternationOperator)
	if !ok {
		return fmt.Errorf("other is not of the same type: %s", reflect.TypeOf(other))
	}
	if alt.key != otherAlt.key {
		return fmt.Errorf("keys do not match: expected %s, got %s", alt.key, otherAlt.key)
	}
	if len(alt.subOperators) != len(otherAlt.subOperators) {
		return fmt.Errorf("lenght of sub operators do not match: expected %d, got %d", len(alt.subOperators), len(otherAlt.subOperators))
	}
	for i, subOperator := range alt.subOperators {
		if err := subOperator.equals(otherAlt.subOperators[i]); err != nil {
			return err
		}
	}
	return nil
}

// parseAlternation converts a raw (nested) alternation node to a (more) readable one.
// ABNF: alternation = concatenation *(*c-wsp "/" *c-wsp concatenation)
func parseAlternation(rawNode *operators.Node) Operator {
	// an alternation has at least one concatenation node
	subOperators := []Operator{
		parseConcatenation(rawNode.GetSubNode("concatenation")),
	}
	// get all other concatenation nodes
	for _, other := range rawNode.GetSubNodesBefore(`*c-wsp "/" *c-wsp concatenation`, "(", "[") {
		if rawConcat := other.GetNode("concatenation"); rawConcat != nil {
			subOperators = append(subOperators, parseConcatenation(rawConcat))
		}
	}
	// not need to return an alternation of one element
	if len(subOperators) == 1 {
		return subOperators[0]
	}
	return AlternationOperator{
		key:          rawNode.String(),
		subOperators: subOperators,
	}
}

// ConcatenationOperator represents a concatenation node of a rule.
type ConcatenationOperator struct {
	key          string
	subOperators []Operator
}

func (concat ConcatenationOperator) Key() string {
	return concat.key
}

func (concat ConcatenationOperator) equals(other Operator) error {
	otherConcat, ok := other.(ConcatenationOperator)
	if !ok {
		return fmt.Errorf("other is not of the same type: %s", reflect.TypeOf(other))
	}
	if concat.key != otherConcat.key {
		return fmt.Errorf("keys do not match: expected %s, got %s", concat.key, otherConcat.key)
	}
	if len(concat.subOperators) != len(otherConcat.subOperators) {
		return fmt.Errorf("lenght of sub operators do not match: expected %d, got %d", len(concat.subOperators), len(otherConcat.subOperators))
	}
	for i, subOperator := range concat.subOperators {
		if err := subOperator.equals(otherConcat.subOperators[i]); err != nil {
			return err
		}
	}
	return nil
}

// parseConcatenation converts a raw (nested) concatenation node to a (more) readable one.
// ABNF: concatenation = repetition *(1*c-wsp repetition)
func parseConcatenation(rawNode *operators.Node) Operator {
	// a concatenation has at least one repetition node
	subOperators := []Operator{
		parseRepetition(rawNode.GetSubNode("repetition")),
	}
	// get all other repetition nodes
	for _, other := range rawNode.GetSubNodesBefore(`1*c-wsp repetition`, "(", "[") {
		if rawConcat := other.GetNode("repetition"); rawConcat != nil {
			subOperators = append(subOperators, parseRepetition(rawConcat))
		}
	}
	// not need to return a concatenation of one element
	if len(subOperators) == 1 {
		return subOperators[0]
	}
	return ConcatenationOperator{
		key:          rawNode.String(),
		subOperators: subOperators,
	}
}

// RepetitionOperator represents a repetition node of a rule.
type RepetitionOperator struct {
	key         string
	min, max    int
	subOperator Operator
}

func (rep RepetitionOperator) Key() string {
	return rep.key
}

func (rep RepetitionOperator) equals(other Operator) error {
	otherRep, ok := other.(RepetitionOperator)
	if !ok {
		return fmt.Errorf("other is not of the same type: %s", reflect.TypeOf(other))
	}
	if rep.key != otherRep.key {
		return fmt.Errorf("keys do not match: expected %s, got %s", rep.key, otherRep.key)
	}
	if rep.min != otherRep.min {
		return fmt.Errorf("min subValues do not match: expected %d, got %d", rep.min, otherRep.min)
	}
	if rep.max != otherRep.max {
		return fmt.Errorf("max subValues do not match: expected %d, got %d", rep.max, otherRep.max)
	}
	return rep.subOperator.equals(otherRep.subOperator)
}

// parseRepetition converts a raw (nested) repetition node to a (more) readable one.
// ABNF: repetition = [repeat] element
func parseRepetition(rawNode *operators.Node) Operator {
	if rawNode.Children[0].IsEmpty() {
		// no repeat
		return parseElement(rawNode.GetSubNode("element"))
	}
	min, max := parseRepeat(rawNode.GetSubNode("repeat"))
	return RepetitionOperator{
		key: rawNode.String(),
		min: min, max: max,
		subOperator: parseElement(rawNode.GetSubNode("element")),
	}
}

// parseRepetition converts a raw (nested) repetition node to a two their respective min and max values.
// ABNF: repeat = 1*DIGIT / (*DIGIT "*" *DIGIT)
func parseRepeat(rawNode *operators.Node) (int, int) {
	if rawNode.Key == "1*DIGIT" {
		i, _ := strconv.Atoi(rawNode.String())
		return i, i
	}
	min, max, asterisk := 0, -1, false
	for _, child := range rawNode.Children[0].Children {
		if child.Key == "*DIGIT" {
			if child.IsEmpty() {
				continue
			}
			if !asterisk {
				min, _ = strconv.Atoi(child.String())
			} else {
				max, _ = strconv.Atoi(child.String())
			}
		} else {
			asterisk = true
		}
	}
	return min, max
}

// parseRepetition converts a raw (nested) element node to a (more) readable one.
// ABNF: element =  rulename / group / option / char-val / num-val / prose-val
func parseElement(rawNode *operators.Node) Operator {
	switch rawNode := rawNode.Children[0]; rawNode.Key {
	case "rulename":
		return parseRuleName(rawNode)
	case "group":
		return parseGroup(rawNode)
	case "option":
		return parseOption(rawNode)
	case "char-val":
		return parseCharacterValue(rawNode)
	case "num-val":
		return parseNumericValue(rawNode)
	case "prose-val":
		panic("not implemented")
	default:
		return nil
	}
}

// RuleNameOperator represents a rule name node of a rule.
type RuleNameOperator struct {
	key string
}

func (name RuleNameOperator) Key() string {
	return name.key
}

func (name RuleNameOperator) equals(other Operator) error {
	otherName, ok := other.(RuleNameOperator)
	if !ok {
		return fmt.Errorf("other is not of the same type: %s", reflect.TypeOf(other))
	}
	if name.key != otherName.key {
		return fmt.Errorf("keys do not match: expected %s, got %s", name.key, otherName.key)
	}
	return nil
}

// parseRuleName converts a raw (nested) rulename node to a (more) readable one.
// ABNF: rulename = ALPHA *(ALPHA / DIGIT / "-")
func parseRuleName(rawNode *operators.Node) Operator {
	return RuleNameOperator{
		key: rawNode.String(),
	}
}

// parseGroup converts a raw (nested) group node to a (more) readable one.
// ABNF: group = "(" *c-wsp alternation *c-wsp ")"
func parseGroup(rawNode *operators.Node) Operator {
	return parseAlternation(rawNode.GetSubNode("alternation"))
}

// OptionOperator represents an option node of a rule.
type OptionOperator struct {
	key         string
	subOperator Operator
}

func (opt OptionOperator) Key() string {
	return opt.key
}

func (opt OptionOperator) equals(other Operator) error {
	otherOpt, ok := other.(OptionOperator)
	if !ok {
		return fmt.Errorf("other is not of the same type: %s", reflect.TypeOf(other))
	}
	if opt.key != otherOpt.key {
		return fmt.Errorf("keys do not match: expected %s, got %s", opt.key, otherOpt.key)
	}
	return opt.subOperator.equals(otherOpt.subOperator)
}

// parseOption converts a raw (nested) option node to a (more) readable one.
// ABNF: option = "[" *c-wsp alternation *c-wsp "]"
func parseOption(rawNode *operators.Node) Operator {
	rawAlternation := rawNode.GetSubNode("alternation")
	return OptionOperator{
		key:         rawNode.String(),
		subOperator: parseAlternation(rawAlternation),
	}
}

// CharacterValueOperator represents a character value node of a rule.
type CharacterValueOperator struct {
	value string
}

func (value CharacterValueOperator) Key() string {
	return value.value
}

func (value CharacterValueOperator) equals(other Operator) error {
	otherValue, ok := other.(CharacterValueOperator)
	if !ok {
		return fmt.Errorf("other is not of the same type: %s", reflect.TypeOf(other))
	}
	if value.value != otherValue.value {
		return fmt.Errorf("subValues do not match: expected %s, got %s", value.value, otherValue.value)
	}
	return nil
}

// parseCharacterValue converts a raw (nested) character value node to a (more) readable one.
// ABNF: char-val = DQUOTE *(%x20-21 / %x23-7E) DQUOTE
func parseCharacterValue(rawNode *operators.Node) Operator {
	rawValue := rawNode.GetSubNode("*(%x20-21 / %x23-7E)")
	return CharacterValueOperator{
		value: rawValue.String(),
	}
}

// NumericValueOperator represents a numeric value node of a rule.
type NumericValueOperator struct {
	key            string
	hyphen, points bool
	numericType    numericType
	value          []string
}

func (value NumericValueOperator) Key() string {
	return value.key
}

func (value NumericValueOperator) equals(other Operator) error {
	otherValue, ok := other.(NumericValueOperator)
	if !ok {
		return fmt.Errorf("other is not of the same type: %s", reflect.TypeOf(other))
	}
	if value.key != otherValue.key {
		return fmt.Errorf("keys do not match: expected %s, got %s", value.key, otherValue.key)
	}
	if value.numericType != otherValue.numericType {
		return fmt.Errorf(
			"numeric types do not match: expected %s, got %s",
			value.numericType, otherValue.numericType,
		)
	}
	if len(value.value) != len(otherValue.value) {
		return fmt.Errorf(
			"lenght of values do not match: expected %d, got %d",
			len(value.value), len(otherValue.value),
		)
	}
	if value.hyphen != otherValue.hyphen || value.points != otherValue.points {
		return fmt.Errorf(
			"value types do not match: expected -%t .%t, got -%t .%t",
			value.hyphen, value.points, otherValue.hyphen, otherValue.points,
		)
	}
	for i, part := range value.value {
		if part != otherValue.value[i] {
			return fmt.Errorf(
				"value parts %d do not match: expected %s, got %s", i, part, otherValue.value[i])
		}
	}
	return nil
}

// parseCharacterValue converts a raw (nested) numeric value node to a (more) readable one.
// ABNF: num-val = "%" (bin-val / dec-val / hex-val)
func parseNumericValue(rawNode *operators.Node) Operator {
	rawValue := rawNode.Children[1].Children[0]
	var numericType numericType
	switch rawValue.Key {
	case "bin-val":
		numericType = binary
	case "dec-val":
		numericType = decimal
	case "hex-val":
		numericType = hexadecimal
	}
	values, hasHyphen, hasPoints := make([]string, 0), false, false
	for _, child := range rawValue.Children {
		if rawNumericValue := child.GetNode(string(numericType)); rawNumericValue != nil {
			// encountered hyphen
			if child.Contains("-") {
				hasHyphen = true
			}
			// encountered point(s)
			if child.Contains(".") {
				hasPoints = true
			}

			if hasHyphen {
				values = append(values, rawNumericValue.String())
			} else if hasPoints {
				for _, part := range rawNumericValue.GetSubNodes(string(numericType)) {
					values = append(values, part.String())
				}
			} else {
				values = []string{rawNumericValue.String()}
			}
		}
	}

	return NumericValueOperator{
		key:         rawNode.String(),
		hyphen:      hasHyphen,
		points:      hasPoints,
		numericType: numericType,
		value:       values,
	}
}

type numericType string

const (
	binary      numericType = "1*BIT"
	decimal     numericType = "1*DIGIT"
	hexadecimal numericType = "1*HEXDIG"
)

func (value NumericValueOperator) toIntegers() [][]int {
	bytes := make([][]int, len(value.value))
	switch value.numericType {
	case binary:
		for i, part := range value.value {
			raw, _ := strconv.ParseInt(part, 2, 64)
			bytes[i] = []int{int(raw)}
		}
	case decimal:
		for i, part := range value.value {
			raw, _ := strconv.Atoi(part)
			bytes[i] = []int{raw}
		}
	case hexadecimal:
		for i, part := range value.value {
			bytes[i] = hexStringToBytes(part)
		}
	default:
		panic("invalid numeric type")
	}
	return bytes
}

func hexStringToBytes(hexStr string) []int {
	n, _ := strconv.ParseInt(hexStr, 16, 64)
	b := make([]int, (len(hexStr)+1)/2)
	for i := range b {
		b[i] = int(byte(n >> uint64(8*(len(b)-i-1))))
	}
	return b
}
