package env

import (
	"fmt"
)

type Env interface {
	Get(string) (any, error)
	Set(string, any)
	Father() Env
	Has(string) bool
}

type GlobalEnv struct {
	VarMap map[string]any
}

func NewGlobalEnv() *GlobalEnv {
	ge := &GlobalEnv{VarMap: make(map[string]any)}
	return ge
}

func (ge *GlobalEnv) Get(name string) (any, error) {
	value, ok := ge.VarMap[name]
	if !ok {
		return nil, fmt.Errorf("%v 变量未定义", name)
	}
	return value, nil
}

func (ge *GlobalEnv) Set(name string, value any) {
	_, ok := ge.VarMap[name]
	if ok || !ge.Has(name) {
		ge.VarMap[name] = value
	} else {
		ge.Father().Set(name, value)
	}
}

func (ge *GlobalEnv) Has(name string) bool {
	_, ok := ge.VarMap[name]
	return ok
}

func (ge *GlobalEnv) Father() Env {
	return nil
}

type DefEnv struct {
	FatherEnv Env
	VarMap    map[string]any
}

func NewDefEnv(fa Env) *DefEnv {
	de := &DefEnv{FatherEnv: fa}
	de.VarMap = make(map[string]any)
	return de
}

func (de *DefEnv) Get(name string) (any, error) {
	value, ok := de.VarMap[name]
	if !ok && de.Father() != nil {
		return de.Father().Get(name)
	}
	return value, nil
}

func (de *DefEnv) Set(name string, value any) {
	_, ok := de.VarMap[name]
	if ok {
		de.VarMap[name] = value
	} else {
		var env Env
		env = de
		for env.Father() != nil {
			env = env.Father()
			switch env.(type) {
			case *GlobalEnv:
				if env.Has(name) {
					env.Set(name, value)
				} else {
					de.VarMap[name] = value
				}
				return
			case *DefClassEnv:
				if env.Has(name) {
					env.Set(name, value)
				} else {
					de.VarMap[name] = value
				}
				return
			}
		}
	}
}

func (de *DefEnv) Father() Env {
	return de.FatherEnv
}

func (de *DefEnv) Has(name string) bool {
	_, ok := de.VarMap[name]
	return ok || (de.Father() != nil && de.Father().Has(name))
}

type DefClassEnv struct {
	DefEnv
}

func NewDefClassEnv(fa Env) *DefClassEnv {
	dce := &DefClassEnv{}
	dce.FatherEnv = fa
	dce.VarMap = make(map[string]any)
	return dce
}

func (dce *DefClassEnv) Get(name string) (any, error) {
	value, ok := dce.VarMap[name]
	if !ok && dce.Father() != nil {
		return dce.Father().Get(name)
	}
	return value, nil
}
