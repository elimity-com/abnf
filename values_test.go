package abnf

import (
	"fmt"
	"strings"
	"testing"
)

func TestBinValues(t *testing.T) {
	for _, test := range []struct {
		abnf     string
		contains string
	}{
		{
			abnf: "b = %b0\n",
			contains: "Terminal(\"b\", []byte{0})",
		},
		{
			abnf: "b = %b1000001\n",
			contains: "Terminal(\"b\", []byte{65})",
		},
		{
			abnf:"b = %b1000001-1011010\n",
			contains: "Range(\"b\", []byte{65}, []byte{90})",
		},
		{
			abnf:"b = %b1000001.1000010.1000011\n",
			contains: "String(\"b\", \"ABC\")",
		},
	} {
		t.Run("BinVal", func(t *testing.T) {
			g := Generator{
				PackageName:  "num",
				RawABNF:      test.abnf,
			}
			f := g.GenerateABNFAsOperators()
			if !strings.Contains(fmt.Sprintf("%#v", f), test.contains) {
				fmt.Printf("%#v", f)
				t.Errorf("did not parse correctly")
			}
		})
	}
}

func TestDecValues(t *testing.T) {
	for _, test := range []struct {
		abnf     string
		contains string
	}{
		{
			abnf: "d = %d0\n",
			contains: "Terminal(\"d\", []byte{0})",
		},
		{
			abnf: "d = %d65\n",
			contains: "Terminal(\"d\", []byte{65})",
		},
		{
			abnf:"d = %d65-90\n",
			contains: "Range(\"d\", []byte{65}, []byte{90})",
		},
		{
			abnf:"d = %d65.66.67\n",
			contains: "String(\"d\", \"ABC\")",
		},
	} {
		t.Run("DecVal", func(t *testing.T) {
			g := Generator{
				PackageName:  "num",
				RawABNF:      test.abnf,
			}
			f := g.GenerateABNFAsOperators()
			if !strings.Contains(fmt.Sprintf("%#v", f), test.contains) {
				t.Errorf("did not parse correctly")
			}
		})
	}
}

func TestHexValues(t *testing.T) {
	for _, test := range []struct {
		abnf     string
		contains string
	}{
		{
			abnf: "x = %x0\n",
			contains: "Terminal(\"x\", []byte{0})",
		},
		{
			abnf: "x = %x41\n",
			contains: "Terminal(\"x\", []byte{65})",
		},
		{
			abnf:"x = %x41-5A\n",
			contains: "Range(\"x\", []byte{65}, []byte{90})",
		},
		{
			abnf:"x = %x41.42.43\n",
			contains: "String(\"x\", \"ABC\")",
		},
		{
			abnf:"x = %x3C0\n",
			contains: "Terminal(\"x\", []byte{3, 192})",
		},
	} {
		t.Run("HexVal", func(t *testing.T) {
			g := Generator{
				PackageName:  "num",
				RawABNF:      test.abnf,
			}
			f := g.GenerateABNFAsOperators()
			if !strings.Contains(fmt.Sprintf("%#v", f), test.contains) {
				fmt.Printf("%#v", f)
				t.Errorf("did not parse correctly")
			}
		})
	}
}
