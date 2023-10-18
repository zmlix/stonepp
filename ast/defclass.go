package ast

import (
	"fmt"
	"log"
	"stone/env"
)

type ClassBody struct {
	ASTList
}

func NewClassBody(ch []ASTNode) *ClassBody {
	cb := &ClassBody{}
	cb.Children = ch
	return cb
}

type DefClass struct {
	ASTList
}

func NewDefClass(ch []ASTNode) *DefClass {
	dc := &DefClass{}
	dc.Children = ch
	return dc
}

func (cb *DefClass) Name() ASTNode {
	return cb.Children[0]
}

func (cb *DefClass) Extends() ASTNode {
	return cb.Children[1]
}

func (cb *DefClass) Body() ASTNode {
	return cb.Children[2]
}

func (cb *DefClass) String() string {
	if cb.Extends() != nil {
		return fmt.Sprintf("class %v extends %v{\n%v}", cb.Name(), cb.Extends(), cb.Body())
	}
	return fmt.Sprintf("class %v{\n%v}", cb.Name(), cb.Body())
}

func (cb *DefClass) Eval(env_ env.Env) any {
	classInfo := NewClassInfo(cb.Name(), cb.Extends(), cb.Body(), env_)
	env_.Set(cb.Name().Value().GetValue().(string), classInfo)
	return fmt.Sprintf("<class: %v>", cb.Name().Value().GetValue())
}

type Dot struct {
	ASTList
}

func NewDot(method ASTNode) *Dot {
	d := &Dot{}
	d.Children = append(d.Children, method)
	return d
}

func (d *Dot) Method() ASTNode {
	return d.Children[0]
}

func (d *Dot) Name() string {
	return d.Method().Eval(nil).(string)
}

func (d *Dot) String() string {
	return fmt.Sprintf(".%v", d.Name())
}

func (d *Dot) Eval(env env.Env) any {
	method, err := env.Get(d.Method().Eval(env).(string))
	if err != nil {
		log.Fatalf("SyntaxError line %4v: %v %v", d.LineNumber(), d.Method(), "成员不存在")
	}
	switch m := method.(type) {
	case *Function:
		m.ftype = Method
		m.env = env
		return m
	default:
		return m
	}
}

type ClassInfo struct {
	ASTLeaf
	name    ASTNode
	extends ASTNode
	body    ASTNode
	env     env.Env
}

func NewClassInfo(name ASTNode, extends ASTNode, body ASTNode, env_ env.Env) *ClassInfo {
	ci := &ClassInfo{name: name, extends: extends, body: body, env: env_}
	if env_ != nil {
		ci.env = env.NewDefClassEnv(env_)
		ci.body.Eval(ci.env)
	}
	return ci
}

func (ci *ClassInfo) LineNumber() int {
	return ci.name.LineNumber()
}

func (ci *ClassInfo) String() string {
	return fmt.Sprintf("<class: %v>", ci.name)
}

func (ci *ClassInfo) Eval(env_ env.Env) any {
	class_env := env.NewDefClassEnv(env_)
	for k, v := range ci.env.(*env.DefEnv).VarMap {
		class_env.Set(k, v)
	}
	return NewClassInfo(ci.name, ci.extends, ci.body, class_env)
}

func (ci *ClassInfo) EvalConstructor(env_ env.Env, p_values []any, ast ASTNode) any {
	constructorName := fmt.Sprintf("%v", ci.name)
	hasConstructor := ci.env.Has(constructorName)
	if !hasConstructor {
		return NewClassObject(ci.name, ci.extends, ci.body, ci.env)
	}
	constructor_, _ := ci.env.Get(constructorName)

	switch c := constructor_.(type) {
	case *ClassInfo:
		cc_, err := c.env.Get(constructorName)
		if err != nil {
			log.Fatalf("TypeError line %4v: %v 类成员 %v 的值为 %v %v", ci.LineNumber(), ci.name, ci.name, constructor_, "不是构造函数")
		}
		cc, ok := cc_.(*Function)
		if !ok {
			return NewClassObject(ci.name, ci.extends, ci.body, ci.env)
		}
		obj := NewClassObject(ci.name, ci.extends, ci.body, ci.env)

		p_names := cc.params.Eval(nil).([]string)
		if len(p_names) != len(p_values) {
			log.Fatalf("SyntaxError line %4v: %v 期望(%v)个 获得(%v)个", ast.LineNumber(), "参数数量不匹配", len(p_names), len(p_values))
		}
		params := make(map[string]any)
		for i := 0; i < len(p_names); i++ {
			params[p_names[i]] = p_values[i]
		}
		cc.ftype = Func
		cc.EvalFunction(obj.env, params)
		return obj
	case *Function:
		obj := NewClassObject(ci.name, ci.extends, ci.body, ci.env)

		p_names := c.params.Eval(nil).([]string)
		if len(p_names) != len(p_values) {
			log.Fatalf("SyntaxError line %4v: %v 期望(%v)个 获得(%v)个", ast.LineNumber(), "参数数量不匹配", len(p_names), len(p_values))
		}
		params := make(map[string]any)
		for i := 0; i < len(p_names); i++ {
			params[p_names[i]] = p_values[i]
		}
		c.ftype = Func
		c.EvalFunction(obj.env, params)
		return obj
	default:
		log.Fatalf("TypeError line %4v: %v 类成员 %v 的值为 %v %v", ci.LineNumber(), ci.name, ci.name, constructor_, "不是构造函数")
	}
	return nil
}

type ClassObject struct {
	ASTLeaf
	name    ASTNode
	extends ASTNode
	body    ASTNode
	env     env.Env
}

func NewClassObject(name ASTNode, extends ASTNode, body ASTNode, env_ env.Env) *ClassObject {
	co := &ClassObject{name: name, extends: extends, body: body, env: env_}
	co.env = env.NewDefClassEnv(env_.Father())
	for k, v := range env_.(*env.DefClassEnv).VarMap {
		co.env.Set(k, v)
	}
	return co
}

func (co *ClassObject) String() string {
	return fmt.Sprintf("<object: %v>", co.name)
}

func (co *ClassObject) EvalMethod(env env.Env, m *Dot) any {

	methodName := fmt.Sprintf("%v", m.Eval(env))
	method, err := co.env.Get(methodName)
	if err != nil {
		return nil
	}
	return method
}
