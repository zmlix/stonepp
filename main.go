package main

import (
	"fmt"
	"log"
	"stone/lexer"
	"stone/parser"
	"stone/utils"
)

func main() {
	file := "code/sum.stone"
	code, err := utils.OpenCodeFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	code = lexer.Preprocessor(code)
	fmt.Printf("%v code: \n%v\n", file, code)
	tokens := lexer.ParseToken(code)
	nodes := parser.Parser(tokens)
	for _, node := range nodes {
		fmt.Println(node)
	}
}
