## Folang Sample 


### Two function example

```
package main
import "fmt"

let hello (msg:string) = 
    GoEval "fmt.Printf(\"Hello %s\\n\", msg)"

let main () =
    hello "World"

```

generated go: [gen_two_func.go](./gen_two_func.go)


### Record

```
package main

type hoge = {X: string; Y: string}

let ika () =
    {X="abc"; Y="def"}

```

generated go: [gen_record.go](./gen_record.go)


### Basic Union sample

```
package main

type IntOrString =
| I of int
| S of string

let ika () =
   I 123 

```

generated go: [gen_union_simple.go](./gen_union_simple.go)


### Union with match

```
package main
import "fmt"

type IorS =
 | IT of int
 | ST of string

 let main () =
  match IT 3 with
  | IT ival -> GoEval "fmt.Printf(\"ival=%d\\n\", ival)"
  | ST sval -> GoEval "fmt.Printf(\"ival=%s\\n\", sval)"


```

generated go: [gen_union_match.go](./gen_union_match.go)


### Match with no value

```
package main

type IorS =
 | IT of int
 | ST of string

 let ika () =
  match IT 3 with
  | IT _ -> "i hit" 
  | ST sval -> sval


```

generated go: [gen_union_match_noval.go](./gen_union_match_noval.go)


### Match with no content

```
package main
import "fmt"

type AorB =
 | A 
 | B


let ika (ab:AorB) =
  match ab with
  | A -> "a match"
  | B -> "b match"


```

generated go: [gen_union_nocontent.go](./gen_union_nocontent.go)


### Folang standard package example

```
package main

import frt
import buf

let main() =
  let bb = buf.New ()
  buf.Write bb "hello"
  buf.Write bb "world"
  let res = buf.String bb
  frt.Println res


```

generated go: [gen_std_pkg.go](./gen_std_pkg.go)


### Standard package with generics

```
package main

import frt
import slice

let main() =
  let s = [5; 6; 7; 8]
  let s2 = slice.Take 2 s
  frt.Printf1 "%v\n" s2


```

generated go: [gen_pkg_generics.go](./gen_pkg_generics.go)


### Tuple example

```
package main

import frt

let ika () =
   frt.Snd (123, "abc")

let main () =
  let s = ika ()
  frt.Println s
```

generated go: [gen_tuple.go](./gen_tuple.go)


### Pipe operator

```
package main

import frt
import slice

let main() =
  let s = [5; 6; 7; 8]
  let s2 = s |> slice.Take 2
  frt.Printf1 "%v\n" s2


```

generated go: [gen_pipe.go](./gen_pipe.go)


### Using go package.

```
package main

import frt
import "path/filepath"

package_info filepath =
  let Dir: string->string
  let Base: string->string

let main () =
  let target = "/home/karino2/src/folang/samples/README.md"
  filepath.Dir target |> frt.Println
  filepath.Base target |> frt.Println

```

generated go: [gen_go_pkg.go](./gen_go_pkg.go)


### Map

```
package main

import frt
import slice

let conv (i:int) =
  frt.Sprintf1 "a %d" i

let main() =
  let s = [5; 6; 7; 8]
  let s2 = slice.Map conv s
  frt.Printf1 "%v\n" s2


```

generated go: [gen_map.go](./gen_map.go)


### Destructuring let example

```
package main

import frt

let ika () =
  (123, "abc")

let main () =
  let (a, _) = ika ()
  frt.Sprintf1 "a=%d" a
  |> frt.Println
  

```

generated go: [gen_destr_let.go](./gen_destr_let.go)


### Explicit type argument

```
package main

import frt
import slice


let main () =
  let s = slice.New<string> ()
  let ss = slice.PushHead "hoge" s
  let e = slice.Head ss
  frt.Println e

```

generated go: [gen_generic_specify.go](./gen_generic_specify.go)


### Dictionary

```
package main

import frt
import dict

let main () =
  let d = dict.New<string, int> ()
  dict.Add d "hoge" 123
  dict.Add d "ika" 456
  let i1 = dict.Item d "hoge"
  let i2 = dict.Item d "ika"
  frt.Printf1 "sum=%d\n" (i1+i2)

```

generated go: [gen_dict_sample.go](./gen_dict_sample.go)


### Generic function

```
package main

import frt
import slice

let hoge a =
  slice.Head a

let main () =
  let b = [1; 2; 3]
  let c = hoge b
  frt.Printf1 "%d\n" c

```

generated go: [gen_generic_func.go](./gen_generic_func.go)


### Type inference

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

generated go: [gen_type_inference.go](./gen_type_inference.go)


### Shorthand property access notation

```
package main

import slice
import frt
import strings

type Rec = {SField: string; IField: int}

let main () =
  let r1 = {SField="abc"; IField=123}
  let r2 = {SField="def"; IField=456}
  [r1; r2] |> slice.Map _.SField |> strings.Concat ", " |> frt.Println 

```

generated go: [gen_shorthand_prop.go](./gen_shorthand_prop.go)

