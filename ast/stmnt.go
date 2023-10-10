package ast

import (
	"fmt"
	"log"
	"stone/env"
)

type BlockStmnt struct {
	ASTList
}

func NewBlockStmnt(ch []ASTNode) *BlockStmnt {
	bs := &BlockStmnt{}
	bs.Children = ch
	return bs
}

func (bs *BlockStmnt) String() string {
	s := ""

	for _, child := range bs.Children {
		s += fmt.Sprintf("%v\n", child)
	}

	return "{\n" + s + "}"
}

func (ast *BlockStmnt) Eval(env env.Env) any {
	var res any
	for _, child := range ast.Children {
		if child == nil {
			log.Fatalf("SyntaxError line %4v: %s", ast.LineNumber(), "语法错误")
		}
		res = child.Eval(env)
		r, ok := res.(*ReturnValue)
		if ok {
			return r
		}
	}
	return res
}

type IfStmnt struct {
	ASTList
}

func NewIfStmnt(cond, then, else_ ASTNode, elif []ASTNode) *IfStmnt {
	is := &IfStmnt{}
	is.Children = append(is.Children, cond)
	is.Children = append(is.Children, then)
	is.Children = append(is.Children, else_)
	is.Children = append(is.Children, elif...)
	return is
}

func (is *IfStmnt) Cond() ASTNode {
	return is.Children[0]
}

func (is *IfStmnt) Then() ASTNode {
	return is.Children[1]
}

func (is *IfStmnt) Else() (ASTNode, error) {
	if is.Children[2] != nil {
		return is.Children[2], nil
	}
	return nil, fmt.Errorf("不存在else块")
}

func (is *IfStmnt) Elif() ([]ASTNode, error) {
	if len(is.Children) >= 5 {
		return is.Children[3:], nil
	}
	return nil, fmt.Errorf("不存在elif块")
}

func (is *IfStmnt) String() string {
	s := fmt.Sprintf("if %v %v\n", is.Cond(), is.Then())
	if elif, err := is.Elif(); err == nil {
		for i := 0; i < len(elif); i += 2 {
			s += fmt.Sprintf("elif %v %v\n", elif[i], elif[i+1])
		}
	}
	if else_, err := is.Else(); err == nil {
		s += fmt.Sprintf("else %v\n", else_)
	}
	return s
}

func (is *IfStmnt) Eval(env env.Env) any {
	cond, ok := is.Cond().Eval(env).(bool)
	if !ok {
		log.Fatalf("TypeError line %4v: %s", is.LineNumber(), "条件返回值必须是\"bool\"类型")
	}
	if cond {
		return is.Then().Eval(env)
	}
	if elif, err := is.Elif(); err == nil {
		for i := 0; i < len(elif); i += 2 {
			elif_cond, elif_then := elif[i], elif[i+1]
			cond, ok := elif_cond.Eval(env).(bool)
			if !ok {
				log.Fatalf("TypeError line %4v: %s", is.LineNumber(), "条件返回值必须是\"bool\"类型")
			}
			if cond {
				return elif_then.Eval(env)
			}
		}
	}
	if else_, err := is.Else(); err == nil {
		return else_.Eval(env)
	}
	return nil
}

type WhileStmnt struct {
	ASTList
}

func NewWhileStmnt(cond, block ASTNode) *WhileStmnt {
	ws := &WhileStmnt{}
	ws.Children = append(ws.Children, cond)
	ws.Children = append(ws.Children, block)
	return ws
}

func (ws *WhileStmnt) Cond() ASTNode {
	return ws.Children[0]
}

func (ws *WhileStmnt) Block() ASTNode {
	return ws.Children[1]
}

func (ws *WhileStmnt) String() string {
	return fmt.Sprintf("while %v %v", ws.Cond(), ws.Block())
}

func (ws *WhileStmnt) Eval(env env.Env) any {
	var res any
	for {
		cond, ok := ws.Cond().Eval(env).(bool)
		if !ok {
			log.Fatalf("TypeError line %4v: %s", ws.LineNumber(), "条件返回值必须是\"bool\"类型")
		}
		if cond {
			res = ws.Block().Eval(env)
		} else {
			break
		}
	}

	return res
}

type ReturnStmnt struct {
	ASTList
}

type ReturnValue struct {
	Value any
}

func NewReturnExpr(expr ASTNode) *ReturnStmnt {
	re := &ReturnStmnt{}
	re.Children = append(re.Children, expr)
	return re
}

func (re *ReturnStmnt) Res() ASTNode {
	return re.Children[0]
}

func (re *ReturnStmnt) Empty() bool {
	return re.Children[0] == nil
}

func (re *ReturnStmnt) String() string {
	return fmt.Sprintf("return %v", re.Children[0])
}

func (re *ReturnStmnt) Eval(env env.Env) any {
	if re.Empty() {
		return &ReturnValue{Value: nil}
	}
	return &ReturnValue{re.Res().Eval(env)}
}

type DefStmnt struct {
	ASTList
}

func NewDefStmnt(ch []ASTNode) *DefStmnt {
	ds := &DefStmnt{}
	ds.Children = ch
	return ds
}

func (ds *DefStmnt) Name() string {
	return ds.Children[0].Value().GetValue().(string)
}

func (ds *DefStmnt) ParamList() ASTNode {
	return ds.Children[1]
}

func (ds *DefStmnt) Block() ASTNode {
	return ds.Children[2]
}

func (ds *DefStmnt) String() string {
	return fmt.Sprintf("def %v %v %v", ds.Name(), ds.ParamList(), ds.Block())
}

func (ds *DefStmnt) Eval(env env.Env) any {
	function := NewFunction(ds.Name(), ds.ParamList(), ds.Block(), env)
	env.Set(ds.Name(), function)
	return nil
}
