package ast

import (
	"fmt"
	"stone/lexer"
)

type Literal struct {
	ASTLeaf
}

func (ast *Literal) String() string {
	return fmt.Sprintf("%v", ast.Token.GetValue())
}

type NumberLiteral struct {
	Literal
}

func NewNumberLiteral(token *lexer.Token) *NumberLiteral {
	nl := &NumberLiteral{}
	nl.Token = token
	return nl
}

type StringLiteral struct {
	Literal
}

type BooleanLiteral struct {
	Literal
}
