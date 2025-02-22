# 1. Getting Started Folang

I will explain how to set up and get the first Hello World working.

## Setup

Currently, I have only built it manually and tried it.

You can create the fc command by following the steps below.

```
$ git clone https://github.com/karino2/folang.git
$ cd folang/fc
$ go build
```

You will also need the file pkg/pkg_all.foi under folang.

For the time being, I am working by copying these two to the development directory.

```
$ cp fc ../../your/target/dir/
$ cp ../pkg/pkg_all.foi ../../your/target/dir/
```

## Hello Folang

Let's start with the simplest Hello World.

### Running go mod

Folang is transpiled into Go and executed as Go, so you need to prepare an environment for running Go.

Also, to use external packages such as frt, which will be described later, you need to set up go mod just like normal Go.

Let's try doing the following.

```
$ mkdir hello_fo
$ cd hello_fo
$ go mod init hello_fo
go: creating new go.mod: module hello_fo
go: to add module requirements and sums:
go mod tidy
$ go mod tidy
```

Also, copy the fc command created in the setup and the pkg_all.foi file, which will be explained later.
Please execute the following, substituting it for your own environment.

```
$ cp ../folang/fc/fc ./
$ cp ../folang/pkg/pkg_all.foi ./
```

This completes the preparation.

### Creating and transpiling the hello.fo file

Create a suitable directory and create a text file named hello.fo with the following contents.

```
$ cat hello.fo
package main
import "fmt"

let main () =
GoEval "fmt.Println(\"Hello World\")"

```

Folang is a language where indentation is important. Make sure the last line is a blank line without any spaces.

Also, be careful about the backslash. By the way, don't worry, this is only done in the first example.

GoEval is a special function that passes the argument string directly through to the Go file.

Transpile this with the fc command created in the setup.

If you copy it to the same directory, run it as follows.

```
$ ./fc hello.fo
transpile: hello.fo
$ ls
fc
gen_hello.go
hello.fo
```

This should create a file called gen_hello.go.
If you look inside, you'll see the following.

```golang
// gen_hello.fo
package main
import "fmt"
func main() {
fmt.Println("Hello World")
}
```

This is a normal Go file, so
you can run it as follows.

```
$ go run gen_hello.go
Hello World
$
```

fc doesn't indent, so you would normally use go fmt.
I've made it into a shell script along with the one that includes pkg_all.foi, which I'll explain later, but
let's run it manually first.

```
$ go fmt gen_hello.go
gen_hello.go
```

This will look like this.

```golang
package main

import "fmt"

func main() {
fmt.Println("Hello World")
}
```

Now it looks nice. This is more about Golang than Folang.

### A little explanation of the contents of hello.fo

If you are familiar with both Golang and F# or Ocaml languages, you can probably understand the meaning just by looking at it, but I will briefly explain the contents.

The first package and import statements are the same as in Golang, so that's fine. The generated result is also the same.

The next statement looks like it defines a function.
In Folang, functions are defined as follows.

``
let main () =
// Write the body here
```

The function definition is `let`. And arguments are written separated by spaces, but if there are no arguments, a special `()` is placed.

Please assume that this is the case until an example with arguments appears.

Another difference from Golang is that there is an "=" after that.

The return type of a function is not usually written. The type of the last expression becomes the type of the function.

The body of the function is indented. The end of the indentation is interpreted as the end of the block.

GoEval and other functions are a bit tricky, so let's move on without going into too much detail here.

A quick summary of the above.

- Function definition is `let`
- No arguments are expressed with `()`
- Arguments must be followed by `=`
- The body is indented (same as Python, etc.)

## Next time: Using the frt package (Explanation of pkg_all.foi)

[2. Using the frt package (Explanation of pkg_all.foi)](2_UseFrtPackage.md)