package generator

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGenerateCore(t *testing.T) {
	raw, _ := ioutil.ReadFile("../testdata/core.abnf")
	f := GenerateABNF("core", string(raw))

	fmt.Printf("%#v", f)

	return

	_ = ioutil.WriteFile("../core/core_abnf.go",  []byte(fmt.Sprintf("%#v", f)), 0644)
}
