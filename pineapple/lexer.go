package pineapple

import (
	"fmt"
	"regexp"
	"strings"
)

// EBNF:
// SourceCharacter ::=  #x0009 | #x000A | #x000D | [#x0020-#xFFFF]
// Name            ::= [_A-Za-z][_0-9A-Za-z]*
// StringCharacter ::= SourceCharacter - '"'
// String          ::= '"' '"' Ignored | '"' StringCharacter '"' Ignored
// Variable        ::= "$" Name Ignored
// Assignment      ::= Variable Ignored "=" Ignored String Ignored
// Print           ::= "print" "(" Ignored Variable Ignored ")" Ignored
// Statement       ::= Print | Assignment
// SourceCode      ::= Statement+
// Ignored        ::= WhiteSpace | LineTerminator
// WhiteSpace     ::= '\t' | ' ' /* ASCII: \t | Space, Horizontal Tab (U+0009), Space (U+0020) */
// LineTerminator ::= '\n' | '\r' | '\r\n'   /* ASCII: \n | \r\n | \r, New Line (U+000A) | Carriage Return (U+000D) [Lookahead != New Line (U+000A)] | Carriage Return (U+000D)New Line (U+000A) */

// tokens
const (
	TOKEN_EOF         = iota // end-of-file
	TOKEN_VAR_PREFIX         // $
	TOKEN_LEFT_PAREN         // (
	TOKEN_RIGHT_PAREN        // )
	TOKEN_EQUAL              // =
	TOKEN_QUOTE              // "
	TOKEN_DOUQUOTE           // ""
	TOKEN_NAME               // Name ::= [_A-Za-z][_0-9A-Za-z]*
	TOKEN_PRINT              // print
	TOKEN_IGNORED            // white space or line terminator
)

// map for better print
var tokenNameMap = map[int]string{
	TOKEN_EOF:         "EOF",
	TOKEN_VAR_PREFIX:  "$",
	TOKEN_LEFT_PAREN:  "(",
	TOKEN_RIGHT_PAREN: ")",
	TOKEN_EQUAL:       "=",
	TOKEN_QUOTE:       "\"",
	TOKEN_DOUQUOTE:    "\"\"",
	TOKEN_NAME:        "Name",
	TOKEN_PRINT:       "print",
	TOKEN_IGNORED:     "Ignored",
}

var keywords = map[string]int{
	"print": TOKEN_PRINT,
}

// lexer struct
type Lexer struct {
	sourceCode       string
	lineNum          int
	nextToken        string
	nextTokenType    int
	nextTokenLineNum int
}

func NewLexer(sourceCode string) *Lexer {
	return &Lexer{sourceCode, 1, "", 0, 0} // start at line 1
}

func (lexer *Lexer) GetLineNum() int {
	return lexer.lineNum
}

func (lexer *Lexer) MatchToken() (lineNum int, tokenType int, token string) {
	// check ignored
	if lexer.isIgnored() {
		return lexer.lineNum, TOKEN_IGNORED, "Ignored"
	}

	// finish
	if len(lexer.sourceCode) == 0 {
		return lexer.lineNum, TOKEN_EOF, tokenNameMap[TOKEN_EOF]
	}

	// check token
	switch lexer.sourceCode[0] {
	case '$':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_VAR_PREFIX, "$"
	case '(':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_LEFT_PAREN, "("
	case ')':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_RIGHT_PAREN, ")"
	case '=':
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_EQUAL, "="
	case '"':
		if lexer.nextSourceCodeIs("\"\"") {
			lexer.skipSourceCode(2)
			return lexer.lineNum, TOKEN_DOUQUOTE, "\"\""
		}
		lexer.skipSourceCode(1)
		return lexer.lineNum, TOKEN_QUOTE, "\""
	}

	// check multiple character token
	if lexer.sourceCode[0] == '_' || isLetter(lexer.sourceCode[0]) {
		token := lexer.scanName()
		if tokenType, isMatch := keywords[token]; isMatch {
			return lexer.lineNum, tokenType, token
		} else {
			return lexer.lineNum, TOKEN_NAME, token
		}
	}

	// unexpected symbol
	err := fmt.Sprintf("MatchToken(): unexpected symbol near '%q'.", lexer.sourceCode[0])
	panic(err)
	return
}

// a pack for MatchToken function
func (lexer *Lexer) GetNextToken() (lineNum int, tokenType int, token string) {
	if lexer.nextTokenLineNum > 0 {
		lineNum = lexer.nextTokenLineNum
		tokenType = lexer.nextTokenType
		token = lexer.nextToken
		lexer.lineNum = lexer.nextTokenLineNum
		lexer.nextTokenLineNum = 0
		return
	}
	return lexer.MatchToken()
}

// Ignored        ::= WhiteSpace | LineTerminator
// WhiteSpace     ::= '\t' | ' ' /* ASCII: \t | Space, Horizontal Tab (U+0009), Space (U+0020) */
/* ASCII: \n | \r\n | \r, New Line (U+000A) | Carriage Return (U+000D) [Lookahead != New Line (U+000A)] | Carriage Return (U+000D)New Line (U+000A) */
// LineTerminator ::= '\n' | '\r' | '\r\n'   /* ASCII: \n | \r\n | \r, New Line (U+000A) | Carriage Return (U+000D) [Lookahead != New Line (U+000A)] | Carriage Return (U+000D)New Line (U+000A) */
func (lexer *Lexer) isIgnored() bool {
	isIgnored := false
	// target pattern
	isNewLine := func(c byte) bool {
		return c == '\r' || c == '\n'
	}
	isWhiteSpace := func(c byte) bool {
		switch c {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			return true
		}
		return false
	}
	// matching
	for len(lexer.sourceCode) > 0 {
		if lexer.nextSourceCodeIs("\r\n") || lexer.nextSourceCodeIs("\n\r") {
			lexer.skipSourceCode(2)
			lexer.lineNum += 1
			isIgnored = true
		} else if isNewLine(lexer.sourceCode[0]) {
			lexer.skipSourceCode(1)
			lexer.lineNum += 1
			isIgnored = true
		} else if isWhiteSpace(lexer.sourceCode[0]) {
			lexer.skipSourceCode(1)
			isIgnored = true
		} else {
			break
		}
	}
	return isIgnored
}

func (lexer *Lexer) nextSourceCodeIs(s string) bool {
	return strings.HasPrefix(lexer.sourceCode, s)
}

// String ::= '"' StringCharacter '"'
// Finding the `StringCharacter`
func (lexer *Lexer) scanBeforeToken(token string) string {
	s := strings.Split(lexer.sourceCode, token)
	if len(s) < 2 {
		panic("unreachable!")
		return ""
	}
	lexer.skipSourceCode(len(s[0]))
	return s[0]
}

var regexName = regexp.MustCompile(`^[_\d\w]+`)

func (lexer *Lexer) scanName() string {
	return lexer.scan(regexName)
}

func (lexer *Lexer) scan(regexp *regexp.Regexp) string {
	if token := regexp.FindString(lexer.sourceCode); token != "" {
		lexer.skipSourceCode(len(token))
		return token
	}
	panic("unreachable")
	return ""
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func (lexer *Lexer) skipSourceCode(n int) {
	lexer.sourceCode = lexer.sourceCode[n:]
}

// handlers

// assert the next token & move forward to "eat" more tokens
func (lexer *Lexer) NextTokenIs(tokenType int) (lineNum int, token string) {
	nowLineNum, nowTokenType, nowToken := lexer.GetNextToken()
	// syntax error
	if tokenType != nowTokenType {
		err := fmt.Sprintf("NextTokenIs(): syntax error near '%s', expected token {%s} but got {%s}.", tokenNameMap[nowTokenType], tokenNameMap[tokenType], tokenNameMap[nowTokenType])
		panic(err)
	}
	return nowLineNum, nowToken
}

// return the next token type, but no moving the cursor
func (lexer *Lexer) LookAhead() int {
	// lexer.nextToken* is already set
	if lexer.nextTokenLineNum > 0 {
		return lexer.nextTokenType
	}
	// set it
	nowLineNum := lexer.lineNum
	lineNum, tokenType, token := lexer.GetNextToken()
	// cursor moving backwards
	lexer.lineNum = nowLineNum
	lexer.nextTokenLineNum = lineNum
	lexer.nextTokenType = tokenType
	lexer.nextToken = token
	return tokenType
}

// only for `Ignored`, just like lexer.LookAhead()
func (lexer *Lexer) LookAheadAndSkip(expectedType int) {
	// get next token
	nowLineNum := lexer.lineNum
	lineNum, tokenType, token := lexer.GetNextToken()
	// not expected type, reverse cursor
	if tokenType != expectedType {
		lexer.lineNum = nowLineNum
		lexer.nextTokenLineNum = lineNum
		lexer.nextTokenType = tokenType
		lexer.nextToken = token
	}
}
