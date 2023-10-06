package parser

import (
	"log"
	"stone/ast"
	"stone/lexer"
)

// simple    : expr [ args ]
// block     : "{" [ statement ] { EOL [ statement ] } "}"
// statement : "if" expr block { "elif" expr block } [ "else" block ]
//           | "while" expr block
//           | simple
//           | "return" [ expr ]

func simpleParser() ast.ASTNode {
	var simple ast.ASTNode

	simple = exprParser()

	return simple
}

func blockParser() ast.ASTNode {
	var left ast.ASTNode
	var blocks []ast.ASTNode
	if tokenUtils.isToken("{", lexer.Symbol) {
		tokenUtils.next()
		left = statementParser()
		if left != nil {
			blocks = append(blocks, left)
			for tokenUtils.isType(lexer.EOL) {
				left = statementParser()
				if left != nil {
					blocks = append(blocks, left)
				} else {
					break
				}
			}
		}
		skipEOL()
		if !tokenUtils.isToken("}", lexer.Symbol) {
			log.Fatalln("line:", tokenUtils.token().GetLineNumber(), "缺少}")
		}
		tokenUtils.next()
	} else {
		log.Fatalln("line:", tokenUtils.token().GetLineNumber(), "缺少{")
	}
	left = ast.NewBlockStmnt(blocks)
	return left
}

func statementParser() ast.ASTNode {
	skipEOL()
	var left ast.ASTNode
	var ifCond ast.ASTNode
	var thenBlock ast.ASTNode
	var elseBlock ast.ASTNode
	var elif []ast.ASTNode

	if tokenUtils.isToken("if", lexer.Symbol) {
		tokenUtils.next()
		ifCond = exprParser()
		if ifCond == nil {
			log.Fatalln("line:", tokenUtils.token().GetLineNumber(), "缺少条件")
		}
		thenBlock = blockParser()
		for tokenUtils.isToken("elif", lexer.Symbol) {
			tokenUtils.next()
			left = exprParser()
			if left == nil {
				log.Fatalln("line:", tokenUtils.token().GetLineNumber(), "缺少条件")
			}
			elif = append(elif, left)
			left = blockParser()
			elif = append(elif, left)
		}
		if tokenUtils.isToken("else", lexer.Symbol) {
			tokenUtils.next()
			elseBlock = blockParser()
		} else {
			elseBlock = ast.NewBlockStmnt([]ast.ASTNode{})
		}
		left = ast.NewIfStmnt(ifCond, thenBlock, elseBlock, elif)
	} else if tokenUtils.isToken("while", lexer.Symbol) {
		tokenUtils.next()
		ifCond = exprParser()
		if ifCond == nil {
			log.Fatalln("line:", tokenUtils.token().GetLineNumber(), "缺少条件")
		}
		thenBlock = blockParser()
		left = ast.NewWhileStmnt(ifCond, thenBlock)
	} else if tokenUtils.isToken("return", lexer.Symbol) {
		tokenUtils.next()
		thenBlock = exprParser()
		if thenBlock != nil {
			left = ast.NewReturnExpr(thenBlock)
		} else {
			left = ast.NewReturnExpr(nil)
		}
	} else {
		left = simpleParser()
	}
	return left
}
