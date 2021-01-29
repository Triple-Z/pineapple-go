# pineapple-go

Try to write a (toy) lexer&parser&backend by following the [tutorial](https://github.com/karminski/pineapple) which is written by [@karminski](https://github.com/karminski).

## Quick rundown

```bash
$ git clone https://github.com/Triple-Z/pineapple-go.git
$ cd pineapple-go
$ make example
This is the original code:
$hello_world = "Hello, world! This is pineapple programming language!"
print( $hello_world )

Run this program, and the following are results!
Hello, world! This is pineapple programming language!
```

:tada: You will see the results at your console!

## Pineapple EBNF

```
SourceCharacter ::=  #x0009 | #x000A | #x000D | [#x0020-#xFFFF]
Name            ::= [_A-Za-z][_0-9A-Za-z]*
StringCharacter ::= SourceCharacter - '"'
String          ::= '"' '"' Ignored | '"' StringCharacter '"' Ignored
Variable        ::= "$" Name Ignored
Assignment      ::= Variable Ignored "=" Ignored String Ignored
Print           ::= "print" "(" Ignored Variable Ignored ")" Ignored
Statement       ::= Print | Assignment
SourceCode      ::= Statement+
Ignored        ::= WhiteSpace | LineTerminator
WhiteSpace     ::= '\t' | ' ' /* ASCII: \t | Space, Horizontal Tab (U+0009), Space (U+0020) */
LineTerminator ::= '\n' | '\r' | '\r\n'   /* ASCII: \n | \r\n | \r, New Line (U+000A) | Carriage Return (U+000D) [Lookahead != New Line (U+000A)] | Carriage Return (U+000D)New Line (U+000A) */
```

## License

MIT, follow the origin repository.
