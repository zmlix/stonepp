package ast

import (
	"fmt"
	"stone/lexer"
)

type NumberLiteral struct {
	ASTLeaf
}

func NewNumberLiteral(token *lexer.Token) *NumberLiteral {
	nl := &NumberLiteral{}
	nl.Token = token
	return nl
}

type IdentifierLiteral struct {
	ASTLeaf
}

func NewIdentifierLiteral(token *lexer.Token) *IdentifierLiteral {
	nl := &IdentifierLiteral{}
	nl.Token = token
	return nl
}

type StringLiteral struct {
	ASTLeaf
}

func NewStringLiteral(token *lexer.Token) *StringLiteral {
	nl := &StringLiteral{}
	nl.Token = token
	return nl
}

func (sl *StringLiteral) String() string {
	return fmt.Sprintf("\"%v\"", sl.Token.GetValue())
}

type BooleanLiteral struct {
	ASTLeaf
}

func NewBooleanLiteral(token *lexer.Token) *BooleanLiteral {
	nl := &BooleanLiteral{}
	nl.Token = token
	return nl
}
