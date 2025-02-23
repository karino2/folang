# 2. Using the frt package (explanation of pkg_all.foi)

GoEval, used in the previous [1. Getting Started](1_GettingStarted.md), is a special function that can be executed without using any external packages, but this is not something you would normally do.

Normally, you would use the standard Folang package called `frt`.

`frt` stands for Folang RunTime, and provides the basic functions required for most Folang programs.

Let's try using frt here.  Here, we will use frt's Println.

## Execute go mod init (if not done)

If you followed the previous [1. Getting Started](1_GettingStarted.md) you should have done this, but this time we will use frt, which is an external package from the Go language's perspective, so we will need to configure go mod init.

Let's do the following.

```
hello_fo $ pwd
~/helo_fo
hello_fo $ go mod init hello_fo
hello_fo $ go mod tidy
```

## Create hello_frt.fo file

Let's create a hello_frt.fo file that uses frt.

The contents are as follows.

```
// hello_frt.fo
package main
import frt

let main () =
  frt.Println "Hello World"

```

First, the import on the second line has been updated.

```
import frt
```

Standard packages can be specified without double quotes,
which is the same as the following import with double quotes:

```
import "github.com/karino2/folang/pkg/frt"
```

Basically, libraries intended to be used from Folang are placed under folang/pkg.

By the way, comments are line comments of `//` and block comments of `/* */`, just like Golang.

## Transpiling hello_frt.fo and pkg_all.foi file

Now, if you transpile the previous file, you will get the following error message.

```
./fc hello_frt.fo
transpile: hello_frt.fo
panic: Unknown var ref: frt.Println

goroutine 1 [running]:
github.com/karino2/folang/pkg/frt.Panic(...)
...
Long stack trace below
...
```

You may be surprised by the large number of error messages, but this is because `fc` is still under development.

The only meaningful message is the top one, "panic: Unknown var ref: frt.Println".

This means that the fc command does not know the information about the frt package.

That information is written in pkg_all.foi.
Specify this before hello_frt.fo.

```
$ ./fc pkg_all.foi hello_frt.fo
$ go fmt gen_hello_frt.go
```

This should generate a file like the one below.

```golang
// gen_hello_frt.go
package main

import "github.com/karino2/folang/pkg/frt"

func main() {
  frt.Println("Hello World")
}
```

Now, if you run go run just like you would normally run Go, you will be told:

```
$ go run gen_hello_frt.go
gen_hello_frt.go:3:8: no required module provides package github.com/karino2/folang/pkg/frt; to add it:
go get github.com/karino2/folang/pkg/frt
```

So run go get.

```
$ go get
go: added github.com/google/go-cmp v0.6.0
go: added github.com/karino2/folang/pkg/frt v0.0.0-20250220122800-4ff80daf0a9a
```

Then run go run to execute it.

```
$ go run gen_hello_frt.go
Hello World
```

To summarize the above

- Folang normally imports frt
- To import frt, put a file called pkg_all.foi before the .fo file in `fc` command
- To use standard packages in Folang (packages external to Go), you need to use go get etc. just like in normal Go.

## A brief explanation of calling frt.Println

This is the first time I've seen Folang-like code, so I'll explain the basics.

Reposting it, it looks like this:

```
let main () =
  frt.Println "Hello World"
```

The function call is shown below.

```
  frt.Println "Hello World"
```

Functions are executed by lining up arguments separated by spaces.
Unlike Golang, no parentheses are used.

frt.Println is a function that takes a string as an argument and does not return a result.
This is similar to Golang's fmt.Println.

## A little about the contents of pkg_all.foi

Now, the pkg_all.foi mentioned earlier is just a text file.
The frt interface is written at the beginning of this text file.
Let's extract that part.

```
package_info frt =
  let Println: string->()
  let Sprintf1<T>: string->T->string
  let Printf1<T>: string->T->()
  let Fst<T, U> : T*U->T
  let Snd<T, U> : T*U->U
  let Assert : bool->string->()
  let Panic : string->()
  let Panicf1<T>: string->T->()
  let Empty<T>: ()->T
```

And the only information needed for this program is Println.

Even if you write this in hello_frt.fo, it can actually be executed.

```
// hello_frt2.fo
package main
import frt

// Add these two lines
package_info frt =
  let Println: string->()

let main () =
  frt.Println "Hello Frt2"
```

This way, you can transpile without specifying pkg_all.foi.
Let's try it out.

```
$ ./fc hello_frt2.fo
$ go run gen_hello_frt2.go
Hello Frt2
```

The transpiler just assumes that there is a function of that type, resolves the type, and generates Go code.
It is the programmer's responsibility to ensure that the generated code actually runs correctly.

This is important when you write and run your own wrapper, so keep it in mind.

## Summary of the second session

- Import standard Folang packages without double quotes
- Import frt normally
- Specify the pkg_all.foi file at the beginning of the fc command
- Comments are `//` and `/* */`

## Next session: Slices, pipes, and maps

[3. Slices, pipes, and maps](3_SlicePipeMap.md)