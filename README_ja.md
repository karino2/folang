# Folang

FolangはGoにトランスパイルする為に新規に設計された関数型言語で、
Goにトランスパイルされます。
仕様はF# に強く影響を受けています。
Folangのトランスパイラ自身もFolangで書かれています（セルフホスト）。

より詳細は[Folangとは何か？](docs/WhatIsFolang_ja.md)を参照ください。

## 簡単な例

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

他の例は[samples/README.md](samples/README.md)を参照

## セットアップ

[tutorials/1_GettingStarted_ja.md](docs/tutorials/1_GettingStarted_ja.md)を参照ください。

## チュートリアル

[tutorials/Index_ja.md](docs/tutorials/Index_ja.md)

## 仕様関連

F#やGolangと違って注意が必要な所を中心としたメモ。

[specs/note_ja.md](docs/specs/note_ja.md)

## レポジトリ構成

- cmd このサイトを作るツールなど
- docs ドキュメント、チュートリアルやスペックなど
- fc Folangのトランスパイラ
- pkg Folangの標準ライブラリ
- samples Folangの開発中に機能確認に使っているサンプル
- tinyfo 初期に使われていたGo言語で書かれたトランスパイラ、現在は使われていないが記録の為に残してある

## ゴールとプライオリティ

どういうものを作りたいのかを最初に書いておこう。

### ゴールとノンゴール

**ゴール**

- 簡潔に書ける
   - パイプラインでスライスを処理していける
- コマンドラインの簡単なツールを書くのをターゲットとする
- Golangの豊富なパッケージを使える
- 生成されるGoのコードが自然で、どういうコードが生成されるか予想出来る
   - あまりリストとか再帰とかは使わず、sliceメインでやっていく
   - 少なくともデバッグ出来る程度のコード
- 軽量なシングルバイナリ（Goのコードとしてデプロイ）
  - 5000LOC 未満程度のコードがサクサク動く

**ノンゴール**

- パフォーマンスはそれほど気にしない
- MLやF#互換は目指さない
- Go無しで全部書くのは目指さない
   - むしろGo向きな処理は気軽にgoで書いてFolangから呼び出すスタイルを推奨したい
- 完全さは目指さない（変な制限があってもあまり使わなければOK）

### プライオリティ

1. 簡潔に書ける＞Goとして自然
2. 生成されるコードがGoとして自然＞ML的な良さ、一貫性
   - let はstatementにする
   - レコードは単なるstructにする
   - Unionはインターフェースとする
   - exhaustive checkなどはほどほど
   - 実用上困らない程度に実現出来るならアドホックな処理で妥協
3. 実装が簡単＞完全性

大なり、の記号は相対的なプライオリティを明確にするためにつけている（左側がプライオリティの項目、つまり1番目は「簡潔に書ける」が重要という意味）。

