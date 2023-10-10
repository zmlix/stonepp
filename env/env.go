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
	if !ok {
		return de.Father().Get(name)
	}
	return value, nil
}

func (de *DefEnv) Set(name string, value any) {
	if _, ok := de.Father().(*GlobalEnv); ok {
		if _, ok2 := de.VarMap[name]; ok2 {
			de.VarMap[name] = value
		} else {
			de.Father().Set(name, value)
		}
	}
	de.VarMap[name] = value
}

func (de *DefEnv) Father() Env {
	return de.FatherEnv
}

func (de *DefEnv) Has(name string) bool {
	_, ok := de.VarMap[name]
	return ok || de.Father().Has(name)
}
