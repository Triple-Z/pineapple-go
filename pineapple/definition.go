package pineapple

type Variable struct {
	LineNum int
	Name    string
}

type Print struct {
	LineNum  int
	Variable *Variable
}

type Assignment struct {
	LineNum  int
	Variable *Variable
	String   string
}

// Statement ::= Print | Assignment
type Statement interface{}

var _ Statement = (*Print)(nil)
var _ Statement = (*Assignment)(nil)

// SourceCode      ::= Statement+
type SourceCode struct {
	LineNum    int
	Statements []Statement
}
