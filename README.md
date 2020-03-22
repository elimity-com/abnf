# Augmented BNF for Syntax Specifications: ABNF
[RFC 5234](https://tools.ietf.org/html/rfc5234)

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