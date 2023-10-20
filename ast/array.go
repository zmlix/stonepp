package ast

import (
	"fmt"
	"log"
	"stone/env"
)

type ArrayRef struct {
	ASTList
	index int
}

func NewArrayRef(expr ASTNode) *ArrayRef {
	ar := &ArrayRef{}
	ar.Children = append(ar.Children, expr)
	return ar
}

func (ar *ArrayRef) Eval(env env.Env) any {
	expr := ar.Children[0].Eval(env)
	idx, ok := expr.(int)
	if !ok {
		log.Panicf("RangeError line %4v: %T %s", ar.LineNumber(), expr, "不合法的索引类型")
	}
	ar.index = idx
	return ar
}

func (ar *ArrayRef) String() string {
	return fmt.Sprintf("[%v]", ar.index)
}

func (ar *ArrayRef) EvalArrayRef(arr *Array) any {
	if arr.Has(ar.index) {
		return arr.Get(ar.index)
	}
	log.Panicf("RangeError line %4v: %s", ar.LineNumber(), "不合法的索引长度")
	return nil
}

type Elements struct {
	ASTList
}

func NewElements(ch []ASTNode) *Elements {
	e := &Elements{}
	e.Children = ch
	return e
}

func (e *Elements) String() string {
	s := ""
	for i, child := range e.Children {
		s += fmt.Sprintf("%v", child)
		if i != len(e.Children)-1 {
			s += ","
		}
	}
	return "[" + s + "]"
}

func (e *Elements) Eval(env env.Env) any {
	elements := []any{}
	for _, element := range e.Children {
		elements = append(elements, element.Eval(env))
	}
	return NewArray(elements)
}

type Array struct {
	ASTLeaf
	array []any
	len   int
	env   env.Env
}

func NewArray(arr []any) *Array {
	a := &Array{}
	a.array = arr
	a.len = len(arr)
	a.env = env.NewGlobalEnv()
	a.env.Set("len", &ArrayMethod{method: a.Len, name: "len"})
	a.env.Set("append", &ArrayMethod{method: a.Append, name: "append"})
	a.env.Set("insert", &ArrayMethod{method: a.Insert, name: "insert"})
	a.env.Set("pop", &ArrayMethod{method: a.Pop, name: "pop"})
	a.env.Set("remove", &ArrayMethod{method: a.Remove, name: "remove"})
	return a
}

func (a *Array) Has(index int) bool {
	return index >= 0 && index < a.len
}

func (a *Array) Get(index int) any {
	return a.array[index]
}

func (a *Array) String() string {
	s := ""
	for i, child := range a.array {
		s += fmt.Sprintf("%v", child)
		if i != a.len-1 {
			s += ","
		}
	}
	return "[" + s + "]"
}

func (a *Array) Len() int {
	return a.len
}

func (a *Array) Append(element any) int {
	a.array = append(a.array, element)
	a.len++
	return a.len
}
func (a *Array) Insert(element any, index int) bool {
	if !a.Has(index) {
		return false
	}
	a.array = append(a.array[:index], append([]any{element}, a.array[index:]...)...)
	a.len++
	return true
}

func (a *Array) Pop() bool {
	if a.len <= 0 {
		return false
	}
	a.array = a.array[:a.len-1]
	a.len--
	return true
}

func (a *Array) Remove(index int) bool {
	if !a.Has(index) {
		return false
	}
	a.array = append(a.array[:index], a.array[index+1:]...)
	a.len--
	return true
}

type ArrayMethod struct {
	method any
	name   string
}

func (am *ArrayMethod) Eval(params []any, ast ASTNode) any {
	switch am.name {
	case "len":
		return am.method.(func() int)()
	case "append":
		if len(params) != 1 {
			log.Panicf("SyntaxError line %4v: %v 期望(%v)个 获得(%v)个", ast.LineNumber(), "参数数量不匹配", 1, len(params))
		}
		return am.method.(func(any) int)(params[0])
	case "insert":
		if len(params) != 2 {
			log.Panicf("SyntaxError line %4v: %v 期望(%v)个 获得(%v)个", ast.LineNumber(), "参数数量不匹配", 2, len(params))
		}
		if index, ok := params[1].(int); ok {
			if !am.method.(func(any, int) bool)(params[0], index) {
				log.Panicf("RangeError line %4v: %s", ast.LineNumber(), "不合法的索引长度")
			}
		} else {
			log.Panicf("RangeError line %4v: %T %s", ast.LineNumber(), params[0], "不合法的索引类型")
		}
	case "pop":
		if !am.method.(func() bool)() {
			log.Panicf("RangeError line %4v: %s", ast.LineNumber(), "数组为空")
		}
	case "remove":
		if len(params) != 1 {
			log.Panicf("SyntaxError line %4v: %v 期望(%v)个 获得(%v)个", ast.LineNumber(), "参数数量不匹配", 1, len(params))
		}
		if index, ok := params[0].(int); ok {
			if !am.method.(func(int) bool)(index) {
				log.Panicf("RangeError line %4v: %s", ast.LineNumber(), "不合法的索引长度")
			}
		} else {
			log.Panicf("RangeError line %4v: %T %s", ast.LineNumber(), params[0], "不合法的索引类型")
		}

	}
	return nil
}
