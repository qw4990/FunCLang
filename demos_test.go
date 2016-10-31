package ez_lang

import (
	"strconv"
	"strings"
	"testing"
)

type mockEZVar struct {
	str string
	num float64
}

func (mv *mockEZVar) Str() string {
	return mv.str
}

func (mv *mockEZVar) Num() float64 {
	return mv.num
}

func (mv *mockEZVar) IsNum() bool {
	return mv.str == ""
}

func (mv *mockEZVar) IsStr() bool {
	return mv.num == 0
}

func TestDemo0(t *testing.T) {
	sum := func(ezVars ...EZVar) EZVar {
		result := 0.0
		for _, v := range ezVars {
			if v.IsNum() {
				result += v.Num()
			}
		}
		return &mockEZVar{"", result}
	}

	ez := NewEZInterpreter()
	ez.RegisterFunc("Sum", sum)

	code := `
        a := 1
        b := 2
        c := 3
        d := 4
        Println(Sum(a, b, c, d, 5, 6, 7, 8, 9, 10))
    `
	reader := strings.NewReader(code)
	ez.Interprete(reader)
}

func TestDemo1(t *testing.T) {
	add := func(ezVars ...EZVar) EZVar {
		result := ""
		for _, v := range ezVars {
			if v.IsNum() {
				result = result + " " + strconv.FormatFloat(v.Num(), 'f', 2, 64)
			} else {
				result = result + " " + v.Str()
			}
		}
		return &mockEZVar{result, 0}
	}

	ez := NewEZInterpreter()
	ez.RegisterFunc("Add", add)

	code := `
        a := 1
        b := 2
        c := 3
        d := 4
        Println(Add(a, b, c, d, 5, 6, 7, 8, 9, 10, "Hello World"))
    `
	reader := strings.NewReader(code)
	ez.Interprete(reader)
}
