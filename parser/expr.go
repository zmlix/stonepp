package parser

import (
	"log"
	"stone/ast"
	"stone/lexer"
)

// args      : expr { "," expr }
// postfix   : "(" [ args ] ")"
// primary   : ("(" expr ")" | Number | Identifier | String | Boolean) { postfix }
// factor    : {"-"} primary
// expr      : factor { Op factor}

type Precedence struct {
	prec      int
	leftAssoc bool
}

var operators = map[string]Precedence{
	"!":  {prec: 0, leftAssoc: false},
	"*":  {prec: -1, leftAssoc: true},
	"/":  {prec: -1, leftAssoc: true},
	"%":  {prec: -1, leftAssoc: true},
	"+":  {prec: -2, leftAssoc: true},
	"-":  {prec: -2, leftAssoc: true},
	"<<": {prec: -3, leftAssoc: true},
	">>": {prec: -3, leftAssoc: true},
	">":  {prec: -4, leftAssoc: true},
	"<":  {prec: -4, leftAssoc: true},
	">=": {prec: -4, leftAssoc: true},
	"<=": {prec: -4, leftAssoc: true},
	"==": {prec: -5, leftAssoc: true},
	"!=": {prec: -5, leftAssoc: true},
	"&":  {prec: -6, leftAssoc: true},
	"^":  {prec: -7, leftAssoc: true},
	"|":  {prec: -8, leftAssoc: true},
	"&&": {prec: -9, leftAssoc: true},
	"||": {prec: -10, leftAssoc: true},
	"=":  {prec: -11, leftAssoc: false},
}

func primaryParser() ast.ASTNode {
	var left ast.ASTNode
	if tokenUtils.isToken("(", lexer.Symbol) {
		tokenUtils.next()
		left = exprParser()
		if !tokenUtils.isToken(")", lexer.Symbol) {
			log.Fatalln("SyntaxError line:", tokenUtils.token().GetLineNumber(), "缺少\")\"")
		}
	} else if tokenUtils.isType(lexer.Number) {
		left = ast.NewNumberLiteral(tokenUtils.token())
	} else if tokenUtils.isType(lexer.Identifier) {
		left = ast.NewIdentifierLiteral(tokenUtils.token())
	} else if tokenUtils.isType(lexer.String) {
		left = ast.NewStringLiteral(tokenUtils.token())
	} else if tokenUtils.isType(lexer.Boolean) {
		left = ast.NewBooleanLiteral(tokenUtils.token())
	} else {
		return left
	}
	tokenUtils.next()
	return left
}

func factorParser() ast.ASTNode {
	var left ast.ASTNode
	if tokenUtils.isToken("+", lexer.Symbol) {
		log.Fatalln("SyntaxError line:", tokenUtils.token().GetLineNumber(), "多余的\"+\"")
	}
	if tokenUtils.isToken("-", lexer.Symbol) {
		tokenUtils.next()
		left = factorParser()
		if left != nil {
			left = ast.NewNegativeExpr([]ast.ASTNode{left})
		}
	} else {
		left = primaryParser()
		if left != nil {
			left = ast.NewPrimaryExpr([]ast.ASTNode{left})
		}
	}
	return left
}

func exprParser() ast.ASTNode {
	skipEOL()
	var right ast.ASTNode
	var doShift func(left ast.ASTNode) ast.ASTNode
	var rightIsExpr = func(leftOp, rightOp Precedence) bool {
		if rightOp.leftAssoc {
			return leftOp.prec < rightOp.prec
		}
		return leftOp.prec <= rightOp.prec
	}
	doShift = func(left ast.ASTNode) ast.ASTNode {
		leftOp := operators[tokenUtils.token().GetValue().(string)]
		op := tokenUtils.token()
		tokenUtils.next()
		right := factorParser()
		rightOp := operators[tokenUtils.token().GetValue().(string)]
		for tokenUtils.isType(lexer.Symbol) && tokenUtils.isOpToken() && rightIsExpr(leftOp, rightOp) {
			right = doShift(right)
			rightOp = operators[tokenUtils.token().GetValue().(string)]
		}
		return ast.NewBinaryExpr(op, []ast.ASTNode{left, right})
	}

	right = factorParser()

	for tokenUtils.isType(lexer.Symbol) && tokenUtils.isOpToken() {
		right = doShift(right)
	}

	return right
}
