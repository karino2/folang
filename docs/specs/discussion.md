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