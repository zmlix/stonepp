package parser

import (
	"log"
	"stone/ast"
	"stone/lexer"

	"golang.org/x/exp/slices"
)

// elements  : expr { "," expr }
// args      : expr { "," expr }
// postfix   : "." Identifier | "(" [ args ] ")" | "[" expr "]"
// primary   : ("fun" param_list block | "[" [ elements ] "]" | "(" expr ")" | Number | Identifier | String | Boolean) { postfix }
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

func elementsParser() ast.ASTNode {
	var left ast.ASTNode
	var elements []ast.ASTNode
	left = exprParser()
	if left == nil {
		return ast.NewElements([]ast.ASTNode{})
	}
	elements = append(elements, left)
	for tokenUtils.isToken(",", lexer.Symbol) {
		tokenUtils.next()
		left = exprParser()
		if left == nil {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "\",\"后缺少参数")
		}
		elements = append(elements, left)
	}
	left = ast.NewElements(elements)
	return left
}

func argsParser() ast.ASTNode {
	var left ast.ASTNode
	var args []ast.ASTNode
	left = exprParser()
	if left == nil {
		return left
	}
	args = append(args, left)
	for tokenUtils.isToken(",", lexer.Symbol) {
		tokenUtils.next()
		left = exprParser()
		if left == nil {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "\",\"后缺少参数")
		}
		args = append(args, left)
	}
	left = ast.NewArgs(args)
	return left
}

func postfixParser() ast.ASTNode {
	var left ast.ASTNode
	if tokenUtils.isToken("[", lexer.Symbol) {
		tokenUtils.next()
		left = exprParser()
		if left == nil {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少索引下标")
		}
		if !tokenUtils.isToken("]", lexer.Symbol) {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少\"]\"")
		}
		tokenUtils.next()
		left = ast.NewArrayRef(left)
	} else if tokenUtils.isToken("(", lexer.Symbol) {
		tokenUtils.next()
		left = argsParser()
		if !tokenUtils.isToken(")", lexer.Symbol) {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少\")\"")
		}
		tokenUtils.next()
		if left == nil {
			left = ast.NewPostfix([]ast.ASTNode{})
		}
	} else if tokenUtils.isToken(".", lexer.Symbol) {
		tokenUtils.next()
		if tokenUtils.isType(lexer.Identifier) {
			left = ast.NewStringLiteral(tokenUtils.token())
			left = ast.NewDot(left)
			tokenUtils.next()
		} else {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少调用方法")
		}
	} else {
		return nil
	}
	left = ast.NewPostfix([]ast.ASTNode{left})
	return left
}

func primaryParser() ast.ASTNode {
	var left ast.ASTNode
	var postfix ast.ASTNode
	var param_list ast.ASTNode
	var block ast.ASTNode
	if tokenUtils.isType(lexer.Identifier) && slices.Contains(ast.Native, tokenUtils.token().GetValue().(string)) {
		left = ast.NewNativeFunction(tokenUtils.token())
	} else if tokenUtils.isToken("fun", lexer.Symbol) {
		tokenUtils.next()
		param_list = paramListParser()
		block = blockParser()
		left = ast.NewFunction(nil, param_list, block, nil, ast.Lambda)
		tokenUtils.back()
	} else if tokenUtils.isToken("[", lexer.Symbol) {
		tokenUtils.next()
		left = elementsParser()
		if !tokenUtils.isToken("]", lexer.Symbol) {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少\"]\"")
		}
	} else if tokenUtils.isToken("(", lexer.Symbol) {
		tokenUtils.next()
		left = exprParser()
		if !tokenUtils.isToken(")", lexer.Symbol) {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少\")\"")
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

	var primaryList []ast.ASTNode
	primaryList = append(primaryList, left)

	for {
		postfix = postfixParser()
		if postfix == nil {
			break
		}
		primaryList = append(primaryList, postfix)
	}

	return ast.NewPrimaryExpr(primaryList)
}

func factorParser() ast.ASTNode {
	var left ast.ASTNode
	if tokenUtils.isToken("+", lexer.Symbol) {
		log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "意外出现的\"+\"")
	}
	if tokenUtils.isToken("-", lexer.Symbol) {
		tokenUtils.next()
		left = factorParser()
		if left != nil {
			left = ast.NewNegativeExpr([]ast.ASTNode{left})
		}
	} else {
		left = primaryParser()
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
