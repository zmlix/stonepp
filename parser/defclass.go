package parser

import (
	"log"
	"stone/ast"
	"stone/lexer"
)

// member    : def | simple
// class_body: "{" [ member ] { EOL [ member ] } "}"
// defclass  : "class" Identifier [ "extends" Identifier ] class_body

func memberParser() ast.ASTNode {
	var def ast.ASTNode
	var simple ast.ASTNode
	def = defParser()
	if def != nil {
		return def
	}
	simple = simpleParser()
	if simple != nil {
		return simple
	}
	return nil
}

func classBodyParser() ast.ASTNode {
	var left ast.ASTNode
	var member ast.ASTNode
	var memberList []ast.ASTNode

	if tokenUtils.isToken("{", lexer.Symbol) {
		tokenUtils.next()
		member = memberParser()
		if member != nil {
			memberList = append(memberList, member)
			for tokenUtils.isType(lexer.EOL) {
				skipEOL()
				member = memberParser()
				if member != nil {
					memberList = append(memberList, member)
				} else {
					break
				}
			}
		}
		skipEOL()
		if !tokenUtils.isToken("}", lexer.Symbol) {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少}")
		}
		tokenUtils.next()
	}
	left = ast.NewClassBody(memberList)
	return left
}

func defclassParser() ast.ASTNode {
	var left ast.ASTNode
	var name ast.ASTNode
	var extends ast.ASTNode
	var body ast.ASTNode
	if tokenUtils.isToken("class", lexer.Symbol) {
		tokenUtils.next()
		if tokenUtils.isType(lexer.Identifier) {
			name = ast.NewIdentifierLiteral(tokenUtils.token())
			tokenUtils.next()
		} else {
			log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少类名")
		}
		if tokenUtils.isToken("extends", lexer.Symbol) {
			tokenUtils.next()
			if tokenUtils.isType(lexer.Identifier) {
				extends = ast.NewIdentifierLiteral(tokenUtils.token())
				tokenUtils.next()
			} else {
				log.Panicf("SyntaxError line %4v: %s", tokenUtils.token().GetLineNumber(), "缺少继承的父类名")
			}
		}
		body = classBodyParser()
		left = ast.NewDefClass([]ast.ASTNode{name, extends, body})
	}

	return left
}
