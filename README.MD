# FunCLang
FunCLang (Function Call Language) is an easy language written by Golang and can be nested in your Golang code;

FunCLang is easily to use and understand; The total lines of code of this project is less than 1000;

You should regard FunCLang codes as the body of your Golang function, which can be modified on runtime ;)

# Types and Variables
FunCLang has only two types: string and number(float64);

You can use &lt;declare&gt; statement to declare variables, for example:

```
name := "zhangyuanjia"
age := 22

Println(name, "'s age is", age)
```

# Interact with Golang
Each variable in FunCLang is represented by an interface in Golang:

```
type Var interface {
	IsNum() bool
	Num() float64
	IsStr() bool
	Str() string
	IsTrue() bool
}
```

And each function in FunCLang is a function whose type is:

```
type Func func(args ...Var) Var
```

So the Println in above code may be implemented by Golang:

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

Please see more example in demos_test.go;


# Statements
FuncClang programs consist of multiple statements;

There are only four kinds of statements in this language: &lt;if&gt;, &lt;for&gt;, &lt;declare&gt;, &lt;assign&gt;, &lt;return&gt;;

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

Each expression will return a value who implement interface Var;

# All operators are function call
For simplification, there are no '+', '<', '==', '&&' and all other operators;

If you want, you can use default functions (like Add, Lt, Euqal, And) or write your own functions to replace these operators;

All default functions are put in <a href="https://github.com/qw4990/FuncLang/blob/master/default_funcs.go">default_funcs.go</a>;

Please read <a href="https://github.com/qw4990/FuncLang/blob/master/default_funcs.go">default_funcs.go</a> to see how to write your own functions;

# Use FuncLang with Golang
FuncLang can be easily nested in your Golang;

<b>You just implement Func and register your function into FuncLang interpreter, then you can use your function in FuncLang;</b>

Please read <a href="https://github.com/qw4990/FuncLang/blob/master/demos_test.go">demos_test.go</a> for more examples;