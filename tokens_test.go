package func_lang

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qw4990/func_lang/tokener"
)

func TestTokenize(t *testing.T) {
	input := strings.NewReader(`
	Lt(Rand(), 5)

	call(12.3, 4)
    x := 23.5 
    if NumEqual(x, 12) {
		y := "12 dx we" 
    }`)
	tk, err := tokener.NewSimpleTokener(_SPLITER_CHARS, tokenizeRules...)
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
