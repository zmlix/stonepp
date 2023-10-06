package lexer

import (
	"fmt"
)

type TokenType int

const (
	Identifier TokenType = iota
	Number
	String
	Boolean
	Symbol
	EOL
	EOF
)

type Token struct {
	lineNumber int
	tokenType  TokenType
	value      any
}

type Tokener interface {
	GetLineNumber() int
	GetType() TokenType
	GetValue() any
	SetToken(int, TokenType, any)
	Print()
}

func (token *Token) GetLineNumber() int {
	return token.lineNumber
}

func (token *Token) GetType() TokenType {
	return token.tokenType
}

func (token *Token) GetValue() any {
	return token.value
}

func (token *Token) SetToken(lineNumber int, tokenType TokenType, value any) {
	token.lineNumber = lineNumber
	token.tokenType = tokenType
	token.value = value
}

func (token *Token) Print() {
	// fmt.Printf("%#v\n", token)
	tokenTypeMap := make(map[TokenType]string)
	tokenTypeMap[Identifier] = "Identifier"
	tokenTypeMap[Number] = "Number"
	tokenTypeMap[String] = "String"
	tokenTypeMap[Boolean] = "Boolean"
	tokenTypeMap[Symbol] = "Symbol"
	tokenTypeMap[EOL] = "EOL"
	tokenTypeMap[EOF] = "EOF"
	fmt.Printf("lineNumber: %v\t", token.lineNumber)
	fmt.Printf("tokenType: %-8v\t", tokenTypeMap[token.tokenType])
	fmt.Printf("value: %#v\n", token.value)
}
