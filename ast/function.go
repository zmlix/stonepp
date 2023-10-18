package ast

import (
	"fmt"
	"stone/env"
)

type FunctionType int

const (
	Func FunctionType = iota
	Lambda
	Method
)

type Postfix struct {
	ASTList
}

func NewPostfix(ch []ASTNode) *Postfix {
	p := &Postfix{}
	p.Children = ch
	return p
}

func (p *Postfix) String() string {
	if len(p.Children) == 0 {
		return "()"
	}
	s := fmt.Sprintf("%v", p.Children[0])
	return s
}

func (p *Postfix) Dot() ASTNode {
	dot, ok := p.Children[0].(*Dot)
	if ok {
		return dot
	}
	return nil
}

func (p *Postfix) Eval(env env.Env) any {
	if len(p.Children) == 0 {
		return []any{}
	}
	if p.Dot() != nil {
		return p.Dot()
	}
	return p.Children[0].Eval(env)
}

type Args struct {
	ASTList
}

func NewArgs(ch []ASTNode) *Args {
	a := &Args{}
	a.Children = ch
	return a
}

func (a *Args) String() string {
	s := ""
	for i, child := range a.Children {
		s += fmt.Sprintf("%v", child)
		if i != len(a.Children)-1 {
			s += ","
		}
	}
	return s
}

func (a *Args) Eval(env env.Env) any {
	args := []any{}
	for _, arg := range a.Children {
		args = append(args, arg.Eval(env))
	}
	return args
}

type Params struct {
	ASTList
}

func NewParams(ch []ASTNode) *Params {
	p := &Params{}
	p.Children = ch
	return p
}

func (p *Params) String() string {
	s := ""
	for i, child := range p.Children {
		s += fmt.Sprintf("%v", child)
		if i != len(p.Children)-1 {
			s += ","
		}
	}
	return s
}

func (p *Params) Eval(env env.Env) any {
	params := []string{}
	for _, param := range p.Children {
		params = append(params, param.Value().GetValue().(string))
	}
	return params
}

type ParamList struct {
	ASTList
}

func NewParamList(ch []ASTNode) *ParamList {
	p := &ParamList{}
	p.Children = ch
	return p
}

func (p *ParamList) String() string {
	if len(p.Children) == 0 {
		return "()"
	}
	s := fmt.Sprintf("(%v)", p.Children[0])
	return s
}

func (p *ParamList) Eval(env env.Env) any {
	if len(p.Children) == 0 {
		return []string{}
	}
	return p.Children[0].Eval(env)
}

type Function struct {
	ASTLeaf
	name   ASTNode
	params ASTNode
	body   ASTNode
	env    env.Env
	ftype  FunctionType
}

func NewFunction(name ASTNode, params ASTNode, body ASTNode, env env.Env, ftype FunctionType) *Function {
	f := &Function{name: name, params: params, body: body, env: env, ftype: ftype}
	return f
}

func (f *Function) LineNumber() int {
	return f.name.Value().GetLineNumber()
}

func (f *Function) Params() ASTNode {
	return f.params
}

func (f *Function) Body() ASTNode {
	return f.body
}

func (f *Function) MakeEnv(env_ env.Env) env.Env {
	return env.NewDefEnv(env_)
}

func (f *Function) String() string {
	switch f.ftype {
	case Func:
		return fmt.Sprintf("<function: %v>", f.name.Value().GetValue().(string))
	case Lambda:
		return fmt.Sprintf("<fun: %v %v>", f.params, f.body)
	case Method:
		return fmt.Sprintf("<method: %v %v>", f.params, f.body)
	default:
		return fmt.Sprintf("<unkown: %T>", f)
	}
}

func (f *Function) EvalFunction(env_ env.Env, params map[string]any) any {
	var new_env env.Env
	switch f.ftype {
	case Func:
		new_env = f.MakeEnv(env_)
	case Lambda:
		new_env = f.env
	case Method:
		new_env = f.MakeEnv(env_)
	}

	for k, v := range params {
		new_env.Set(k, v)
	}
	res := f.body.Eval(new_env)
	r, ok := res.(*ReturnValue)
	if ok {
		return r.Value
	}
	return res
}

func (f *Function) Eval(env env.Env) any {
	return NewFunction(nil, f.params, f.body, env, Lambda)
}
