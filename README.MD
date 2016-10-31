# EZLang
EZLang is an easy language written by Golang and can be nested in your Golang code;

EZLang is easily to use and understand; The total lines of code of this project is less than 1000;

# Syntax
## Types
EZLang has only two types: string and number(float64);

You can use &lt;declare&gt; statement to declare variables, for example:

```
name := "zhangyuanjia"
age := 22

Println(name, "'s age is", age)
```

### Bool
Empty string and zero number will be regarded as False in EZLang, and all other values are True;

You can see function "IsTure" in interface.go;

Bool value will be used in &lt;if&gt; and &lt;for&gt; statements; please see below;

## Statements
EZLang programs consist of multiple statements;

There are only four kinds of statements in this language: &lt;if&gt;, &lt;for&gt;, &lt;declare&gt;, &lt;assign&gt;;

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
    Printf("Lim is less than 12;")  #expression
}
```

### Scope
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

## Expression
Expression in EZLang is extremely simple; There are only four typies of expression: 

+ literal number: 123
+ literal string: "Hello World"
+ identification: x
+ function call : Println(123, "Hello", x)

Each expression will return a value;

### All operators are function call
For simplification, there are no '+', '<', '==', '&&' and all other operators in EZLang;

If you want, you can use default functions (like Add, Lt, Euqal, And) or write your own functions to replace these operators;

All default functions are put in default_funcs.go; 

Please read default_funcs.go to see how to write your own functions;

# Use EZLang with Golang
EZLang can be easily nested in your Golang;

Please read demos_test.go for more examples;