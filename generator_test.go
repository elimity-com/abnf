package abnf

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGenerateCore(t *testing.T) {
	rawABNF, err := ioutil.ReadFile("./testdata/core.abnf")
	if err != nil {
		t.Error(err)
		return
	}
	f := GenerateABNFAsOperators("core", string(rawABNF))
	// _ = ioutil.WriteFile("./core/core_abnf.go", []byte(fmt.Sprintf("%#v", f)), 0644)

	raw, err := ioutil.ReadFile("./core/core_abnf.go")
	if err != nil {
		t.Error(err)
		return
	}

	// NOTE: 1-index based line/position numbers!
	originalSplit := strings.Split(fmt.Sprintf("%#v", f), "\n")
	generatedSplit := strings.Split(string(raw), "\n")
	if len(originalSplit) != len(generatedSplit) {
		t.Errorf("no equal amount of lines: %d, %d", len(originalSplit), len(generatedSplit))
		return
	}
	for i := range originalSplit {
		if originalSplit[i] != generatedSplit[i] {
			if len(originalSplit[i]) != len(generatedSplit[i]) {
				t.Errorf("no equal amount of characters on line %d", i+1)
				return
			}
			for j := range originalSplit[i] {
				if originalSplit[i][j] != generatedSplit[i][j] {
					t.Errorf("line %d: characters do not match at position %d", i+1, j+1)
				}
			}
		}
	}
}

func TestGenerateDefinition(t *testing.T) {
	rawDef, err := ioutil.ReadFile("./testdata/definition.abnf")
	if err != nil {
		t.Error(err)
		return
	}
	rawABNF, err := ioutil.ReadFile("./testdata/core.abnf")
	if err != nil {
		t.Error(err)
		return
	}

	f := GenerateABNFAsAlternatives("abnf", string(rawDef) + string(rawABNF))
	// _ = ioutil.WriteFile("./abnf_definition.go", []byte(fmt.Sprintf("%#v", f)), 0644)

	raw, err := ioutil.ReadFile("./abnf_definition.go")
	if err != nil {
		t.Error(err)
		return
	}

	// NOTE: 1-index based line/position numbers!
	originalSplit := strings.Split(fmt.Sprintf("%#v", f), "\n")
	generatedSplit := strings.Split(string(raw), "\n")
	if len(originalSplit) != len(generatedSplit) {
		t.Errorf("no equal amount of lines: %d, %d", len(originalSplit), len(generatedSplit))
		return
	}
	for i := range originalSplit {
		if originalSplit[i] != generatedSplit[i] {
			if len(originalSplit[i]) != len(generatedSplit[i]) {
				t.Errorf("no equal amount of characters on line %d", i+1)
				return
			}
			for j := range originalSplit[i] {
				if originalSplit[i][j] != generatedSplit[i][j] {
					t.Errorf("line %d: characters do not match at position %d", i+1, j+1)
				}
			}
		}
	}
}