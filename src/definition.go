package pineapple

type Lexer struct {
	sourceCode       string
	lineNum          int
	nextToken        string
	nextTokenType    int
	nextTokenLineNum int
}
