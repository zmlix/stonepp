package ast

import (
	"fmt"
	"log"
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
	if len(is.Children) >= 3 {
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
	} else {
		log.Println(err)
	}
	if else_, err := is.Else(); err == nil {
		s += fmt.Sprintf("else %v\n", else_)
	} else {
		log.Println(err)
	}
	return s
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

type ReturnStmnt struct {
	ASTList
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

type NullStmnt struct {
	ASTList
}

func NewNullStmnt() ASTNode {
	ns := &NullStmnt{}
	return ns
}

func (ns *NullStmnt) String() string {
	return "EOF"
}
