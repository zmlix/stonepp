package ast

import (
	"fmt"
	"log"
	"stone/env"
	"stone/lexer"
)

type ASTNode interface {
	Value() *lexer.Token
	ChildrenList() []ASTNode
	Eval(env.Env) any
	IsVar() bool
	LineNumber() int
}

type ASTLeaf struct {
	Token *lexer.Token
}

type ASTList struct {
	Children []ASTNode
}

func (ast *ASTLeaf) Value() *lexer.Token {
	return ast.Token
}

func (ast *ASTLeaf) ChildrenList() []ASTNode {
	return nil
}

func (ast *ASTLeaf) Eval(env env.Env) any {
	return ast.Token.GetValue()
}

func (ast *ASTLeaf) IsVar() bool {
	return false
}

func (ast *ASTLeaf) LineNumber() int {
	return ast.Token.GetLineNumber()
}

func (ast *ASTLeaf) String() string {
	return fmt.Sprintf("%v", ast.Token.GetValue())
}

func (ast *ASTList) Value() *lexer.Token {
	return nil
}

func (ast *ASTList) ChildrenList() []ASTNode {
	return ast.Children
}

func (ast *ASTList) Eval(env env.Env) any {
	var res any
	for _, child := range ast.Children {
		if child == nil {
			log.Fatalf("SyntaxError line %4v: %s", ast.LineNumber(), "语法错误")
		}
		res = child.Eval(env)
	}
	return res
}

func (ast *ASTList) IsVar() bool {
	return false || (len(ast.Children) == 1 && ast.Children[0].IsVar())
}

func (ast *ASTList) LineNumber() int {
	return ast.Children[0].LineNumber()
}

func (ast *ASTList) String() string {
	s := ""
	for i := 0; i < len(ast.Children); i++ {
		s += fmt.Sprintf("%v\n", ast.Children[i])
	}
	return s
}
