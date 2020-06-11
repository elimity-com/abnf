# Augmented BNF for Syntax Specifications: ABNF
Internet technical specifications often need to define a formal syntax and are free to employ whatever notation their 
authors deem useful. Over the years, a modified version of Backus-Naur Form (BNF), called Augmented BNF (ABNF), has been
popular among many Internet specifications. It balances compactness and simplicity with reasonable representational power.

[RFC 5234](https://tools.ietf.org/html/rfc5234)

## Contents
### Generator
Both the [Core ABNF](./core/core_abnf.go) and the [ABNF Definition](./definition/abnf_definition.go) contained within this package 
where created by the generator.
```go
corePkg := externalABNF{
	operator:    true,
	packageName: "github.com/elimity-com/abnf/core",
}
g := Generator{
	PackageName: "definition",
	RawABNF:     string(rawDef),
	ExternalABNF: map[string]externalABNF{
		"ALPHA":  corePkg,
		"BIT":    corePkg,
		// etc.
	},
}
f := g.GenerateABNFAsAlternatives()
_ = ioutil.WriteFile("./definition/abnf_definition.go", []byte(fmt.Sprintf("%#v", f)), 0644)
```
##### (Currently) Not Supported
- free-form prose
- incremental alternatives

### [Core ABNF](https://godoc.org/github.com/elimity-com/abnf/core)
"Core" rules that are used variously among higher-level rules. The "core" rules might be formed into a lexical analyzer 
or simply be part of the main ruleset.
### [Operators](https://godoc.org/github.com/elimity-com/abnf/operators)
Elements form a sequence of one or more rule names and/or value definitions, combined according to the various operators
defined in this package, such as alternative and repetition.
   
## EOL
Text files created on DOS/Windows machines have different line endings than files created on Unix/Linux. 
DOS uses carriage return and line feed (`\r\n`) as a line ending, which Unix uses just line feed (`\n`).

This is why this package also allows LF which is **NOT** compliant with the specification.
```abnf
CRLF = CR LF / LF
```

## Operator Precedence
[RFC 5234 3.10](https://tools.ietf.org/html/rfc5234#section-3.10)

`highest`

1. Rule name, prose-val, Terminal value
2. Comment
3. Value range
4. Repetition
5. Grouping, Optional
6. Concatenation
7. Alternative

`lowest`