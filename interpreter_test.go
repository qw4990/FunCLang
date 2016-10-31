package ez_lang

import (
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	codes := `
		x := 2
		y := 1
		if y {
			Println("Yes")
			x := 3
			Println(x)
		}

		Println(x)`
	reader := strings.NewReader(codes)

	interpreter := NewEZInterpreter()
	err := interpreter.Interprete(reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSimpleLoop(t *testing.T) {
	codes := `
		x := 1
		for Lt(x, 5) {
			Println("Hello ", x)
			x = Add(x, 1)
		}
		Println(x, " <<<<<<")
		Println(Add("ssss", 3))
		Printf("Hello, %v, xxx sdf", 123)`
	reader := strings.NewReader(codes)

	interpreter := NewEZInterpreter()
	err := interpreter.Interprete(reader)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFibNum(t *testing.T) {
	codes := `
		a := 1
		b := 1
		cnt := 0
		for Lt(cnt, 20) {
			Println(a)
			c := Add(a, b)
			a = b
			b = c
			cnt = Add(cnt, 1)
		}
		`
	reader := strings.NewReader(codes)

	interpreter := NewEZInterpreter()
	err := interpreter.Interprete(reader)
	if err != nil {
		t.Fatal(err)
	}
}
