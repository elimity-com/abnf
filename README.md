# Augmented BNF for Syntax Specifications: ABNF
Internet technical specifications often need to define a formal syntax and are free to employ whatever notation their 
authors deem useful. Over the years, a modified version of Backus-Naur Form (BNF), called Augmented BNF (ABNF), has been
popular among many Internet specifications. It balances compactness and simplicity with reasonable representational power.

[RFC 5234](https://tools.ietf.org/html/rfc5234)

## Contents

**!** `[]byte(...)` should be UTF-8 encoded!

### Function Generator
A way to generate the operators in memory.
```go
g := ParserGenerator{
	RawABNF: rawABNF,
}
functions := g.GenerateABNFAsOperators()
// e.g. functions["ALPHA"]([]byte("a"))
```
### Code Generator
Both the [Core ABNF](./core/core_abnf.go) and the [ABNF Definition](./definition/abnf_definition.go) contained within this package 
where created by the generator.
```go
corePkg := externalABNF{
	operator:    true,
	packageName: "github.com/elimity-com/abnf/core",
}
g := Generator{
	PackageName:  "definition",
	RawABNF:      rawABNF,
	ExternalABNF: map[string]ExternalABNF{
		"ALPHA":  corePkg,
		"BIT":    corePkg,
		// etc.
	},
}
f := g.GenerateABNFAsAlternatives()
// e.g. ioutil.WriteFile("./definition/abnf_definition.go", []byte(fmt.Sprintf("%#v", f)), 0644)
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

## HEXDIG
In the spec HEXDIG is case insensitive. \
i.e. `0x6e != 0x6E`
```abnf
HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
```
In this implementation it is so that `0x6e == 0x6E`.
```abnf
HEXDIG = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
               / "a" / "b" / "c" / "d" / "e" / "f"
```

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


## Automation
It is possible to specify a configuration using a YAML file to automate the code generation. To install the command:

```abnf
go install github.com/elimity-com/abnf/cmd/abnf@latest
abnf generate
```

The generate command expects to find an `abnf.yml` file in its working directory. This YAML file defines and controls
the code generation. To specify a different YAML file:

```abnf
 abnf -f ./path/to/foo.yml generate
```

Creation of an empty YAML file is done by invoking the `init` function. 

### Generating an empty YAML file
```abnf
abnf init
```

The `abnf.yml` file is written out in the current working directory. To specify a different location to write out 
the file, use the same `-f` flag:

```abnf
abnf -f ./path/to/foo.yml init
```

#### Code Generation Configuration Properties
- **version**: must be defined as "1" 
- **spec**: is the path to the ABNF specification file, e.g. `./testdata/core.abnf` 
- **gofile**: the name of the Go file to generated, e.g. `core_abnf.go`
- **package**: the name of the package for the Go file generated, e.g. `core`
- **output**: the output path where the package (folder) and Go file are written, e.g. `.`
- **generate**: the type of generation, {`operators`, `alternatives`}
- **verbose**: displays additional debugging information

The `generate` option either invokes `GenerateABNFAsOperators` or `GenerateABNFAsAlternatives`