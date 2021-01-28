package pineapple

import "errors"

// Print ::= "print" "(" Ignored Variable Ignored ")" Ignored
func parsePrint() {
	// "print"
	NextTokenIs(TOKEN_PRINT)
	// "("
	NextTokenIs(TOKEN_LEFT_PAREN)
	// Variable
	parseVariable()
	// ")"
	NextTokenIs(TOKEN_RIGHT_PAREN)
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
