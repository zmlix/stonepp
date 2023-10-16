package ast

import (
	"fmt"
	"log"
	"stone/env"
	"stone/lexer"
)

type BinaryExpr struct {
	ASTList
	Operator *lexer.Token
}

func NewBinaryExpr(op *lexer.Token, ch []ASTNode) *BinaryExpr {
	be := &BinaryExpr{}
	be.Operator = op
	be.Children = ch
	return be
}

func (be *BinaryExpr) Left() ASTNode {
	return be.Children[0]
}

func (be *BinaryExpr) Op() *lexer.Token {
	return be.Operator
}

func (be *BinaryExpr) Right() ASTNode {
	return be.Children[1]
}

func (be *BinaryExpr) String() string {
	return fmt.Sprintf("(%v%v%v)", be.Left(), be.Op().GetValue(), be.Right())
}

func (be *BinaryExpr) Eval(env env.Env) any {
	if be == nil || be.Op() == nil || be.Left() == nil || be.Right() == nil {
		log.Fatalf("SyntaxError line %4v: %s", be.LineNumber(), "语法错误")
	}
	op := be.Op().GetValue().(string)
	if op == "=" {
		if !be.Left().IsName() {
			log.Fatalf("TypeError line %4v: %v 不可被赋值", be.LineNumber(), be.Left())
		}
		left := be.Left()
		right := be.Right().Eval(env)
		env.Set(fmt.Sprintf("%v", left), right)
		return right
	}

	typeAssert := func(value any) string {
		switch v := value.(type) {
		case int:
			return "int"
		case float64:
			return "float"
		case string:
			return "string"
		case bool:
			return "bool"
		default:
			return fmt.Sprintf("%T", v)
		}
	}

	left := be.Left().Eval(env)
	right := be.Right().Eval(env)

	leftType := typeAssert(left)
	rightType := typeAssert(right)

	if leftType == "int" {
		if rightType == "int" {
			lv := left.(int)
			rv := right.(int)
			switch op {
			case "+":
				return lv + rv
			case "-":
				return lv - rv
			case "*":
				return lv * rv
			case "/":
				return lv / rv
			case "%":
				return lv % rv
			case ">":
				return lv > rv
			case "<":
				return lv < rv
			case ">=":
				return lv >= rv
			case "<=":
				return lv <= rv
			case "==":
				return lv == rv
			case "!=":
				return lv != rv
			case "<<":
				return lv << rv
			case ">>":
				return lv >> rv
			case "^":
				return lv ^ rv
			case "|":
				return lv | rv
			case "&":
				return lv & rv
			}
		} else if rightType == "float" {
			lv := float64(left.(int))
			rv := right.(float64)
			switch op {
			case "+":
				return lv + rv
			case "-":
				return lv - rv
			case "*":
				return lv * rv
			case "/":
				return lv / rv
			case ">":
				return lv > rv
			case "<":
				return lv < rv
			case ">=":
				return lv >= rv
			case "<=":
				return lv <= rv
			case "==":
				return lv == rv
			case "!=":
				return lv != rv
			}
		}
	} else if leftType == "float" {
		if rightType == "int" {
			lv := left.(float64)
			rv := float64(right.(int))
			switch op {
			case "+":
				return lv + rv
			case "-":
				return lv - rv
			case "*":
				return lv * rv
			case "/":
				return lv / rv
			case ">":
				return lv > rv
			case "<":
				return lv < rv
			case ">=":
				return lv >= rv
			case "<=":
				return lv <= rv
			case "==":
				return lv == rv
			case "!=":
				return lv != rv
			}
		} else if rightType == "float" {
			lv := left.(float64)
			rv := right.(float64)
			switch op {
			case "+":
				return lv + rv
			case "-":
				return lv - rv
			case "*":
				return lv * rv
			case "/":
				return lv / rv
			case ">":
				return lv > rv
			case "<":
				return lv < rv
			case ">=":
				return lv >= rv
			case "<=":
				return lv <= rv
			case "==":
				return lv == rv
			case "!=":
				return lv != rv
			}
		}
	} else if leftType == "string" {
		if rightType == "string" {
			lv := left.(string)
			rv := right.(string)
			switch op {
			case "+":
				return lv + rv
			case "==":
				return lv == rv
			case "<":
				return lv < rv
			case "<=":
				return lv <= rv
			case ">":
				return lv > rv
			case ">=":
				return lv >= rv
			}
		}
	} else if rightType == "string" {
		if leftType == "string" {
			lv := left.(string)
			rv := right.(string)
			switch op {
			case "+":
				return lv + rv
			case "==":
				return lv == rv
			case "<":
				return lv < rv
			case "<=":
				return lv <= rv
			case ">":
				return lv > rv
			case ">=":
				return lv >= rv
			}
		}
	} else if leftType == "bool" {
		if rightType == "bool" {
			lv := left.(bool)
			rv := right.(bool)
			switch op {
			case "==":
				return lv == rv
			case "&&":
				return lv && rv
			case "||":
				return lv || rv
			}
		}
	} else if rightType == "bool" {
		if leftType == "bool" {
			lv := left.(bool)
			rv := right.(bool)
			switch op {
			case "==":
				return lv == rv
			case "&&":
				return lv && rv
			case "||":
				return lv || rv
			}
		}
	}

	log.Fatalf("TypeError line %4v: %v %v 不能使用\"%v\"运算符\n", be.LineNumber(), leftType, rightType, op)
	return nil
}

type NegativeExpr struct {
	ASTList
}

func NewNegativeExpr(ch []ASTNode) *NegativeExpr {
	ne := &NegativeExpr{}
	ne.Children = ch
	return ne
}

func (ne *NegativeExpr) String() string {
	return fmt.Sprintf("-(%v)", ne.Children[0])
}

func (ne *NegativeExpr) Eval(env env.Env) any {
	res := ne.Children[0].Eval(env)
	switch r := res.(type) {
	case int:
		return -r
	case float64:
		return -r
	default:
		log.Fatalf("TypeError line %4v: %T 类型不能使用\"-\"运算符\n", ne.LineNumber(), r)
	}
	return nil
}

type PrimaryExpr struct {
	ASTList
}

func NewPrimaryExpr(ch []ASTNode) *PrimaryExpr {
	pe := &PrimaryExpr{}
	pe.Children = ch
	return pe
}

func (pe *PrimaryExpr) Postfix() []ASTNode {
	return pe.Children[1:]
}

func (pe *PrimaryExpr) String() string {
	s := fmt.Sprintf("%v", pe.Children[0])
	for i := 1; i < len(pe.Children); i++ {
		s += fmt.Sprintf("(%v)", pe.Children[i])
	}
	return s
}

func (pe *PrimaryExpr) Eval(env_ env.Env) any {
	if len(pe.Children) > 1 {
		return pe.EvalSub(env_, len(pe.Children)-1)
	}
	return pe.Children[0].Eval(env_)
}

func (pe *PrimaryExpr) EvalSub(env env.Env, k int) any {
	if k == 0 {
		return pe.Children[0].Eval(env)
	}
	res := pe.EvalSub(env, k-1)
	p_values, _ := pe.Children[k].Eval(env).([]any)
	nfun, ok := res.(*NativeFunction)
	if ok {
		return nfun.EvalFunction(p_values)
	}
	fun, ok := res.(*Function)
	if !ok {
		log.Fatalf("TypeError line %4v: %T %v", pe.LineNumber(), res, "不可调用")
	}
	p_names, _ := fun.Params().Eval(env).([]string)
	if len(p_names) != len(p_values) {
		log.Fatalf("SyntaxError line %4v: %v 期望(%v)个 获得(%v)个", pe.LineNumber(), "参数数量不匹配", len(p_names), len(p_values))
	}
	params := make(map[string]any)
	for i := 0; i < len(p_names); i++ {
		params[p_names[i]] = p_values[i]
	}
	return fun.EvalFunction(env, params)
}
