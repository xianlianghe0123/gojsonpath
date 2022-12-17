package parser

type Token struct {
	TokenType tokenType
	Value     string
}

type tokenType int

const (
	TokenStart tokenType = iota
	TokenRoot
	TokenAll
	TokenDot
	TokenRecursion
	TokenField
	TokenSquare
)
