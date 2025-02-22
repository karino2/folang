# What is Folang?

Folang is a transpiler to Golang, a functional language similar to F# and ReScript.

F# is OCaml for .NET, ReScript is OCaml for JavaScript.
If so, Folang aims to be a OCaml for Golang.
However, the syntax is more heavily influenced by F# than OCaml.

If you look at the [sample page](../samples/README.md), you can get a sense of the atmosphere.

Folang is self-hosted on Folang. You can try it only with Golang. It is a big advantage that there is no need for an OCaml or .NET environment.

Also, after transpilation, it will be a normal Go file.
Simply put the file on github and you can get into the normal Go module system.

Let's talk about each of them in a little more detail below.

## The generated data type or function type is natural as Go.

The following Folang code is

```fsharp
let a = 3
```

It will be transpiled into the following Go code:

```golang
a := 3
```

The following record type is

```fsharp
type Rec = {A: int; B: string}
```


It expands to the following struct:

```golang
type Rec struct {
  A int
  B string
}
```

The following function is

```fsharp
let add3 (a:int) =
  a + 3
```

It will be converted to the following Go code:

```golang
func add3(a int) int {
  return a+3
}
```

Though expression "if else" becomes an ugly code wrapped in a function because it becomes an expression,
basically you can call the generated function or use the data type as natural Go entities.
Calling function and using type is just work.

I wouldn't say that the generated code is easy to read, but At least it's debuggable in the debugger.

## Once generated, it's just a Go file

Since it is a transpiler, the generated result is a Go source file.

The command to transpile is called `fc`, but once transpiled, `fc` is not required thereafter.
The runtime is also standard Go module in `github.com/karino2/folang/pkg/frt`, which is a standard regular Go package.
It's no different from using any other package.

As a result, after generating it, it can be used as a Go module.
If you put it on github, you can also `go get` and `go install`.

What kind of package will it be and what kind of mod will it be?
It is a problem on the side of the Go language, and Folang is not relevant.
And Go is a very good language with this package and deployment mechanism.

In addition, deployment of the result is easy, just a single binary.
It also starts up quickly. GO Single Binary is really awesome!

All of these characteristics of the Go language, which makes it a great tool for writing command line tools, are inherited.

However, in terms of execution speed, the code to be generated is code with many functions unnecessarily.
I think it's much slower than writing Golang by hand.
I think Folang is a language that doesn't perform very well (although I haven't measured it).

## You don't need an environment other than Go.

Folang's transpiler is written in Folang.
You only need a Go environment to run the generated Go code. You don't need to set up an environment such as OCaml or Haskell.
I think that's a big advantage.

Go develop environment is well maintained in many environments and does not clog up.
The generated binaries also work trouble-free in various environments.

And since many programmers set up the Go development environment anyway, in many cases, you already have it at hand from the beginning.

If you make it on one machine, you can easily deploy it to other machines, just call Go install to each machine.
I think it's a big advantage to make a chore tool with Golang.

## Intended for use as a small-scale tool.

The main target of Folang is a command line tool with a scale of 100 to 5000 lines.
At the time of this writing, the largest Folang program is [fc](../fc), which is the Folang transpiler itself.
It currently consists of 3612 lines of fo code and 529 lines of Go code. Even if it is large, it is assumed to be about this scale.
Those who want to create something on a larger scale are better off choosing other languages.

We do not think that we should handle too large data efficiently.
I don't really care about runtime performance either.
These characteristics also do not seem to be suitable for creating too large tools and web services.

However, it is important to start a small program quickly.
I'm not doing anything for that in particular, but Go is great because it gets up to speed quickly.

Personally, I made it to replace everyday tools that I usually write in F# or Python.

## Why OCaml-like languages?

Most of the benefits listed so far are just saying that the Go language is great.
There are no particular benefits of Folang.
In fact, there are many people who prefer Go because of their preference.

In the following, I will write why I wanted an OCaml like language.
It's just my preference.

The Go language emphasizes that it works as you see it.
Don't do things that cost you money out of sight.
I think that the language design is done.
This is also a very good feature.
I think that the fact that it works so well in system programs is the reason why the Go language has become so widespread.

However, if you want to write a small tool for your own use, I don't really care if it costs me out of sight.
Small, fast-paced single binaries and mod systems, good packages, well-maintained dev environment, I choose Golang for small-mid chore.
It doesn't have to be this low level language at all.

I'm happier to be able to write concisely than that.
In particular, it is nice to write small tools that can be written concisely with type inference, pipe operators, and partial applications.

If it's a language that can be written slowly and shortly and concisely, there are various candidates other than Golang.
However, getting slightly larger code size, a few thousand lines, many of these languages become a little tough.

For me, I can start small and concise.
When things get a little more complicated, I want to use the powerful type modeling capabilities of Discriminated Union.
I want a concise language of high level, that could be slow.

In addition, there are languages that can write code to some extent with only basic functions.
It is desirable from the standpoint of implementing the transpiler individually.
Languages like kotlin that can only be used conveniently after having complex things such as inline returns, it's tough to implement.
It is desirable to have a language that you can manage with only the minimum number of functions to start using.

Also, OCaml like language has a track record as a transpiler like ReScript.
I thought it would be a good reference when I was worried about how to realize it.

Finally, I was familiar with F# and liked it.
I want something similar.

## Transpiler open to Go

I mentioned earlier that the generated Go code is plain as a type or function.

In addition, Folang does not aim to include everything necessary as a language from the beginning.
The missing part is supposed to be written in Golang side.
Instead, the Folang language specification is kept simple.
It is designed to be easy to interact with Golang's code.

Especially in the case of more than 1000 lines of code, it is recommended to divide labor so that the parts that are suitable for Go are written in Go.

For example, Folang's transpiler, `fc`, at the time of this writing, it consists of 3612 lines of fo code and 529 lines of Go code.
This is a typical Folang configuration.
Instead of forcibly adding syntax such as loops, destructive assignments, and pointers, the idea is to write such things in Go.

It's ideomatic, though not required, to create a file called wrapper.go and we implement functions for Folang with a free-standing function and an argument adjustment for the pipeline in this file,
then writing interface information in `.fo` file and call them.

In Folang, onion architecture is often used.
Write the inside side of the onion architecture in Folang, writing the outside in Golang.
It's not suitable to just call the Go package as it is.
Write a wrapper for the problem and hand-write the interface information for it, It's designed to make it easy.

Also, because it is assumed that the division of labor with Go will be done from the beginning, even if there are unimplemented or missing features in Folang, it can be implemented on the Go side at any time.
Therefore, you can use it with the existing functions without waiting for the completion of Folang transpiler.

In order to aim for a language that is open to this Go language, we are not going to make the backend of the existing language for Golang.
Instead, I decided to create my own language for Go.

## Tools for generating the Go source

I don't intend to do everything entirely in our own language.
Instead, Folang transpiler is just convenient tool for generating Go source.

In fact, all you need to do is generate the desired code. Some can lie about the interface definition.
You can also use features such as GoEval to create a hole for the Go language.

For example, you can lie that Sprintf is a two-argument function like this.

```
import "fmt"

package_info fmt =
  let Sprintf<T>: string->T->string

let Sprintf1 fmt arg =
  fmt.Sprintf fmt arg
```

You can lie with a three-argument function as shown below.

```
package_info fmt =
  let Sprintf<T, U>: string->T->U->string

let Sprintf2 msg arg1 arg2 =
  fmt.Sprintf msg arg1 arg2
```

If the type is consistent in Folang, and if the generated Golang is legitimate code in Golang world, that's fine.

Rather than thinking of it as a closed world, something that generates Go code as an extension of the editor.
I think it's closer to the reality if you think about it that way.