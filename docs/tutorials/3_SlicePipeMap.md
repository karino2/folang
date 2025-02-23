# 3. Slices, pipes, and maps

Previous: [2. Using the frt package (explanation of pkg_all.foi)](2_UseFrtPackage.md)

Slices are the main data structure in Folang, as in Golang.

Folang slices are the same internally as Golang, but there are some differences in syntax.

Also, slices are usually handled more like functional languages in Folang,
and maps and pipelines are typical examples of this.

## How to create slices and let

Slices are created by separating them with semicolons as follows.

```fsharp
let main () =
  let s = [1; 2; 3]
  frt.Printf1 "%v\n" s
```

In addition to slices, there are several new elements.

1. Define variables with let
2. There is something called Printf1, which is like Printf in Golang with one additional argument
3. Slice literals are enclosed in brackets with a semicolon between them (not a comma!)

The generated code is as follows.

```golang
func main() {
  s := ([]int{1, 2, 3})
  frt.Printf1("%v\n", s)
}
```

`let s =` is transpiled to `s :=`.

`%v` and other expressions are the same as in normal Golang.

## Map, strings.Concat, and pipe operator

Next, let's look at a slightly more complicated example.

```
package main
import frt

// Add the following two lines
import slice
import strings

let main () =
  [1; 2; 3]
  |> slice.Map (frt.Sprintf1 "This is %d")
  |> strings.Concat ", "
  |> frt.Println

```

When executed, the following will be displayed.

```
This is 1, This is 2, This is 3
```

### Explanation for people who know F#

Anyone familiar with F# should be able to read this as is, but first, for those who know F#, I will explain some of the differences.

- Package names start with lowercase and function names start with uppercase (this is the Golang convention, unlike the .NET convention)
- slice.XXX has roughly the same API as List.XXX

If you understand this much, it should be enough for F# users.

If you think of Folang as a language that replaces the .NET-like parts of F# with Golang-like ones, you will understand about 90% of the code correctly.

### Pipe operator (for people who don't know F#)

Below is a simple explanation of the code for people who don't know F#.

However, in this tutorial, we will prioritize trying to run it without paying too much attention to the details, so we will keep the explanation short.

First of all, the following two lines are quite difficult to understand.

```fsharp
  [1; 2; 3]
  |> slice.Map (frt.Sprintf1 "This is %d")
```

First of all, `|>` is called the pipe operator, and is a binary operator like `+`.
What it does is call the element on the left as an argument to the function on the right.

The pipe operator is very difficult to understand for first-timers, but you will quickly get used to it as you use it, so it is recommended that you move on after you have a vague understanding of the explanation.

The pipe operator passes the element on the left to the function on the right, so two conditions are necessary:

1. The left of the pipe operator is an element
2. The right of the pipe operator is a function with one argument

Please read the explanation below while keeping these two points in mind.

### Function composition by partial application

Now, with the following code

```fsharp
  |> slice.Map (frt.Sprintf1 "This is %d")
```

I would like to explain slice.Map after the pipe operator, but before that, the difficult part is `(frt.Sprintf1 "This is %d")`.

Sprintf1 is like Sprintf that takes one additional argument, as you can guess from the name.

Therefore, it is originally used as follows.

```fsharp
frt.Sprintf1 "This is %d" 123 //<- Note the last 123.
```

This is how it is normally used, with the argument 123 after the string.

So what happens if there is no 123, that is, if there is one missing argument?

```fsharp
let f1 = frt.Sprintf1 "This is %d"
f1 123
```

This is called partial application, and it becomes a function with one remaining argument (in this case, it is interpreted as a function of `func[T any](arg T) string`).
When you actually transpile this code, it becomes the following code.

```golang
f1 := (func (_r0 int) string{
        return frt.Sprintf1("This is %d", _r0)
      })
f1(123)
```

The type of `_r0` is int, but it would take a long time to explain, so just assume that it is.

Side note: `fc` doesn't look at `%d` properly (it's not that smart at the moment).

Anyway, in Folang, when a function does not have enough arguments, a new function object is generated as a func with the missing arguments.
This is called partial application.

The important thing here is that `frt.Sprinf1 "This is %d"` is a function called `func(int)string`.

### slice.Map

The slice package is also used in almost all Folang programs.

`slice.Map` is particularly common.

slice.Map takes two arguments as follows:

`slice.Map fn slice1`

The second argument is a slice, and fn is executed for each element in turn, and the results are returned as a slice.

```fsharp

let add10 a =
   a+10

let main () =
  let res = slice.Map add10 [1; 2; 3]
  frt.Printf1 "%v\n" res

```

This will display the following.

```
[11 12 13]
```

This is Map.

### Rewrite with pipe operator

Now, the following line:

```fsharp
let res = slice.Map add10 [1; 2; 3]
```

You can use the pipe operator to bring the last element to the front.

```fsharp
let res = [1; 2; 3] |> slice.Map add10
```

It's hard to see the whole thing, but if you take out just the following part,

```fsharp
slice.Map add10 [1; 2; 3]
```

You can move this last element to the left and write it as follows.

```fsharp
[1; 2; 3] |> slice.Map add10
```

The result is the same.

With this in mind, the following code should be easy to understand, right?

```fsharp
  [1; 2; 3]
  |> slice.Map (frt.Sprintf1 "This is %d")
  |> strings.Concat ", "
  |> frt.Println
```

Note that strings.Concat is also a bit special.
This function takes a slice of strings and returns a string that is concatenated with arguments in between.
In other words, it takes `[]string` and returns the concatenated `string`.

Processing various things with this Map and calling strings.Concat at the end is the basic idiom for string processing in Folang.

## Summary

- Slice literals are separated by brackets and semicolons
- Pipe operators and slice.Map can be used
- Function calls with insufficient arguments become partial applications
- Creating a slice of strings with Map and concatenating them with strings.Concat is a common idiom

## Next time: Calling a wrapper written in Go

[4. Calling a wrapper written in Go](4_CallingGoWrapper.md)