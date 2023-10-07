package lexer

import (
	"golang.org/x/exp/slices"
)

func deleteComment(code_ string) string {
	var code string = ""
	for i := 0; i < len(code_); i++ {
		if i+1 < len(code_) && code_[i] == '/' && code_[i+1] == '/' {
			for i < len(code_) && code_[i] != '\n' {
				i++
			}
		}
		if i < len(code_) {
			code += string(code_[i])
		}
	}
	// code += "\n"
	return code
}

func format(code_ string) string {
	var code string = ""

	isString := false
	for i := 0; i < len(code_); i++ {
		// fmt.Println(isString, i, string(code_[i]), slices.Contains(symbols, string(code_[i])))
		if code_[i] == '"' {
			if isString && i-1 >= 0 && code_[i-1] != '\\' {
				isString = false
			} else {
				isString = true
			}
		}
		if !isString && code_[i] != '"' && slices.Contains(symbols, string(code_[i])) {
			if i+1 < len(code_) && slices.Contains(symbols, code_[i:i+2]) {
				code += " " + code_[i:i+2] + " "
				i += 1
			} else if code_[i] == '-' {
				space := false
				for k := i - 1; k >= 0; k-- {
					if code_[k] == ' ' {
						continue
					}
					if !slices.Contains(symbols, string(code_[k])) {
						space = true
					}
					break
				}
				if i+1 < len(code_) && code_[i+1] == '"' {
					space = true
				}
				code += " " + string(code_[i])
				if space {
					code += " "
				}
			} else {
				code += " " + string(code_[i]) + " "
			}
			continue
		}
		code += string(code_[i])
	}

	return code
}

func Preprocessor(code string) string {
	code = deleteComment(code)
	code = format(code)
	return code
}
