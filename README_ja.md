# Folang

FolangはGoにトランスパイルする為に新規に設計された関数型言語で、
Goにトランスパイルされます。
仕様はF# に強く影響を受けています。
Folangのトランスパイラ自身もFolangで書かれています（セルフホスト）。

より詳細は[Folangとは何か？](docs/WhatIsFolang_ja.md)を参照ください。

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

## セットアップ

[tutorials/1_GettingStarted_ja.md](docs/tutorials/1_GettingStarted_ja.md)を参照ください。

## チュートリアル

[tutorials/Index_ja.md](docs/tutorials/Index_ja.md)

## レポジトリ構成

- cmd このサイトを作るツールなど
- docs ドキュメント、チュートリアル
- fc Folangのトランスパイラ
- pkg Folangの標準ライブラリ
- samples Folangの開発中に機能確認に使っているサンプル
- tinyfo 初期に使われていたGo言語で書かれたトランスパイラ、現在は使われていないが記録の為に残してある
