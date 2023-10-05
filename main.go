package main

import (
	"fmt"
	"log"
	"stone/lexer"
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

}
