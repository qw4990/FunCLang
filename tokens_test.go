package ez_lang

import (
	"fmt"
	"lab/ez_lang/tokener"
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	input := strings.NewReader(`call(12.3, 4)
    x := 23.5 
    if NumEqual(x, 12) {
		y := "12 dx we" 
    }`)
	tk, err := tokener.NewSimpleTokener(_EZ_SPLITER_CHARS, ezTokenizeRules...)
	if err != nil {
		t.Fatal(err)
	}
	err = tk.Tokenize(input)
	if err != nil {
		panic(err)
	}

	for tk.HasNext() {
		token, err := tk.Next()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(token.Token())
	}
}
