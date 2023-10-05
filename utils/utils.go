package utils

import "io/ioutil"

func OpenCodeFile(file string) (string, error) {
	code, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(code), nil
}
