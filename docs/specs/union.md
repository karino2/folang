# Implementing a Discriminated Union

For example, consider the following simple case.

```
type IntOrBool =
| I of int
| B of bool
```

This creates functions I and B, and the result is of type IntOrBool, which can be determined by pattern matching at runtime.

It seems that int and bool can sometimes have the same type (see the example of EquilateralTriangle and Square from Microsoft Learn above).

So they cannot be distinguished by a simple type assertion.

So, we will create structs IntOrBool_I and IntOrBool_B with IntOrBool as the interface.

Implementation as follows.

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

In this case, IntOrBool can distinguish between I and B at run time with type assertion, right?
Let's try it.

```fsharp
match iob with
| I ival -> printfn "i=%d" ival
| B bval -> printfn "b=%v" bval
```

This simple case seems like it can be achieved with a simple type assertion.

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

Of course, in reality, there may be more complicated patterns, so there is the question of whether they can be written with type switch, but perhaps we can write all of them with further conditions inside the case?
Well, I won't be using complicated patterns for a while, so I should make this simple case work first.

## Case without of

```
type AorB =
| A
| B
```

You can also do things like this. In this case, A is a variable, not a function (fsharp does not distinguish between functions without arguments and variables, but it is distinguished from functions with Unit arguments).

For now, I'll try making the Golang side var as follows.

```golang
var New_AorB_A AorB = AorB_A{}
```

It's strange to have New in the variable name, but I don't want to change the code too much depending on whether there is of or not, so I'll leave it like this.

This name won't come up in Folang anyway.

## Related links for F#

An explanation of Discriminated Union in the first place.

I'll post a link to a document explaining the implementation of F#, which I'm familiar with.

- [Discriminated Unions - F# for fun and profit](https://fsharpforfunandprofit.com/posts/discriminated-unions/) Explanation of F# as a function
- [Discriminated Unions - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/discriminated-unions)