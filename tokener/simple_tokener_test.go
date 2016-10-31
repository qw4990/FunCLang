package tokener

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type mockRule struct {
	typ    int
	regexp string
}

func (mr *mockRule) Type() int {
	return mr.typ
}

func (mr *mockRule) RegExp() string {
	return mr.regexp
}

func TestSimpleTokenerTest(t *testing.T) {
	r0 := &mockRule{0, "type"}
	r1 := &mockRule{1, "struct"}
	r2 := &mockRule{2, "string"}
	r3 := &mockRule{3, "int"}
	r4 := &mockRule{4, "func"}
	r5 := &mockRule{5, "[*]"}
	r6 := &mockRule{6, "[(]"}
	r7 := &mockRule{7, "[)]"}
	r8 := &mockRule{8, "[{]"}
	r9 := &mockRule{9, "[}]"}
	r10 := &mockRule{10, "[.]"}
	r11 := &mockRule{11, "[a-zA-Z_]+"}
	r12 := &mockRule{12, "\".*\""}

	input := strings.NewReader(`type simpleToken struct {
		token string
		typeX int
	}
	
	aaa bbb ccc
	"aaa bbb ccc"
	append
	`)

	type resultToken struct {
		token string
		typ   int
	}

	results := []*resultToken{
		&resultToken{"type", 0},
		&resultToken{"simpleToken", 11},
		&resultToken{"struct", 1},
		&resultToken{"{", 8},
		&resultToken{"token", 11},
		&resultToken{"string", 2},
		&resultToken{"typeX", 11},
		&resultToken{"int", 3},
		&resultToken{"}", 9},
		&resultToken{"aaa", 11},
		&resultToken{"bbb", 11},
		&resultToken{"ccc", 11},
		&resultToken{`"aaa bbb ccc"`, 12},
		&resultToken{`append`, 11}}

	tker, err := NewSimpleTokener(" \t\n",
		r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11, r12)
	if err != nil {
		t.Fatal("Err")
	}

	err = tker.Tokenize(input)
	if err != nil {
		t.Fatal("Err")
	}

	for {
		tk, err := tker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		if len(results) == 0 {
			t.Fatal("Error")
		}

		fmt.Println(tk.Token())

		if tk.Type() != results[0].typ || tk.Token() != results[0].token {
			t.Fatal(tk.Token(), tk.Type(), results[0].token, results[0].typ)
		}
		results = results[1:]
	}
}
