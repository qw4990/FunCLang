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

This program prints all fibonacci numbers less than 1000:

```
a := 1
b := 1
for Lt(b, 1000) {
    Println(b)
    c := Add(a, b)
    a = b
    b = c
}
```

### Scope
&lt;if&gt; and &lt;for&gt; statements have their own scope; Please see this example:

```
x := 2
if 1 {
    x := 55
    Printf(x)
}
Printf(x)
```

The output of this program is "55 2";

## Expression
