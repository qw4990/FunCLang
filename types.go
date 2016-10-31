package ez_lang

// type scope interface {
// 	LookupVar(ident string) EZVar
// 	LookupFunc(ident string) EZFunc
// }

const (
	_EZVAR_TYPE_NUM = iota
	_EZVAR_TYPE_STR = iota
	_EZVAR_TYPE_NIL = iota
)

var ezNil = &ezVar{
	typex: _EZVAR_TYPE_NIL,
}

type ezVar struct {
	typex int
	num   float64
	str   string
}

func cloneEZVar(v EZVar) *ezVar {
	ev := new(ezVar)
	if v.IsNum() {
		ev.typex = _EZVAR_TYPE_NUM
		ev.num = v.Num()
	} else {
		ev.typex = _EZVAR_TYPE_STR
		ev.str = v.Str()
	}
	return ev
}

func newEZNum(num float64) *ezVar {
	return &ezVar{
		typex: _EZVAR_TYPE_NUM,
		num:   num,
	}
}

func newEZStr(str string) *ezVar {
	return &ezVar{
		typex: _EZVAR_TYPE_STR,
		str:   str,
	}
}

func (ev *ezVar) IsNum() bool {
	return ev.typex == _EZVAR_TYPE_NUM
}

func (ev *ezVar) Num() float64 {
	if ev.IsNum() == false {
		return 0 // return zero type of float64
	}
	return ev.num
}

func (ev *ezVar) IsStr() bool {
	return ev.typex == _EZVAR_TYPE_STR
}

func (ev *ezVar) Str() string {
	if ev.IsStr() == false {
		return "" // return zero type of string
	}
	return ev.str
}

func (ev *ezVar) isTrue() bool {
	return IsTure(ev)
}

func (ev *ezVar) clone(v *ezVar) {
	ev.typex = v.typex
	ev.num = v.num
	ev.str = v.str
}

type ezScope struct {
	varTable  map[string]*ezVar
	funcTable map[string]EZFunc
	parent    *ezScope
}

func newEZScope(parent *ezScope) *ezScope {
	return &ezScope{
		varTable:  make(map[string]*ezVar),
		funcTable: make(map[string]EZFunc),
		parent:    parent,
	}
}

// lookupVar lookups a var in this scope and its ancestors
func (es *ezScope) lookupVar(name string) *ezVar {
	if v, ok := es.varTable[name]; ok {
		return v
	}
	if es.parent != nil {
		return es.parent.lookupVar(name)
	}
	return nil // not found
}

// lookupVar lookups a func in this scope and its ancestors
func (es *ezScope) lookupFunc(name string) EZFunc {
	if f, ok := es.funcTable[name]; ok {
		return f
	}
	if es.parent != nil {
		return es.parent.lookupFunc(name)
	}
	return nil // not found
}
