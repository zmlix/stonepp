package parser

import (
	"log"
	"stone/ast"
	"stone/lexer"

	"golang.org/x/exp/slices"
)

type TokenUtils struct {
	tokens []*lexer.Token
	pos    int
}

func (tu *TokenUtils) isToken(value any, tokenType lexer.TokenType) bool {
	// fmt.Println(tu.pos, tu.tokens[tu.pos])
	return tu.tokens[tu.pos].GetValue() == value && tu.tokens[tu.pos].GetType() == tokenType
}

func (tu *TokenUtils) isOpToken() bool {

	var Ops = []string{
		"+", "-", "*", "/", ">", "<", ">=", "<=", "==", "=", "!=", "<<", ">>", "^",
		"%", "||", "&&", "|", "&", "!"}

	return slices.Contains(Ops, tu.tokens[tu.pos].GetValue().(string))
}

func (tu *TokenUtils) isType(tokenType lexer.TokenType) bool {
	return tu.tokens[tu.pos].GetType() == tokenType
}

func (tu *TokenUtils) token() *lexer.Token {
	return tu.tokens[tu.pos]
}

func (tu *TokenUtils) next() {
	tu.pos++
	if tu.pos >= len(tu.tokens) {
		log.Fatalln("已读取全部token")
	}
}

var tokenUtils *TokenUtils

func skipEOL() {
	for tokenUtils.isType(lexer.EOL) {
		tokenUtils.next()
	}
}

func programParser() ast.ASTNode {
	var left ast.ASTNode
	left = defParser()
	if left != nil {
		return left
	}
	left = statementParser()
	return left
}

func Parser(tokens []*lexer.Token) []ast.ASTNode {
	for _, token := range tokens {
		token.Print()
	}
	tokenUtils = &TokenUtils{tokens, 0}

	var astNodes []ast.ASTNode
	var node ast.ASTNode
	for {
		node = programParser()
		if node == nil {
			break
		}
		astNodes = append(astNodes, node)
	}

	return astNodes
}