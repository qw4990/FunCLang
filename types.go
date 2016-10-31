package func_lang

const (
	_EZVAR_TYPE_NUM = iota
	_EZVAR_TYPE_STR = iota
	_EZVAR_TYPE_NIL = iota
)

var fcNil = &fcVar{
	typex: _EZVAR_TYPE_NIL,
}

type fcVar struct {
	typex int
	num   float64
	str   string
}

func cloneVar(v Var) *fcVar {
	ev := new(fcVar)
	if v.IsNum() {
		ev.typex = _EZVAR_TYPE_NUM
		ev.num = v.Num()
	} else {
		ev.typex = _EZVAR_TYPE_STR
		ev.str = v.Str()
	}
	return ev
}

func newFCNum(num float64) *fcVar {
	return &fcVar{
		typex: _EZVAR_TYPE_NUM,
		num:   num,
	}
}

func newFCStr(str string) *fcVar {
	return &fcVar{
		typex: _EZVAR_TYPE_STR,
		str:   str,
	}
}

func (ev *fcVar) IsNum() bool {
	return ev.typex == _EZVAR_TYPE_NUM
}

func (ev *fcVar) Num() float64 {
	if ev.IsNum() == false {
		return 0 // return zero type of float64
	}
	return ev.num
}

func (ev *fcVar) IsStr() bool {
	return ev.typex == _EZVAR_TYPE_STR
}

func (ev *fcVar) Str() string {
	if ev.IsStr() == false {
		return "" // return zero type of string
	}
	return ev.str
}

func (ev *fcVar) IsTrue() bool {
	if ev.IsNum() {
		return ev.Num() != 0
	}
	if ev.IsStr() {
		return ev.Str() != ""
	}
	return false
}

func (ev *fcVar) clone(v *fcVar) {
	ev.typex = v.typex
	ev.num = v.num
	ev.str = v.str
}

type fcScope struct {
	varTable  map[string]*fcVar
	funcTable map[string]Func
	parent    *fcScope
}

func newFCScope(parent *fcScope) *fcScope {
	return &fcScope{
		varTable:  make(map[string]*fcVar),
		funcTable: make(map[string]Func),
		parent:    parent,
	}
}

// lookupVar lookups a var in this scope and its ancestors
func (fs *fcScope) lookupVar(name string) *fcVar {
	if v, ok := fs.varTable[name]; ok {
		return v
	}
	if fs.parent != nil {
		return fs.parent.lookupVar(name)
	}
	return nil // not found
}

// lookupVar lookups a func in this scope and its ancestors
func (fs *fcScope) lookupFunc(name string) Func {
	if f, ok := fs.funcTable[name]; ok {
		return f
	}
	if fs.parent != nil {
		return fs.parent.lookupFunc(name)
	}
	return nil // not found
}
