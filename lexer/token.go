package lexer

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

type TokenType int

const (
	Identifier TokenType = iota
	Number
	String
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
	fmt.Printf("%#v\n", token)
}

func printTokens(tokens []*Token) {
	for _, token := range tokens {
		token.Print()
	}
}

func parseNumber(word string) (any, error) {

	if word[0] == '-' {
		value, err := parseNumber(word[1:])
		if err != nil {
			return 0, err
		}
		switch v := value.(type) {
		case int:
			return -v, err
		case float64:
			return -v, err
		}
	}

	pointCnt := 0
	for i := 0; i < len(word); i++ {
		if !unicode.IsNumber(rune(word[i])) && word[i] != '.' {
			return 0, fmt.Errorf("%q 数值不合法", word)
		}
		if word[i] == '.' {
			pointCnt++
			if pointCnt > 1 {
				return 0, fmt.Errorf("%q 数值不合法", word)
			}
		}
		if word[i] == '-' {
			return 0, fmt.Errorf("%q 数值不合法", word)
		}
	}

	if pointCnt == 1 {
		var value float64 = 0
		point := 0
		base := 1.0
		for i := 0; i < len(word); i++ {
			if word[i] == '.' {
				point = 1
				continue
			}
			if point == 0 {
				value = value*10 + float64(word[i]-'0')
			} else {
				base = base / 10
				value = value + float64(word[i]-'0')*base
			}
		}
		return value, nil
	} else {
		var value int
		for i := 0; i < len(word); i++ {
			value = value*10 + int(word[i]-'0')
		}
		return value, nil
	}
}

func parseIdentifier(word string) (string, error) {

	for _, c := range word {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_' {
			return "", fmt.Errorf("%q 是不合法的标识符", word)
		}
	}

	return word, nil
}

func parseString(word string) (string, error) {

	if word[0] != '"' || word[len(word)-1] != '"' {
		return "", fmt.Errorf("%v 字符串解析错误", word)
	}

	for i := 1; i < len(word)-1; i++ {
		if word[i] == '"' {
			return "", fmt.Errorf("%v 字符串解析错误", word)
		}
		if word[i] == '\\' {
			i++
		}
	}
	return word[1 : len(word)-1], nil
}

func ParseToken(code string) []*Token {
	lines := strings.Split(code, "\n")
	var words [][]string
	for _, line := range lines {
		words = append(words, strings.Fields(line))
	}
	fmt.Printf("words: %d, %q\n", len(words), words)

	symbols := []string{
		"+", "-", "*", "/", ">", "<", ">=", "<=", "==", "=", "!=", "<<", ">>", "%",
		"{", "}", "[", "]", "(", ")", "\"", "'", "||", "&&", "|", "&"}
	var tokens []*Token
	for lineNumber, line := range words {
		// fmt.Printf("lineNumber: %v\n", lineNumber)
		// fmt.Printf("line: %q\n", line)
		for _, word := range line {
			// fmt.Printf("word: %v\n", word)
			if slices.Contains(symbols, word) {
				tokens = append(tokens, &Token{lineNumber + 1, Symbol, word})
			} else if word[0] == '_' || unicode.IsLetter(rune(word[0])) {
				tokenValue, err := parseIdentifier(word)
				if err != nil {
					log.Fatalln(err)
				}
				tokens = append(tokens, &Token{lineNumber + 1, Identifier, tokenValue})
			} else if word[0] == '-' || unicode.IsNumber(rune(word[0])) {
				tokenValue, err := parseNumber(word)
				if err != nil {
					log.Fatalln(err)
				}
				tokens = append(tokens, &Token{lineNumber + 1, Number, tokenValue})
			} else if word[0] == '"' {
				tokenValue, err := parseString(word)
				if err != nil {
					log.Fatalln(err)
				}
				tokens = append(tokens, &Token{lineNumber, String, tokenValue})
			}
		}
		tokens = append(tokens, &Token{lineNumber + 1, EOL, -1})
	}
	tokens = append(tokens, &Token{len(words), EOF, -2})
	printTokens(tokens)
	return tokens
}
