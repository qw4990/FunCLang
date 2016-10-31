package ez_lang

import "io"

// EZVar represents a variable used in EZLang
// An EZVar is either a number or a string
type EZVar interface {
	IsNum() bool
	Num() float64 // if this var is not a number, return the zero type of float: 0
	IsStr() bool
	Str() string // if this var is not a string, return an empty string: ""
}

// EZFunc x
type EZFunc func(args ...EZVar) EZVar

// EZInterpreter x
type EZInterpreter interface {
	Interprete(input io.Reader) error
	RegisterFunc(funcName string, ezFunc EZFunc)
}

// IsTure checks if this var is true in EZ_LANG
func IsTure(ev EZVar) bool {
	if ev.IsNum() {
		return ev.Num() != 0
	}
	return ev.Str() != ""
}
