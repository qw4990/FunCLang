package func_lang

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

		Println(x)
		return 5`
	reader := strings.NewReader(codes)

	fc := NewFunCaller()
	result, err := fc.Call(reader)
	if err != nil {
		t.Fatal(err)
	}

	if result.Num() != 5 {
		t.Fatal("Error")
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

	fc := NewFunCaller()
	_, err := fc.Call(reader)
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

		return "happy"
		`
	reader := strings.NewReader(codes)

	fc := NewFunCaller()
	result, err := fc.Call(reader)
	if err != nil {
		t.Fatal(err)
	}
	if result.Str() != "happy" {
		t.Fatal("error")
	}
}

func TestReturn(t *testing.T) {
	codes := `
		return "happy"
		Println("hhh")
		return "not happy"
		`
	reader := strings.NewReader(codes)

	fc := NewFunCaller()
	result, err := fc.Call(reader)
	if err != nil {
		t.Fatal(err)
	}
	if result.Str() != "happy" {
		t.Fatal("error")
	}
}

func TestSingleQuote(t *testing.T) {
	codes := `
		str := 'hello'
		Println(str)
	`
	reader := strings.NewReader(codes)

	fc := NewFunCaller()
	_, err := fc.Call(reader)
	if err != nil {
		t.Fatal("error")
	}
}
