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
