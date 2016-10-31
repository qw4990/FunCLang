package func_lang

import "io"

// Var represents a variable used in FunCLang
// An Var is either a number or a string
type Var interface {
	IsNum() bool
	Num() float64
	IsStr() bool
	Str() string
	IsTrue() bool
}

// Func ...
type Func func(args ...Var) Var

// FunCaller ...
type FunCaller interface {
	Call(funcBody io.Reader) (Var, error)
	RegisterFunc(funcName string, f Func)
}
