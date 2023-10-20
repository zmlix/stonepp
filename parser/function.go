package parser

import (
	"log"
	"stonepp/ast"
	"stonepp/lexer"
)

// param     : Identifier
// params    : param { "," param }
// param_list: "(" [ params ] ")"
// def       : "def" Identifier param_list block

func paramParser() ast.ASTNode {
	var left ast.ASTNode
	if tokenUtils.isType(lexer.Identifier) {
		left = ast.NewIdentifierLiteral(tokenUtils.token())
		tokenUtils.next()
	} else {
		return nil
	}
	return left
}

func paramsParser() ast.ASTNode {
	var left ast.ASTNode
	var params []ast.ASTNode
	left = paramParser()
	if left == nil {
		return nil
	}
	params = append(params, left)
	for tokenUtils.isToken(",", lexer.Symbol) {
		tokenUtils.next()
		left = paramParser()
		if left == nil {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "\",\"后缺少参数")
		}
		params = append(params, left)
	}
	left = ast.NewParams(params)
	return left
}

func paramListParser() ast.ASTNode {
	var left ast.ASTNode
	if tokenUtils.isToken("(", lexer.Symbol) {
		tokenUtils.next()
		left = paramsParser()
		if !tokenUtils.isToken(")", lexer.Symbol) {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少\")\"")
		}
		tokenUtils.next()
		if left == nil {
			return ast.NewParamList([]ast.ASTNode{})
		}
	} else {
		return nil
	}
	left = ast.NewParamList([]ast.ASTNode{left})
	return left
}

func defParser() ast.ASTNode {
	skipEOL()
	var left ast.ASTNode
	var name ast.ASTNode
	var param_list ast.ASTNode
	var block ast.ASTNode
	if tokenUtils.isToken("def", lexer.Symbol) {
		tokenUtils.next()
		if tokenUtils.isType(lexer.Identifier) {
			name = ast.NewIdentifierLiteral(tokenUtils.token())
			tokenUtils.next()
		} else {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少函数名")
		}
		param_list = paramListParser()
		block = blockParser()
		return ast.NewDefStmnt([]ast.ASTNode{name, param_list, block})
	}

	return left
}
