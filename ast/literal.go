package ast

import (
	"fmt"
	"log"
	"stone/env"
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

func (l *NumberLiteral) Eval(env env.Env) any {
	return l.Token.GetValue()
}

type IdentifierLiteral struct {
	ASTLeaf
}

func NewIdentifierLiteral(token *lexer.Token) *IdentifierLiteral {
	nl := &IdentifierLiteral{}
	nl.Token = token
	return nl
}

func (l *IdentifierLiteral) Eval(env env.Env) any {
	v, err := env.Get(l.Token.GetValue().(string))
	if err != nil {
		log.Panicf("ReferenceError line %4v: %v %v", l.LineNumber(), l.Token.GetValue().(string), "变量未定义")
	}
	return v
}

func (l *IdentifierLiteral) IsVar() bool {
	return true
}

type StringLiteral struct {
	ASTLeaf
}

func NewStringLiteral(token *lexer.Token) *StringLiteral {
	nl := &StringLiteral{}
	nl.Token = token
	return nl
}

func (sl *StringLiteral) Eval(env env.Env) any {
	return sl.Token.GetValue().(string)
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

func (bl *BooleanLiteral) Eval(env env.Env) any {
	return bl.Token.GetValue().(bool)
}
