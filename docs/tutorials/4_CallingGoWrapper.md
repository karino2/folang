# 4. Calling a wrapper written in Go

Last time: [3. Slices, pipes, and maps](3_SlicePipeMap.md)

This time, we will look at an example of calling your own code written in Go from Folang.

## Folang is basically developed in a division of labor with Go

Folang does not aim to include all functions.

For example, it does not have loops, destructive assignment to variables, method calls, pointers, etc.

The basic style is to write the parts that use these functions in Golang,
and call them from Folang.

For Folang, collaboration with Golang is not the exception, but the norm.
Folang programs of a certain level or higher will likely include parts written in Go.

Here, we will look at how to collaborate with Golang.

## Example of calling a Go package directly

As mentioned briefly in the second article [2. Using the frt package (explanation of pkg_all.foi)] (2_UseFrtPackage.md),
you can inform `fc` of the existence of functions defined outside by writing something called `package_info`.

In general, Folang can be used as is for Go libraries with many freestanding functions,
but object-like libraries need to be wrapped before use.

First, let's try using the Golang `filepath` package, which is a commonly used package.

Here, we'll try using

- `filepath.Dir`
- `filepath.Join`
- `filepath.Base`

```
package main

import frt
import "path/filepath"

package_info filepath =
  let Dir: string->string
  // The actual type is vararg, but it is used as two arguments
  let Join: string->string->string
  let Base: string->string

let main () =
  let test = "/home/karino2/src/folang/README.md"
  filepath.Dir test |> frt.Println
  filepath.Base test |> frt.Println
  let dir = filepath.Dir test
  filepath.Join dir "hello.txt" |> frt.Println

```

This will print the following:

```
/home/karino2/src/folang
README.md
/home/karino2/src/folang/hello.txt
```

### How to write package_info

The first part of package_info is as follows.

```
package_info filepath =
  let Dir: string->string
```

First, write the package name after the `package_info` keyword.
If you want the generated Go code to be, for example, `filepath.Dir`, write `filepath`.

The next line is the function signature.
Write the name after `let`, then write the function type after `:`.

The function type is connected to the argument types with `->`, and finally the return type with `->`.

For example, the following type in Golang is written as:

```golang
func add(a int, b int) int {...}
```

It will look like this.

```
int->int->int
```

The distinction between the argument and return types is difficult to understand, so it takes some getting used to,
but once you get used to it, it is a method of notation that works well with partial application.

To give another example, the following function is

```golang
func toStr(str string, a int) string {...}
```

It will look like this.

```
string->int->string
```

When reading, interpret the part after the last arrow as the return type,
and read the rest as being separated by `,`, and you can translate it into a Golang function type.

## Write a function in wrapper.go and call it

Calling a Golang package directly is convenient when you want to use a simple function,
but it can be a problem in many cases.
Folang does not support many features, such as pointers, method calls, variable length arguments, and multi-value returns.

If these are used in package, you cannot use it as is.

Even if you could use it, Golang functions usually do not have arguments in an order that assumes currying and pipelines.
(For more information, see Designing functions for partial application in [Partial application - F# for fun and profit](https://fsharpforfunandprofit.com/posts/partial-application/#designing-functions-for-partial-application)).

So, usually, you include a file named wrapper.go in the package,
prepare functions and types for Folang in it,
and use it from Folang.

Here, we will create a JoinTail in wrapper.go that simply changes the order of the arguments in the filepath Join we saw earlier, and call it.

### Creating a new package

Since we are going to create a Go command that uses multiple files, let's create a new directory and use go mod init in the usual Golang style.

Here, do the following:

- Create a directory called `call_go`
- Run `go mod init call_go`
- Copy `fc` and `pkg_all.foi`

### Contents of wrapper.go

When you are ready to create a new command of Go, the next step is to prepare the wrapper.go file.

This time, we will create a function called JoinTail with two arguments.
This function looks like this:

```golang
package main
import "path/filepath"

func JoinTail(tail string, head string) string {
  return filepath.Join(head, tail)
}
```

Note that the order of the arguments passed to filepath.Join is reversed.

The reason why this is called JoinTail and not JoinHead is because the order and names are based on the assumption of a pipeline.

You can see what it means by actually using it in a pipeline, so let's call it from Folang.

### Calling Go functions in the same package from Folang

Create a file called call_go.fo as follows.

```
// call_go.fo
package main

import frt

package_info _ =
  let JoinTail: string->string->string

let main () =
  "/home/karino2"
  |> JoinTail "src/folang"
  |> JoinTail "samples"
  |> JoinTail "README.txt"
  |> frt.Println

```

The output is as follows.

```
/home/karino2/src/folang/samples/README.txt
```

Let's take a look at the code in order.

### Specifying an underscore in package_info

package_info is as follows.

```
package_info _ =
  let JoinTail: string->string->string
```

The difference from last time is that an underscore is specified where the package name is written.

This way, JoinTail will be added to the same name space.

This means that functions of files in the same package can be called like this.

Basically, it is common to write various Go functions in wrapper.go, and write only the functions you will use like this in package_info underscore in Folang.

package_info can be written multiple times in the same file, so it is common to place package_info near the functions that use them.
(See [fc's parse_state.fo](../../fc/parse_state.fo), etc.)

This ability to easily call Go functions in the same package is an important feature of Folang, and is a feature that I would like to actively use.

### Let's look at the pipeline code

Next, let's look at the pipeline code.

```fsharp
  "/home/karino2"
  |> JoinTail "src/folang"
  |> JoinTail "samples"
  |> JoinTail "README.txt"
  |> frt.Println
```

If you want to add "src/folang" to the end of a string, make it the first argument and the destination to add it to the second,
and the second argument will be received by the pipeline, so you can concatenate them by the pipeline.

This is the order of arguments assuming a pipeline.

This is the same for all other functions, like AddPrefix or TrimSuffix.

When adding an element to a slice, the element comes first and the slice comes after.

This order is reversed in normal languages,
so if you are writing a program of a certain size, you will end up wrapping all the libraries.

### A few more cases for using wrapper.go

wrapper.go is not only needed to change the order of arguments, but is also needed in various cases.

Flang does not support pointers, but if you want to pass something by reference, you can create a type alias for the pointer type in wrapper.go and refer to it to treat it as a reference type (Folang has no knowledge of this).

For example, to pass a Golang bytes.Buffer to multiple functions, you pass a pointer to the bytes.Buffer,
which you can write in wrapper.go as follows:

```golang
type Buffer = *bytes.Buffer
```

You can then use it from the .fo file as follows:

```
package_info _ =
  type Buffer
  let BufferNew: ()->Buffer
  let Write: Buffer->string->()
  let String: Buffer->string
```

[pkg/buf](../../pkg/buf) does exactly this.

Also, it is customary to write a file that contains only the package_info in a file with the extension .foi.

The `fc` command does not generate the corresponding .go file when the file extension is `.foi`.

Other than that, it actually does exactly the same thing as the `.fo` file.

This is the true nature of pkg_all.foi, which is always passed to the fc command.

## Summary

- You can call functions by writing the type of the function in package_info.
- If you specify an underscore in the package name of package_info, it will be added to the namespace of the current package.
- You can call functions in Go files in the same package.
- For programs of a certain size or larger, use a .go file wrapped in a pipeline-friendly form from the .fo file.
- Use a .go file to avoid NYI and missing functions in .fo.

## Further readings

This is the end of this tutorial.

At the time of writing this document, there is not much other documentation, but here are some links to other documents.

[What is Folang? ](../WhatIsFolang.md) describes the concept and philosophy.

[Sample README.md](../../samples/README.md) lists the code used to verify basic functions,
and is a good way to see what functions are available and what kind of Go code is generated.

As an example of a simple tool, [cmd/build_sample_md](../../cmd/build_sample_md/)
is a tool that generates the Markdown mentioned above from sample source,
and can be used as an example of writing such a tool in a throw-away manner.

Also, [fc command's main.fo](../../fc/main.fo)
can be used as a reference when writing command line tools.

Differences with F# and Golang that require attention are described in the [spec notes](../specs/note.md) (although these are just development notes, so they are scribbled down).

The most accurate way to find out what functions have been implemented is the code in
TestTranspileContain and TestTranspileContainsMulti in [fc/fc_parser_test.go](../../fc/fc_parser_test.go).

It's not very easy to read, but basically all of the implemented functions are in here.

The most complete source code is the [fc command itself](../../fc/).
This contains everything you need to create thousands of lines of code with Folang.