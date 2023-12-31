package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"stonepp/ast"
	"stonepp/env"
	"stonepp/lexer"
	"stonepp/parser"
	"stonepp/utils"
)

var stonepp string = " ____    __                              __        __      " + "\n" +
	"/\\  _`\\ /\\ \\__                          /\\ \\      /\\ \\     " + "\n" +
	"\\ \\,\\L\\_\\ \\ ,_\\   ___     ___      __   \\_\\ \\___  \\_\\ \\___ " + "\n" +
	" \\/_\\__ \\\\ \\ \\/  / __`\\ /' _ `\\  /'__`\\/\\___  __\\/\\___  __\\" + "\n" +
	"   /\\ \\L\\ \\ \\ \\_/\\ \\L\\ \\/\\ \\/\\ \\/\\  __/\\/__/\\ \\_/\\/__/\\ \\_/" + "\n" +
	"   \\ `\\____\\ \\__\\ \\____/\\ \\_\\ \\_\\ \\____\\   \\ \\_\\     \\ \\_\\ " + "\n" +
	"    \\/_____/\\/__/\\/___/  \\/_/\\/_/\\/____/    \\/_/      \\/_/ "

func Shell(env env.Env) {
	defer func() {
		if err := recover(); err != nil {
			Shell(env)
		}
	}()

	var code string
	var code_ string
	var tokens []*lexer.Token
	var nodes []ast.ASTNode

	counter := func(s string) int {
		cnt := 0
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c == '{' {
				cnt++
			}
			if c == '}' {
				cnt--
			}
		}
		return cnt
	}

	for {
		fmt.Print(">>> ")
		reader := bufio.NewReader(os.Stdin)
		code, _ = reader.ReadString('\n')
		code = lexer.Preprocessor(code)
		for counter(code) != 0 {
			fmt.Print("... ")
			code_, _ = reader.ReadString('\n')
			code_ = lexer.Preprocessor(code_)
			code += "\n" + code_
		}

		tokens = lexer.ParseToken(code)
		nodes = parser.Parser(tokens)
		for _, node := range nodes {
			utils.PrintResult(node.Eval(env))
		}
	}
}

func RunShell() {
	fmt.Println(stonepp)
	env := &env.GlobalEnv{VarMap: make(map[string]any)}
	Shell(env)
}

func RunInterpreter(path string) {
	file := path
	code, err := utils.OpenCodeFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	code = lexer.Preprocessor(code)
	tokens := lexer.ParseToken(code)
	nodes := parser.Parser(tokens)
	env := &env.GlobalEnv{VarMap: make(map[string]any)}
	for _, node := range nodes {
		node.Eval(env)
	}
}

func main() {
	var runFlag string
	flag.StringVar(&runFlag, "run", "", "请输入文件路径")
	flag.Parse()
	if runFlag == "" {
		RunShell()
	} else {
		RunInterpreter(runFlag)
	}
}
