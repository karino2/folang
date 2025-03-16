# Notes from when considering specifications

Collect notes from when considering specifications.
Write the finalized specifications in [note.md](note.md).

Since these are just notes, don't worry too much about the format.

## Separator for slice expressions

Commas are used in ReScript, semicolons in FSharp.

- [Array & List - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/array-and-list)
   - [Record - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/record) Records are commas
   - [Tuple - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/tuple) Tuples are commas in parentheses
- FSharp is semicolon [Lists - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/lists)
   - Records are semicolons
   - Tuples are commas in parentheses [Tuples - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/tuples)

Since I'm using semicolons to separate records, I've decided to standardize on semicolons.

When I first used F# in a while, I tend to use commas, but when the record field is a tuple, it's still confusing, so I think semicolons are better.

I decided to use semicolons to separate slices.

## Considering the specifications for function calls

For now, I want to think about the specifications with something like hello world in mind.

First of all, I want golang function calls to be like F# dotnet function calls.

```
import "fmt"

let main () =
   fmt.Println("hoge")
```

I want this to be expanded as follows.

```golang
import "fmt"

func main() {
   fmt.Println("hoge")
}
```

Next, consider the function definition in folang.

Let's start without type inference. Then the function definition is as follows.

```
import "fmt"

let hello (msg: string) =
    fmt.Println(msg)

let main() =
    hello "hoge"
```

For folang, functions can be called without parentheses, and partially applied like FSharp.

I want to distinguish between calling with parentheses and calling directly.

I'm not sure what kind of code this will be expanded into at this point.

```golang
import "fmt"

func hello(msg string) {
    fmt.Println(msg)
}

func main() {
   hello("hoge")
}
```

Whether this hello("hoge") will be expanded or partially applied can probably be determined at compile time.

If there are multiple arguments, it will look like this.

```
let hello (msg: string) (num: int) = 
   fmt.Printf(msg, num)

let main () =
   let temp = hello "hoge%d"
   temp 123
```

If it is partially applied, will it look like this?

```golang
func hello(msg string, num int) {
    fmt.Printf(msg, num)
}

func main() {
   temp := func(num int) { hello("hello%d", num) }
   temp(123)
}
```

Let's start by making it possible to generate something like this.

No, wouldn't it be better to define functions in a more golang-like way?

```
import "fmt"

func hello (msg string, num int) = 
   fmt.Printf(msg, num)

func main () =
   let temp = hello "hoge%d"
   temp 123
```

No, it's not good that the syntax for function definition is so cumbersome.

## Generics syntax (type parameters)

Go uses square brackets, but FSharp uses square brackets.

I would like to align it with golang as much as possible, but FSharp's syntax makes the distinction between variables and functions unclear,
so it becomes complicated with index access.

```
let Length[T any] (args: []T) =
   ...

Length[int] listOfList[3]
```

This is still tough. I'll use square brackets.

```
let Length<T any> (args: []T) =
   ...

Length<int> listOfList[3]
```

This makes parsing quite tricky, but there's no other way.

## External type information

I want to think about package access soon. I can create and add type information files manually, but what kind of syntax should I use?

### Related information

FSharp's signature files are a similar story.

- [Signature files - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/signature-files)

Functions become vals. Hmm. I don't really like using vals just for this. Module and namespace are the same as normal definitions.

In ReScript, they are defined with let.

- [Module - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/module#signatures)

But the syntax becomes a colon. It's weird, but if you think of variables as being defined, let seems correct.
The fact that the module type comes at the beginning is different from a normal module.

In Borgo, function definitions are keywords, so they look like normal function definitions.

- [Borgo Programming Language](https://borgo-lang.github.io/#package-definitions)

There is an advantage in being able to use the same parser, but the language syntax is so different that it's not much of a reference.

### External type information in folang

So, what should we do with folang?

Should we call it module? No, maybe it's better to change it.
In golang, it's package.

Let's call it package_info.
And let's imitate ReScript and use let with colon.

```
package_info slice =
   let Length<T>: []T -> int
   let Take<T> : int->[]T->[]T 
```

Slice requires generics right from the start... Length is fine with any, but Take is necessary. There's no choice, let's just give up and deal with it.
It's fine to specify it explicitly when calling for a while.

It's hard to decide whether to use `<T>` or `[T]` for the type parameter, but since index access is in square brackets, I'll align it with FSharp.

Let's write Buffer as is.

```
package_info buf =
    type Buffer
    let New: ()->Buffer
    let Write: Buffer->()
    let String: ()->string
```

This one seems better. There's nothing to write to the right of type. There's no need to support defining type aliases inside.

I'll make the New function create one since constructors will be wrapped anyway.

I think I'll keep the file extension as foi for now.

## Comments

Comments are supported in both C style and C++ one-line comments, just like golang. (However, only C style is currently implemented.)

## Implementation policy for discriminated union

First, here's a link to the implementation documentation for F#, which I'm familiar with.

- [Discriminated Unions - F# for fun and profit](https://fsharpforfunandprofit.com/posts/discriminated-unions/) Explanation as a function of F#
- [Discriminated Unions - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/discriminated-unions)

For example, consider the following simple case.

```
type IntOrBool =
| I of int
| B of bool
```

This creates functions I and B, and the result is of type IntOrBool, which can be determined by pattern matching at runtime.

It seems that int and bool may have the same type (see the example of EquilateralTriangle and Square on Microsoft Learn above).

So they cannot be distinguished by a simple type assertion.

How about the following implementation?

```golang
type IntOrBool interface {
  IntOrBool_Union()
}

func (IntOrBool_I) IntOrBool_Union(){}
func (IntOrBool_B) IntOrBool_Union(){}

type IntOrBool_I struct {
   Value int
}

type IntOrBool_B struct {
   Value bool
}

func New_IntOrBool_I(v int) IntOrBool { return IntOrBool_I{v} }
func New_IntOrBool_B(v bool) IntOrBool { return IntOrBool_B{v} }
```

It seems fine to map I and B to NewXXX function calls.

At first I used pointers, but it was difficult to distinguish between interface and struct internally in cases where the definition comes later, so I unified them all into entities.

In this case, IntOrBool can distinguish between I and B at runtime using type assertion, right?
Let's try it.

```fsharp
match iob with
| I ival -> printfn "i=%d" ival
| B bval -> printfn "b=%v" bval
```

This simple case seems like it can be achieved with a simple type assert.

```golang
switch iob.(type) {
case IntOrBool_I:
   ival := iob.Value
   fmt.Printf("i=%d", ival)
case IntOrBool_B:
   bval := iob.Value
   fmt.Printf("b=%v", bval)
}
```

Of course, there are more complicated patterns in reality, so the question is whether they can be written with type switch, but I think you can probably write all of them with further conditions inside the case.
Well, I won't be using complicated patterns for a while, so I should make this simple case work first.

### Case without of

```
type AorB =
| A
| B
```

You can also do things like this. In this case, A is a variable, not a function (fsharp does not distinguish between functions without arguments and variables, but it is distinguished from functions with Unit arguments).

For now, let's try making the golang side var as follows.

```golang
var New_AorB_A AorB = AorB_A{}
```

It's strange to have New in the variable name, but I don't want to change the code too much depending on whether there is of or not, so I'll leave it like this.
This name won't come up in flang anyway.

### Let's take a look at the altJS implementation of Union.

- [Fable Features](https://fable.io/docs/typescript/features.html#f-unions) Fable's Union implementation
- [Pattern Matching / Destructuring - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/pattern-matching-destructuring) ReScript's Union implementation and payload are useful references.

In both fable and ReScript, you first prepare a tag to represent each case and then enter a string or index into it.
In JS, type information is lost at runtime, so this seems necessary, but isn't it unnecessary in go?

The language Borgo has something similar to Rust's enum, so I'll take a look at it. [borgo/compiler/test/snapshot/codegen-emit/enums.exp at main · borgo-lang/borgo](https://github.com/borgo-lang/borgo/blob/main/compiler/test/snapshot/codegen-emit/enums.exp)

Tags are used.

## String literals

F# has three double quotes, but golang has backquotes.

I also need interpolation.

[Interpolated strings - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/interpolated-strings)

I think I'll implement backquotes and dollar prefixes for now.

```
let a = `This is
Multiline
string`

let b = $"String {a} interpolation"

let c = $`This
is
also {a}
interpolation. {{}} for brace pair.`
```

Are these two enough?

## Union generics

I wonder if generics can be used in Golang interfaces? I looked it up but couldn't really find out, but when I asked chatGPT, they gave me the code and it worked.

```golang
package main

import "fmt"

// Interface with type parameter T
type Printer[T any] interface {
	Print(value T)
}

// Printer implementation of int type
type IntPrinter struct{}

func (p IntPrinter) Print(value int) {
	fmt.Println("Printing int:", value)
}

// Printer implementation of string type
type StringPrinter struct{}

func (p StringPrinter) Print(value string) {
	fmt.Println("Printing string:", value)
}

func main() {
	var intPrinter Printer[int] = IntPrinter{}
	intPrinter.Print(42)

	var stringPrinter Printer[string] = StringPrinter{}
	stringPrinter.Print("Hello, World!")
}
```

If this works, it's not that difficult, right?

I googled to see if there was an implementation of Optional somewhere, and found the following.

[Generic Go Optionals · Preslav Rachev](https://preslav.me/2021/11/18/generic-golang-optionals/)

It seems like it might be possible to support generics for records and unions.

I want to make something like this.

```
type Result<T> =
| Success of T
| Failure of string
```

This seems like the following would be fine as Go code.

```golang
type Result[T any] interface {
   Result_Union()
}

func (Result_Success[T]) Result_Union(){}
func (Result_Failure[T]) Result_Union(){}

type Result_Success[T any] struct {
  Value T
}

type Result_Failure[T any] struct {
  Value string
}

func New_Result_Success[T any](v T) Result[T] { return Result_Success[T]{v} }
func New_Result_Failure[T any](v string) Result[T] { return Result_Failure[T]{v} }
```

I was able to confirm that it works.

But type inference on the Folang side is not easy.

Looking at the following example from [Understanding Parser Combinators - F# for fun and profit](https://fsharpforfunandprofit.com/posts/understanding-parser-combinators/),

```fsharp
type ParseResult<'a> =
  | Success of 'a
  | Failure of string

let pchar (charToMatch,str) =
  if String.IsNullOrEmpty(str) then
    Failure "No more input"
  else
    let first = str.[0]
    if first = charToMatch then
      let remaining = str.[1..]
      Success (charToMatch,remaining)
    else
      let msg = sprintf "Expecting '%c'. Got '%c'" charToMatch first
      Failure msg
```

The type parameter of this Failure is only determined by Success. No, I can just assign type variables to everything separately and unify them with transitive law.

I'll also post the original Result type.

- [Result<'T, 'TError> (FSharp.Core) - FSharp.Core](https://fsharp.github.io/fsharp-core-docs/reference/fsharp-core-fsharpresult-2.html)
- [Result (FSharp.Core) - FSharp.Core](https://fsharp.github.io/fsharp-core-docs/reference/fsharp-core-resultmodule.html)
- [Results - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/results)

### This doesn't work for cases where there are no values

When I tried to implement it, it didn't work for cases where there are no values. [folang/docs/specs/union_ja.md at main · karino2/folang](https://github.com/karino2/folang/blob/main/docs/specs/union_ja.md)

Originally, for the following Union,

```fsharp
type AorB =
| A
| B
```

The following Go code was generated.

```golang
var New_AorB_A AorB = AorB_A{}
```

But this doesn't allow you to specify T.

You can't create variables like this.

```golang
var New_AorB_A AorB[T] = AorB_A[T}{}
```

I guess the only way to make the case where there is no value a function is to use it. It's mostly fine as long as it's like this.

```golang
func New_AorB_A[T any]() AorB[T] { return AorB_A[T]{} }
```

Of course, in Folang, you have to specify it explicitly, but

```fsharp
AorB_A<int> ()
```

How does it work in F#?

```
> type AorB<'t> =
- | A
- | B of 't
-
- ;;
type AorB<'t> =
  | A
  | B of 't

> A ;;
val it: AorB<'a>

> B 123 ;;
val it: AorB<int> = B 123
```

Hmm, A is a generics type variable. This probably can't be achieved in golang. What should I do?

You can also use the same variable for arguments of different types in ReScript.

```rescript

type result<'a> =
  | Ok('a)
  | Failure
  | Other


module App = {
  let iToS = (i) => {
    switch(i) {
      | Ok(arg) => Int.toString(arg)
      | Failure => "int fail"
      | Other => "int other"
    }
  }
  
  let sToS = (s) => {
    switch(s) {
      | Ok(arg) => arg
      | Failure => "s fail"
      | Other => "s other"
    }
  }
    
  let make = (cond) => {
    let f = Failure
    let o = Other
    let a = if cond { Ok(123) } else { f }
    let b = if cond { Ok("abc") } else { f }
    iToS(a) ++ sToS(b) ++ iToS(o) ++ sToS(o)
  }
}
```

The type is determined when the variable is referenced, and I feel that it would be fine to just put the same value in at runtime and cast it.

I feel that in this case, the type is determined by the reference, not the variable definition.

### In Golang, if there is a type parameter, it is a function.

I thought that it would not be good for a transpiler to include too many concepts that do not exist in Golang, so in the following case,

```
type AorB<T> =
 | A
 | B of T
```

When creating A, assume that there is an argument of `()`.

```
let a = A<int> ()
```

If it is determined by inference, int is not necessary, but it is a function call anyway.

This means that once a variable has been determined, it cannot be of a different type, but that is the specification.

If there is no type parameter, it becomes a variable, so there is no consistency. I think I should have made them all functions, but I don't feel like fixing it now, so I'll just make them special for generics.