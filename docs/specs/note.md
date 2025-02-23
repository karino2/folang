# Notes on specifications

As specifications are decided, they are added to the top.

The order is from the bottom up, so it's a bit tricky.

The writing style is not yet consistent.

I'd like to organize it once it's more solidified.

The background to these decisions is in [discussion.md](discussion.md).

## Slice and tuple type precedence

The precedence of slices and tuples is a Folang-specific problem that occurs because slice syntax is Go.

It is decided as follows.

- `[]T*U` is parsed as `[](T*U)`.
- `T*[]U` is legal

The above is a bit tricky, but I use it often so I didn't want to put parentheses.

## Slice expression

Semicolon delimited with brackets, no type name. (Note that the separator is not a comma)

```
let ika () =
  [1; 2; 3]
```

## Equal comparison

Equal is done with `=`, and not equal is done with `<>`.

Equal also compares the contents of slices (uses go-comp's Equal [cmp package - github.com/google/go-cmp/cmp - Go Packages](https://pkg.go.dev/github.com/google/go-cmp/cmp))

## Import system library

Imports that are not enclosed in double quotes are considered to be system imports, and internally the path to Folang/pkg is considered to be prefixed.

Specifically, the following two have the same meaning. (The Go code generated will be on the second line in both cases)

```
import frt
import "github.com/karino2/folang/pkg/frt"
```

## Generics syntax (type parameters)

Go uses square brackets, but FSharp uses square brackets.
Folang use square bracket.

```
package_info slice =
  let New<T>: ()->[]T

slice.New<int> ()
```


## External type information

A language element for accessing external packages, etc. It is similar to FSharp's signature file.

The file extension is written in .foi (however, it can also be written in the .fo file, it just won't generate the corresponding Go file if it is written in .foi).

It is defined using package_info, and the function is written as a colon with let, imitating ReScript.

```
package_info slice =
   let Length<T>: []T -> int
   let Take<T> : int->[]T->[]T 
```

Write the type as follows with type.

```
package_info buf =
    type Buffer
    let New: ()->Buffer
    let Write: Buffer->()
    let String: ()->string
```

The file extension is foi. (However, you can also write it in a .fo file. If you use foi, the corresponding Go file will not be generated.)

### Add external type information to your own namespace

By naming the package with an underscore, it will be added to the current namespace without a prefix.

```
package_info _ =
   type wrappedType
   let New: ()->wrappedType
   let doWork: wrappedType -> ()
```

This is used when you create a wrapper for Folang in the same package using wrapper.go, etc., and refer to that information.
If it is a method, wrap it in a normal function and use it.

When doing this, it is a good idea to decide the order of arguments in F# terms, taking currying into consideration.

## Comments

Comments support both C-style and C++ one-line comments, just like Golang.

``
package main

/*
This is a comment
*/

let ika () =
  123 // This is also a comment.
``

## GoEval

Write Go code as a string, similar to an inline assembler.

The following code is expanded:

```
package main
import "fmt"

let main () =
  GoEval "fmt.Println(\"Hello World\")"
```

Expands to the following code:

```
package main
import "fmt"

func main() {
  fmt.Println("Hello World")
}
```

### Specifying the return type

By default, GoEval is considered to be Unit. If you want to specify the return type, specify it with the type parameter.

```
// This s is a string
let s = GoEval<string> "fmt.Sprintf(\"hoge %d\", 123)"
```

The expanded code is as follows.

```
s := fmt.Sprintf("hoge %d", 123)
```

There is no type specification on the Go side, so you must specify it to be the same type as s on Folang.

### Using arguments

The argument identifier is carried over as is in Go, so you can write it as follows. Note that a is used.

```
let ika (a:int) =
  GoEval<string> "fmt.Sprintf(\"hoge %d\", a)"
```

This is expanded as follows.

```
func ika(a int) string {
   return fmt.Sprintf("hoge %d", a)
}
```

Use arguments while keeping in mind the Go code generated below. This is the same as inline assembler.

Functions that are not supported by Folang can be wrapped with GoEval and used quite easily.

## Implementing Discriminated Union

This is long, so go to [union.md](union.md).

## Function definition

The basic function definition is as follows.

```
let ika (a:string) (b:string) =
a+b
```

This is expanded as follows.

```golang
func ika(a string, b string) string {
  return a+b
}
```

Note that no arguments are defined as one argument unit.
GoEval is an expression that passes the argument code directly to Go, and the type is specified by the type parameter, but if no type is specified, it is Unit.

Using this, the following code will be

```
package main
import "fmt"

let main () =
  GoEval "fmt.Println(\"Hello World\")"
```

The following code will be generated.

```golang
import "fmt"

func main() {
  fmt.Println("Hello World")
}
```

Methods are not supported (wrap them manually)

### Type inference

Arguments are inferred, and those that cannot be determined become generics.
The following Go code

```fsharp
let add10 a =
  a + 10
```

The following code will result.

```golang
func add10(a int) int {
  return a + 10
}
```

If the type is unknown, it will be a generics type parameter.

```fsharp
let secondHead s =
  slice.Item 2 s
```

In this case, it is certain that s is a slice, but the elements are not, so the following code will be generated.

```golang
func sceondHead[T0 any](s []T0) T0 {
  slice.Item(2, s)
}
```

Currently, type constraints are not supported, so everything will be any.

Therefore, the following code will result in a compilation error in Go.

```
let ika a b =
  a+b
```

If one of the two is specified as shown below, it will be resolved as both are the same type, so it will work.

```
let ika (a:int) b =
  a+b
```

### Function call

Function calls are F# style and partial application occurs when there are not enough arguments.

Let's start with a basic call. Pay attention to the call to the hello function below.

```
import "fmt"

let hello (msg: string) =
  GoEval "fmt.Println(msg)"

let main() =
  hello "hoge"
```

Folang calls functions without parentheses. The part with `hello "hoge"`.

This is expanded to the following code.

```golang
import "fmt"

func hello(msg string) {
  fmt.Println(msg)
}

func main() {
  hello("hoge")
}
```

### Partial application of function calls

Partial application with multiple arguments is as follows.

```
let hello (msg: string) (num: int) =
  GoEval "fmt.Printf(msg, num)"

let main () =
  let temp = hello "hoge%d"
  temp 123
```

In the following line, only one argument is passed to hello, which has two arguments.

```
let temp = hello "hoge%d"
```

The generated code is as follows.

```
temp := func(num int) { hello("hello%d", num) }
```