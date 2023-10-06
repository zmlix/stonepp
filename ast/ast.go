package ast

import (
	"fmt"
	"stone/lexer"
)

type ASTNode interface {
	Value() *lexer.Token
	ChildrenList() []ASTNode
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

func (ast *ASTLeaf) String() string {
	return fmt.Sprintf("%v", ast.Token.GetValue())
}

func (ast *ASTList) Value() *lexer.Token {
	return nil
}

func (ast *ASTList) ChildrenList() []ASTNode {
	return ast.Children
}

func (ast *ASTList) String() string {
	s := ""
	for i := 0; i < len(ast.Children); i++ {
		s += fmt.Sprintf("%v", ast.Children[i])
	}
	return " " + s + " "
}

type Name struct {
	Token *lexer.Token
}

func (ast *Name) Value() *lexer.Token {
	return ast.Token
}

func (ast *Name) ChildrenList() ASTNode {
	return nil
}
