package ast

import (
	"fmt"
	"log"
	"stone/env"
	"stone/lexer"
)

var Native = []string{"print", "println", "read", "len", "int", "float"}

type NativeFunction struct {
	ASTLeaf
}

func NewNativeFunction(name *lexer.Token) *NativeFunction {
	nf := &NativeFunction{}
	nf.Token = name
	return nf
}

func (nf *NativeFunction) Name() string {
	return nf.Token.GetValue().(string)
}

func (nf *NativeFunction) Eval(env env.Env) any {
	return nf
}

func (nf *NativeFunction) EvalFunction(params []any) any {
	switch nf.Name() {
	case "print":
		nf.Print(params)
	case "println":
		nf.Println(params)
	case "read":
		return nf.Read()
	case "len":
		return nf.Len(params)
	case "int":
		return nf.Int(params)
	case "float":
		return nf.Float(params)
	}
	return nil
}

func (np *NativeFunction) Print(params []any) {
	fmt.Print(params...)
}

func (np *NativeFunction) Println(params []any) {
	fmt.Println(params...)
}

func (nf *NativeFunction) Read() string {
	var value string
	_, err := fmt.Scan(&value)
	if err != nil {
		log.Panicf("%v", err)
	}
	return value
}

func (nf *NativeFunction) Len(params []any) int {
	if len(params) > 1 {
		log.Panicf("SyntaxError line %4v: %v 期望(%v)个 获得(%v)个", nf.LineNumber(), "参数数量不匹配", 1, len(params))
	}
	switch v := params[0].(type) {
	case string:
		return len(v)
	case *Array:
		return v.Len()
	default:
		log.Panicf("TypeError line %4v: %T %v", nf.LineNumber(), params[0], "不可计算长度")
	}
	return 0
}

func (nf *NativeFunction) Int(params []any) int {
	value := nf.Float(params)
	return int(value)
}

func (nf *NativeFunction) Float(params []any) float64 {
	if len(params) > 1 {
		log.Panicf("SyntaxError line %4v: %v 期望(%v)个 获得(%v)个", nf.LineNumber(), "参数数量不匹配", 1, len(params))
	}
	var value float64
	_, err := fmt.Sscanf(fmt.Sprintf("%v", params[0]), "%f", &value)
	if err != nil {
		log.Panicf("TypeError line %4v: %T %v", nf.LineNumber(), params[0], "不可转换成整数")
	}
	return value
}
