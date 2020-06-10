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
			contains: "Rune(\"b\", 0)",
		},
		{
			abnf: "b = %b1000001\n",
			contains: "Rune(\"b\", 65)",
		},
		{
			abnf:"b = %b1000001-1011010\n",
			contains: "Range(\"b\", 65, 90)",
		},
		{
			abnf:"b = %b1000001.1000010.1000011\n",
			contains: "String(\"b\", \"ABC\")",
		},
	} {
		t.Run("BinVal", func(t *testing.T) {
			g := GenerateABNFAsOperators("num", test.abnf)
			if !strings.Contains(fmt.Sprintf("%#v", g), test.contains) {
				t.Errorf("did not parse correctly")
				fmt.Printf("%#v", g)
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
			contains: "Rune(\"d\", 0)",
		},
		{
			abnf: "d = %d65\n",
			contains: "Rune(\"d\", 65)",
		},
		{
			abnf:"d = %d65-90\n",
			contains: "Range(\"d\", 65, 90)",
		},
		{
			abnf:"d = %d65.66.67\n",
			contains: "String(\"d\", \"ABC\")",
		},
	} {
		t.Run("DecVal", func(t *testing.T) {
			g := GenerateABNFAsOperators("num", test.abnf)
			if !strings.Contains(fmt.Sprintf("%#v", g), test.contains) {
				t.Errorf("did not parse correctly")
				fmt.Printf("%#v", g)
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
			contains: "Rune(\"x\", 0)",
		},
		{
			abnf: "x = %x41\n",
			contains: "Rune(\"x\", 65)",
		},
		{
			abnf:"x = %x41-5A\n",
			contains: "Range(\"x\", 65, 90)",
		},
		{
			abnf:"x = %x41.42.43\n",
			contains: "String(\"x\", \"ABC\")",
		},
	} {
		t.Run("DecVal", func(t *testing.T) {
			g := GenerateABNFAsOperators("num", test.abnf)
			if !strings.Contains(fmt.Sprintf("%#v", g), test.contains) {
				t.Errorf("did not parse correctly")
				fmt.Printf("%#v", g)
			}
		})
	}
}