package ast

import (
	"fmt"
	"stone/lexer"
)

type Expr struct {
	ASTList
}

func (ast *Expr) String() string {
	s := ""
	for i := 0; i < len(ast.Children); i++ {
		s += fmt.Sprintf("%v", ast.Children[i])
	}
	return " " + s + " "
}

type BinaryExpr struct {
	Expr
	Operator *lexer.Token
}

func NewBinaryExpr(op *lexer.Token, ch []ASTNode) *BinaryExpr {
	be := &BinaryExpr{}
	be.Operator = op
	be.Children = ch
	return be
}

func (ast *BinaryExpr) Left() ASTNode {
	return ast.Children[0]
}

func (ast *BinaryExpr) Op() *lexer.Token {
	return ast.Operator
}

func (ast *BinaryExpr) Right() ASTNode {
	return ast.Children[1]
}

func (ast *BinaryExpr) String() string {
	return fmt.Sprintf("(%v%v%v)", ast.Left(), ast.Op().GetValue(), ast.Right())
}

type NegativeExpr struct {
	Expr
}

func NewNegativeExpr(ch []ASTNode) *NegativeExpr {
	ne := &NegativeExpr{}
	ne.Children = ch
	return ne
}

func (ast *NegativeExpr) String() string {
	return fmt.Sprintf("(-%v)", ast.Children[0])
}

type PrimaryExpr struct {
	Expr
}

func NewPrimaryExpr(ch []ASTNode) *PrimaryExpr {
	pe := &PrimaryExpr{}
	pe.Children = ch
	return pe
}

func (ast *PrimaryExpr) String() string {
	return fmt.Sprintf("(%v)", ast.Children[0])
}
