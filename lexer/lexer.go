package lexer

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"golang.org/x/exp/slices"
)

var symbols = []string{
	"+", "-", "*", "/", ">", "<", ">=", "<=", "==", "=", "!=", "<<", ">>", "^",
	"%", "{", "}", "[", "]", "(", ")", "\"", "'", "||", "&&", "|", "&", ",", "!",
	"if", "elif", "else", "while", "def", "return", "fun", "class", "extends", "."}

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

func splitWord(lineNumber int, line string) []string {
	line = strings.TrimSpace(line)
	var words []string
	var word string
	isString := false
	for i := 0; i < len(line); i++ {
		if line[i] == '"' {
			if isString && i-1 >= 0 && line[i-1] != '\\' {
				isString = false
				words = append(words, "\""+word+"\"")
				word = ""
			} else {
				isString = true
			}
			continue
		}

		if isString {
			word += string(line[i])
			continue
		}

		if line[i] == ' ' && word != "" {
			if isString {
				log.Panicf("SyntaxError line %4v: %v", lineNumber+1, "未闭合的字符串")
			}
			words = append(words, word)
			word = ""
			continue
		}
		if line[i] != ' ' {
			word += string(line[i])
		}
	}
	if word != "" {
		if isString {
			log.Panicf("SyntaxError line %4v: %v", lineNumber+1, "未闭合的字符串")
		}
		words = append(words, word)
	}
	return words
}

func ParseToken(code string) []*Token {
	lines := strings.Split(code, "\n")
	var words [][]string
	for lineNumber, line := range lines {
		words = append(words, splitWord(lineNumber, line))
	}
	// fmt.Printf("words: %d, %q\n", len(words), words)

	var tokens []*Token
	for lineNumber, line := range words {
		for _, word := range line {
			if word == "true" || word == "false" {
				if word == "true" {
					tokens = append(tokens, &Token{lineNumber + 1, Boolean, true})
				} else {
					tokens = append(tokens, &Token{lineNumber + 1, Boolean, false})
				}
			} else if slices.Contains(symbols, word) {
				tokens = append(tokens, &Token{lineNumber + 1, Symbol, word})
			} else if word[0] == '_' || unicode.IsLetter(rune(word[0])) {
				tokenValue, err := parseIdentifier(word)
				if err != nil {
					log.Panicf("SyntaxError line %4v: %v", lineNumber+1, err)
				}
				tokens = append(tokens, &Token{lineNumber + 1, Identifier, tokenValue})
			} else if unicode.IsNumber(rune(word[0])) {
				tokenValue, err := parseNumber(word)
				if err != nil {
					log.Panicf("SyntaxError line %4v: %v", lineNumber+1, err)
				}
				tokens = append(tokens, &Token{lineNumber + 1, Number, tokenValue})
			} else if word[0] == '-' {
				if unicode.IsNumber(rune(word[1])) {
					tokenValue, err := parseNumber(word)
					if err != nil {
						log.Panicf("SyntaxError line %4v: %v", lineNumber+1, err)
					}
					tokens = append(tokens, &Token{lineNumber + 1, Number, tokenValue})
				} else if unicode.IsLetter(rune(word[1])) {
					tokenValue, err := parseIdentifier(word[1:])
					if err != nil {
						log.Panicf("SyntaxError line %4v: %v", lineNumber+1, err)
					}
					tokens = append(tokens, &Token{lineNumber + 1, Symbol, "-"})
					tokens = append(tokens, &Token{lineNumber + 1, Identifier, tokenValue})
				}
			} else if word[0] == '"' {
				tokenValue, err := parseString(word)
				if err != nil {
					log.Panicf("SyntaxError line %4v: %v", lineNumber+1, err)
				}
				tokens = append(tokens, &Token{lineNumber + 1, String, tokenValue})
			}
		}
		tokens = append(tokens, &Token{lineNumber + 1, EOL, "EOL"})
	}
	tokens = append(tokens, &Token{len(words), EOF, "EOF"})
	return tokens
}
