package main

import (
	"log"
	"stone/env"
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
	// fmt.Printf("%v code: \n%v\n", file, code)
	tokens := lexer.ParseToken(code)
	nodes := parser.Parser(tokens)
	env := &env.GlobalEnv{VarMap: make(map[string]any)}
	for _, node := range nodes {
		// fmt.Printf("%v %T\n", node, node)
		// utils.PrintResult(node.Eval(env))
		node.Eval(env)
	}
}
