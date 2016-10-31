package func_lang

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type mockVar struct {
	str string
	num float64
}

func (mv *mockVar) Str() string {
	return mv.str
}

func (mv *mockVar) Num() float64 {
	return mv.num
}

func (mv *mockVar) IsNum() bool {
	return mv.str == ""
}

func (mv *mockVar) IsStr() bool {
	return mv.num == 0
}

func (mv *mockVar) IsTrue() bool {
	if mv.IsStr() {
		return mv.Str() != ""
	}
	return mv.Num() != 0
}

func TestDemo0(t *testing.T) {
	sum := func(vs ...Var) Var {
		result := 0.0
		for _, v := range vs {
			if v.IsNum() {
				result += v.Num()
			} else {
				panic(fmt.Sprintf("%v is non a number", v.Str())) // use panic to report error
			}
		}
		return &mockVar{"", result}
	}

	fc := NewFunCaller()
	fc.RegisterFunc("Sum", sum)

	code := `
        a := 1
        b := 2
        c := 3
        d := 4
        Println(Sum(a, b, c, d, 5, 6, 7, 8, 9, 10))
    `
	reader := strings.NewReader(code)
	fc.Call(reader)
}

func TestDemo1(t *testing.T) {
	add := func(vs ...Var) Var {
		result := ""
		for _, v := range vs {
			if v.IsNum() {
				result = result + " " + strconv.FormatFloat(v.Num(), 'f', 2, 64)
			} else {
				result = result + " " + v.Str()
			}
		}
		return &mockVar{result, 0}
	}

	fc := NewFunCaller()
	fc.RegisterFunc("Add", add)

	code := `
        a := 1
        b := 2
        c := 3
        d := 4
        Println(Add(a, b, c, d, 5, 6, 7, 8, 9, 10, "Hello World"))
    `
	reader := strings.NewReader(code)
	fc.Call(reader)
}
