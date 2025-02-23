# Folang

Folang is a functional language designed from the ground up to be transpiled into Go.

Its specification is heavily influenced by F#.

The Folang transpiler itself is written in Folang (self-hosted).

For more information, see [What is Folang?](docs/WhatIsFolang.md).

[日本語版 README_ja.md](README_ja.md)

## Simple example


```
package main
import frt

import slice
import strings

let main () =
  [1; 2; 3]
  |> slice.Map (frt.Sprintf1 "This is %d")
  |> strings.Concat ", "
  |> frt.Println

```

=>

```golang
package main

import "github.com/karino2/folang/pkg/frt"
import "github.com/karino2/folang/pkg/slice"
import "github.com/karino2/folang/pkg/strings"

func main() {
	frt.PipeUnit(
    frt.Pipe(
      frt.Pipe(
        ([]int{1, 2, 3}),
        (func(_r0 []int) []string {
		      return slice.Map((func(_r0 int) string { return frt.Sprintf1("This is %d", _r0) }), _r0)
	  })),
    (func(_r0 []string) string {
       return strings.Concat(", ", _r0)
    })), frt.Println)
}
```

```
package main

import frt

let ApplyL fn tup =
  let nl = frt.Fst tup |> fn
  (nl, frt.Snd tup)


let add (a:int) b = 
  a+b

let main () =
  (123, "hoge")
  |> ApplyL (add 456)
  |> frt.Printf1 "%v\n" 
```

=>

```golang
package main

import "github.com/karino2/folang/pkg/frt"

func ApplyL[T0 any, T1 any, T2 any](fn func(T0) T1, tup frt.Tuple2[T0, T2]) frt.Tuple2[T1, T2] {
	nl := frt.Pipe(frt.Fst(tup), fn)
	return frt.NewTuple2(nl, frt.Snd(tup))
}

func add(a int, b int) int {
	return (a + b)
}

func main() {
	frt.PipeUnit(
    frt.Pipe(
      frt.NewTuple2(123, "hoge"),
      (func(_r0 frt.Tuple2[int, string]) frt.Tuple2[int, string] {
    		return ApplyL((func(_r0 int) int { return add(456, _r0) }), _r0)
	 })),
   (func(_r0 frt.Tuple2[int, string]) { frt.Printf1("%v\n", _r0) }))
}
```

For other examples, see [samples/README.md](samples/README.md).

## Setup

See [tutorials/1_GettingStarted.md](docs/tutorials/1_GettingStarted.md).

## Tutorial

[tutorials/Index.md](docs/tutorials/Index.md)

## Specifications

Notes focusing on points that require attention, diff from F# and Golang.

[specs/note.md](docs/specs/note.md)

## Repository structure

- cmd Tools for making this site
- docs Documents, tutorials and specs.
- fc Folang transpiler
- pkg Folang standard library
- samples Samples used to check functionality during Folang development
- tinyfo Transpiler written in Go language used in the early days, no longer in use but kept for record keeping

## Goals and priorities

Write down what I want to make.

### Goals and non-goals

**Goals**

- Can be written succinctly
  - Can process slices in a pipeline
- Targets writing simple command line tools
- Can use Golang's rich packages
- Generated Go code is natural, and you can predict what kind of code will be generated
  - Doesn't use lists or recursion much, and mainly uses slices
  - At least code that can be debugged
- Lightweight single binary (deployed as Go code)
  - Code of less than 5000 LOC runs smoothly

**Non-goals**

- Doesn't care much about performance
- Doesn't aim for ML or F# compatibility
- Doesn't aim to write everything without Go
  - Rather, I would like to recommend a style where Go-friendly processing is written casually in Go and called from Folang
- Doesn't aim for completeness (it's OK even if there are strange restrictions as long as they are not used much)

### Priority

1. Can be written succinctly `>` Natural as Go
2. Generated code is natural as Go `>` ML-like goodness, consistency
  - let should be a statement
  - Records should be simple structs
  - Unions should be interfaces
  - Exhaustive checks should be kept to a minimum
  - Compromise with ad-hoc processing if practical
3. Easy to implement `>` Completeness

The greater than symbol is used to clarify relative priorities (the left side is the priority item, meaning that the first one is important for "being concise to write").