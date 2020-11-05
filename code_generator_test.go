package abnf

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCodeGenerator_core(t *testing.T) {
	rawABNF, err := ioutil.ReadFile("./testdata/core.abnf")
	if err != nil {
		t.Error(err)
		return
	}

	g := CodeGenerator{
		PackageName: "core",
		RawABNF:     rawABNF,
	}
	b := &bytes.Buffer{}
	g.writer = b
	g.isOperator = true
	g.generate()

	coreABNF, err := ioutil.ReadFile("./core/core_abnf.go")
	if err != nil {
		t.Error(err)
		return
	}

	var (
		expected = strings.Split(string(coreABNF), "\n")
		actual   = strings.Split(b.String(), "\n")
	)

	for row, expectedLine := range expected {
		actualLine := []rune(actual[row])
		if len(expectedLine) != len(actualLine) {
			t.Errorf("Lines %d do not have an equal length", row+1)
			continue
		}
		for column, expectedChar := range expectedLine {
			actualChar := actualLine[column]
			if expectedChar != actualChar {
				t.Errorf(
					"Characters do not match on %d-%d: %s != %s",
					row+1, column+1, string(expectedChar), string(actualChar),
				)
				break
			}
		}
	}
}

func TestCodeGenerator_definition(t *testing.T) {
	rawABNF, err := ioutil.ReadFile("./testdata/definition.abnf")
	if err != nil {
		t.Error(err)
		return
	}

	corePkg := ExternalABNF{
		IsOperator:  true,
		PackageName: "core",
		PackagePath: "github.com/elimity-com/abnf/core",
	}
	g := CodeGenerator{
		PackageName: "definition",
		RawABNF:     rawABNF,
		ExternalABNF: map[string]ExternalABNF{
			"ALPHA":  corePkg,
			"BIT":    corePkg,
			"CRLF":   corePkg,
			"DIGIT":  corePkg,
			"DQUOTE": corePkg,
			"HEXDIG": corePkg,
			"VCHAR":  corePkg,
			"WSP":    corePkg,
		},
	}
	b := &bytes.Buffer{}
	g.writer = b
	g.generate()

	coreABNF, err := ioutil.ReadFile("./definition/abnf_definition.go")
	if err != nil {
		t.Error(err)
		return
	}

	var (
		expected = strings.Split(string(coreABNF), "\n")
		actual   = strings.Split(b.String(), "\n")
	)

	for row, expectedLine := range expected {
		actualLine := []rune(actual[row])
		if len(expectedLine) != len(actualLine) {
			t.Errorf("Lines %d do not have an equal length", row+1)
			continue
		}
		for column, expectedChar := range expectedLine {
			actualChar := actualLine[column]
			if expectedChar != actualChar {
				t.Errorf(
					"Characters do not match on %d-%d: %s != %s",
					row+1, column+1, string(expectedChar), string(actualChar),
				)
				break
			}
		}
	}
}
