# FunCLang
FunCLang (Function Call Language) is an easy language written in Golang and can be nested in your Golang code;

FunCLang is easily to use and understand; The total lines of code of this project is less than 1000;

You should regard FunCLang codes as the body of your Golang function, but can be modified on runtime ;)

# Types and Variables
FunCLang has only two types: string and number(float64);

You can use &lt;declare&gt; statement to declare variables, for example:

```
name := "zhangyuanjia"
age := 22

Println(name, "'s age is", age)
```

<b>Zero number or empty string will be regared as False, and all other values are True;</b>

True and False will be used in &lt;if&gt; and &lt;for&gt; statements; Please see below;

# Interact with Golang
Each variable in FunCLang is represented by an interface in Golang:

```
type Var interface {
	IsNum() bool
	Num() float64
	IsStr() bool
	Str() string
}
```

And each function in FunCLang is a Golang function whose type is:

```
type Func func(args ...Var) Var
```

So the Println in above codes may be implemented in Golang:

```
func Println(args ...Var) Var {
	for _, arg := range args {
        if arg == nil {  // if there is error, use panic to expose it
            panic("Println err: arg can't be nil")
        }
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
```

So it's easy for you to write some functions in Golang and register them into FunCLang and then use them in FuncClang;

We will see more example soon;

# Statements
FuncClang programs consist of multiple statements;

There are only five kinds of statements in this language: &lt;if&gt;, &lt;for&gt;, &lt;declare&gt;, &lt;assign&gt;, &lt;return&gt;;

This program prints some fibonacci numbers:

```
lim := 0                            #declare
a := 1                              #declare
b := 1                              #declare
for And(Lt(b, 1000), Lt(lim, 12)) { #for
    Println(b)                      #expression, please see below
    c := Add(a, b)                  #declare
    a = b                           #assign
    b = c                           #assign
    lim = Add(lim, 1)               #assign
}

if Lt(lim, 12) {                    #if
    return lim                      #return
}
return -1                           #return
```

# Scope
&lt;if&gt; and &lt;for&gt; statements have their own scope; Please see this example:

```
x := 2
if 1 {
    x := 55
    Printf("%v ", x)
}
Printf("%v ", x)
```

The output of this program is "55 2";

# Expression
Expression in FuncClang is extremely simple; There are only four typies of expression: 

+ literal number: 123
+ literal string: "Hello World"
+ identification: x
+ function call : Println(123, "Hello", x)

Each expression will return a value which implement interface Var;

# All operators are function call
For simplification, there are no '+', '<', '==', '&&' and all other operators;

If you want, you can use default functions (like Add, Lt, Euqal, And) or write your own functions to replace these operators;

All default functions are put in <a href="https://github.com/qw4990/FuncLang/blob/master/default_funcs.go">default_funcs.go</a>;

Please read <a href="https://github.com/qw4990/FuncLang/blob/master/default_funcs.go">default_funcs.go</a> to see how to write your own functions;

# Demo
FuncLang can be easily nested in your Golang;

<b>You just implement Func in Golang and register it into FuncLang, then you can use it in FuncLang;</b>

There is a demo (or five demos actually):

```
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/qw4990/func_lang"
)

type fcVar struct {
	num float64
	str string
	err error
}

var (
	// zero number or empty string will be regared as False
	// please see IsTure() function in interface.go
	fcFalse = &fcVar{num: 0, str: ""}
	fcTrue  = &fcVar{num: 1, str: "1"}
)

func (fv *fcVar) IsNum() bool {
	return fv.num != 0
}

func (fv *fcVar) IsStr() bool {
	return fv.str != ""
}

func (fv *fcVar) Num() float64 {
	return fv.num
}

func (fv *fcVar) Str() string {
	return fv.str
}

func fcGt(vs ...func_lang.Var) func_lang.Var {
	if len(vs) != 2 { // if there is error, use panic to expose it
		panic("args error")
	}
	if (vs[0].IsNum() && vs[1].IsNum()) == false {
		panic("args error")
	}
	if vs[0].Num() > vs[1].Num() {
		return fcTrue
	}
	return fcFalse
}

func fcAtoi(vs ...func_lang.Var) func_lang.Var {
	if len(vs) != 1 {
		panic("args error")
	}
	if vs[0].IsNum() {
		return vs[0]
	}
	str := vs[0].Str()
	num, err := strconv.Atoi(str)
	if err != nil {
		panic("Atoi err: " + err.Error())
	}
	return &fcVar{num: float64(num)}
}

func fcRand(vs ...func_lang.Var) func_lang.Var {
	rand.Seed(time.Now().Unix())
	x := rand.Float64()
	return &fcVar{num: x}
}

var (
	func1 = `Println("Hello, World")`
	func2 = `
            for Lt(Rand(), 0.5) {
                Println("Try again!!!")
            }
    `
	func3 = `
            strx := "233"
            numx := Atoi(strx)
            if Gt(numx, 222) {
                Println("Are you happy?")
            }`
	func4 = `
            x := Rand()
            if Lt(x, 5) {
                return "Happy!!!"
            }
            return "Not Happy!!!"
    `
	defaultBody = `
            Println("Default function body; ", "xxxxx")
    `
)

func call(body string) {
	fc := func_lang.NewFunCaller()
	fc.RegisterFunc("Rand", fcRand)
	fc.RegisterFunc("Atoi", fcAtoi)
	fc.RegisterFunc("Gt", fcGt)

	reader := strings.NewReader(body)

	result, err := fc.Call(reader)
	if err != nil {
		panic(err)
	}
	if result != nil && result.IsStr() {
		fmt.Println("FUNC RET: ", result.Str())
	}
}

func main() {
	call(func1)
	call(func2)
	call(func3)
	call(func4)

	getFuncBodyFromNet := func() string { return "" }
	for i := 0; i < 10; i++ {
		body := getFuncBodyFromNet()
		if body == "" {
			body = defaultBody
		}
		call(body)
		time.Sleep(time.Second)
	}
}
```

Please read <a href="https://github.com/qw4990/FuncLang/blob/master/demos_test.go">demos_test.go</a> for more examples;