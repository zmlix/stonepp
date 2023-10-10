package utils

import (
	"fmt"
	"io/ioutil"
)

func OpenCodeFile(file string) (string, error) {
	code, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(code), nil
}

func PrintResult(result any) {
	if result != nil {
		fmt.Printf("%v\n", result)
	}
}
