package ast

import (
	"fmt"
	"stone/env"
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
	s := fmt.Sprintf("(%v)", p.Children[0])
	return s
}

func (p *Postfix) Eval(env env.Env) any {
	if len(p.Children) == 0 {
		return []any{}
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
	name   string
	params ASTNode
	body   ASTNode
	env    env.Env
}

func NewFunction(name string, params ASTNode, body ASTNode, env env.Env) *Function {
	f := &Function{name: name, params: params, body: body, env: env}
	return f
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
	return fmt.Sprintf("<function: %v>", f.name)
}

func (f *Function) Eval(env_ env.Env, params map[string]any) any {
	new_env := f.MakeEnv(env_)
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
