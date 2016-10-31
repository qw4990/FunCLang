package ez_lang

import (
	"fmt"
	"strconv"
)

func RegisterDefaultFuncs(ezi EZInterpreter) {
	ezi.RegisterFunc("Add", Add)
	ezi.RegisterFunc("Sub", Sub)
	ezi.RegisterFunc("Eq", Eq)
	ezi.RegisterFunc("Lt", Lt)
	ezi.RegisterFunc("Gt", Gt)
	ezi.RegisterFunc("And", And)
	ezi.RegisterFunc("Or", Or)
	ezi.RegisterFunc("Println", Println)
	ezi.RegisterFunc("Printf", Printf)
}

var EZTrue EZVar = &ezVar{
	typex: _EZVAR_TYPE_NUM,
	num:   1,
}

// EZFalse :zero number or empty string will be regared as false, see function IsTrue
var EZFalse EZVar = &ezVar{
	typex: _EZVAR_TYPE_NUM,
	num:   0,
}

func Add(args ...EZVar) EZVar {
	if len(args) != 2 {
		panic("the length of args is not 2")
	}

	a, b := args[0], args[1]
	if a.IsNum() && b.IsNum() {
		return newEZNum(a.Num() + b.Num())
	}

	var strA string
	if a.IsNum() {
		strA = strconv.FormatFloat(a.Num(), 'f', 3, 64)
	} else {
		strA = a.Str()
	}

	var strB string
	if b.IsNum() {
		strB = strconv.FormatFloat(b.Num(), 'f', 3, 64)
	} else {
		strB = b.Str()
	}

	return newEZStr(strA + " " + strB)
}

func Sub(args ...EZVar) EZVar {
	if len(args) != 2 {
		panic("the length of args is not 2")
	}

	a, b := args[0], args[1]
	if a.IsNum() && b.IsNum() {
		return newEZNum(a.Num() - b.Num())
	}
	if a.IsStr() {
		panic(a.Str() + " is not a number")
	}
	panic(b.Str() + " is not a number")
}

func Mul(args ...EZVar) EZVar {
	if len(args) != 2 {
		panic("the length of args is not 2")
	}

	a, b := args[0], args[1]
	if a.IsNum() && b.IsNum() {
		return newEZNum(a.Num() * b.Num())
	}
	if a.IsStr() {
		panic(a.Str() + " is not a number")
	}
	panic(b.Str() + " is not a number")
}

func Div(args ...EZVar) EZVar {
	if len(args) != 2 {
		panic("the length of args is not 2")
	}

	a, b := args[0], args[1]
	if a.IsNum() && b.IsNum() {
		if b.Num() == 0 {
			panic("divise zero")
		}
		return newEZNum(a.Num() / b.Num())
	}
	if a.IsStr() {
		panic(a.Str() + " is not a number")
	}
	panic(b.Str() + " is not a number")
}

func Eq(args ...EZVar) EZVar {
	if len(args) != 2 {
		panic("the length of args is not 2")
	}

	if args[0].IsNum() && args[1].IsNum() {
		if args[0].Num() == args[1].Num() {
			return EZTrue
		}
	} else if args[0].IsStr() && args[1].IsStr() {
		if args[0].Str() == args[1].Str() {
			return EZTrue
		}
	}
	return EZFalse
}

func Gt(args ...EZVar) EZVar {
	if len(args) != 2 {
		panic("the length of args is not 2")
	}

	if args[0].IsNum() && args[1].IsNum() {
		if args[0].Num() > args[1].Num() {
			return EZTrue
		}
	}
	return EZFalse
}

func Lt(args ...EZVar) EZVar {
	if len(args) != 2 {
		panic("the length of args is not 2")
	}

	if args[0].IsNum() && args[1].IsNum() {
		if args[0].Num() < args[1].Num() {
			return EZTrue
		}
	}
	return EZFalse
}

func Println(args ...EZVar) EZVar {
	for _, arg := range args {
		if arg.IsNum() {
			fmt.Print(arg.Num())
		} else if arg.IsStr() {
			fmt.Print(arg.Str())
		}
		fmt.Print(" ")
	}
	fmt.Println()
	return nil
}

func Printf(args ...EZVar) EZVar {
	if len(args) == 0 {
		return nil
	}
	if args[0].IsStr() == false {
		return Println(args...)
	}

	format := args[0].Str()

	fmt.Println(format, " <<<")

	is := make([]interface{}, 0, len(args)-1)
	for i := 1; i < len(args); i++ {
		if args[i].IsNum() {
			is = append(is, args[i].Num())
		} else {
			is = append(is, args[i].Str())
		}
	}
	fmt.Printf(format, is...)
	return nil
}

func And(args ...EZVar) EZVar {
	for _, a := range args {
		if IsTure(a) == false {
			return EZFalse
		}
	}
	return EZTrue
}

func Or(args ...EZVar) EZVar {
	for _, a := range args {
		if IsTure(a) == true {
			return EZTrue
		}
	}
	return EZFalse
}
