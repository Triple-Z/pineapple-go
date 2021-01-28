package pineapple

import "errors"

// Print ::= "print" "(" Ignored Variable Ignored ")" Ignored
func parsePrint(lexer *Lexer) (*Print, error) {
	var print Print
	var err error

	print.lineNum = lexer.GetLineNum()
	// "print"
	lexer.NextTokenIs(TOKEN_PRINT)
	// "("
	lexer.NextTokenIs(TOKEN_LEFT_PAREN)
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	// Variable
	if print.Variable, err = parseVariable(lexer); err != nil {
		return nil, err
	}
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	// ")"
	lexer.NextTokenIs(TOKEN_RIGHT_PAREN)
	lexer.LookAheadAndSkip(TOKEN_IGNORED)
	return &print, nil
}

// Statement ::= Print | Assignment
func parseStatement() {
	switch LookAhead() {
	// "print"
	case TOKEN_PRINT:
		return parsePrint()
	// "$"
	case TOKEN_VAR_PREFIX:
		return parseAssignment()
	default:
		return nil, errors.New("parseStatement(): unknown Statement")
	}
}
