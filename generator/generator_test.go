package generator

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestGenerateCore(t *testing.T) {
	rawABNF, err := ioutil.ReadFile("../testdata/core.abnf")
	if err != nil {
		t.Error(err)
		return
	}
	f := GenerateABNF("core", string(rawABNF))
	raw, err := ioutil.ReadFile("../core/core_abnf.go")
	if err != nil {
		t.Error(err)
		return
	}

	// NOTE: 1-index based line/position numbers!
	originalSplit := strings.Split(fmt.Sprintf("%#v", f), "\n")
	generatedSplit := strings.Split(string(raw), "\n")
	if len(originalSplit) != len(generatedSplit) {
		t.Error("no equal amount of lines")
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
