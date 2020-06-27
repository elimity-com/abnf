package abnf

import (
	"github.com/elimity-com/abnf/operators"
	"io/ioutil"
	"testing"
)

func TestParserGeneratorGenerateABNFAsOperators(t *testing.T) {
	rawABNF, err := ioutil.ReadFile("./testdata/core.abnf")
	if err != nil {
		t.Error(err)
		return
	}
	g := ParserGenerator{
		RawABNF: rawABNF,
	}
	functions := g.GenerateABNFAsOperators()

	testRanges(t, []characterRange{
		{0, 64, false},
		{65, 90, true},
		{91, 96, false},
		{97, 122, true},
		{123, 255, false},
	}, functions["ALPHA"])

	testRanges(t, []characterRange{
		{0, 47, false},
		{48, 49, true},
		{50, 255, false},
	}, functions["BIT"])

	testRanges(t, []characterRange{
		{0, 0, false},
		{1, 127, true},
		{128, 255, false},
	}, functions["CHAR"])

	testRanges(t, []characterRange{
		{0, 12, false},
		{13, 13, true},
		{14, 255, false},
	}, functions["CR"])

	newLine := functions["CRLF"]
	if newLine([]byte("\r\n")).Best().IsEmpty() {
		t.Error("no matches found")
	}
	testRanges(t, []characterRange{
		{0, 9, false},
		{10, 10, true},
		{11, 255, false},
	}, newLine)

	testRanges(t, []characterRange{
		{0, 31, true},
		{32, 126, false},
		{127, 127, true},
		{128, 255, false},
	}, functions["CTL"])

	testRanges(t, []characterRange{
		{0, 47, false},
		{48, 57, true},
		{58, 255, false},
	}, functions["DIGIT"])

	testRanges(t, []characterRange{
		{0, 33, false},
		{34, 34, true},
		{35, 255, false},
	}, functions["DQUOTE"])

	testRanges(t, []characterRange{
		{0, 47, false},
		{48, 57, true},
		{58, 64, false},
		{65, 70, true},
		{71, 255, false},
	}, functions["HEXDIG"])

	testRanges(t, []characterRange{
		{0, 8, false},
		{9, 9, true},
		{10, 255, false},
	}, functions["HTAB"])

	testRanges(t, []characterRange{
		{0, 9, false},
		{10, 10, true},
		{11, 255, false},
	}, functions["LF"])

	lWhiteSpace := functions["LWSP"]
	if lWhiteSpace([]byte("\r\n ")).Best().IsEmpty() {
		t.Error("no matches found")
	}
	if lWhiteSpace([]byte("\n\t")).Best().IsEmpty() {
		t.Error("no matches found")
	}
	if lWhiteSpace([]byte(" ")).Best().IsEmpty() {
		t.Error("no matches found")
	}

	testRanges(t, []characterRange{
		{0, 255, true},
	}, functions["OCTET"])

	testRanges(t, []characterRange{
		{0, 31, false},
		{32, 32, true},
		{34, 255, false},
	}, functions["SP"])

	testRanges(t, []characterRange{
		{0, 32, false},
		{33, 126, true},
		{127, 255, false},
	}, functions["VCHAR"])

	whiteSpace := functions["WSP"]
	if whiteSpace([]byte(" ")).Best().IsEmpty() {
		t.Error("no matches found")
	}
	if whiteSpace([]byte("\t")).Best().IsEmpty() {
		t.Error("no matches found")
	}
}

type characterRange struct {
	min, max int
	isValid  bool
}

func testRanges(t *testing.T, ranges []characterRange, operator operators.Operator) {
	for _, test := range ranges {
		for i := test.min; i <= test.max; i++ {
			best := operator([]byte{byte(i)}).Best()
			if test.isValid && best.IsEmpty() {
				t.Errorf("no matches found: %d (%s)", i, string(i))
			}
			if !test.isValid && !best.IsEmpty() {
				t.Errorf("matches found: %v", best)
			}
		}
	}
}
