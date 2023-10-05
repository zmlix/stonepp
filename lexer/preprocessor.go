package lexer

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
	return code
}

func Preprocessor(code string) string {

	code = deleteComment(code)
	ParseToken(code)
	return code
}
