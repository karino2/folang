# 4. Goで書いたラッパーを呼ぶ

前回: [3. スライスとパイプとMap](3_SlicePipeMap_ja.md)

今回はGoで書いた自分のコードをFolangから呼ぶ例を見ていきます。

## FolangはGoと分業して開発するのが基本

Folangは全ての機能を含む事を目指していません。

例えばループや変数への破壊的代入、メソッド呼び出し、ポインタなどがありません。

これらの機能を使う部分はGolangで書き、
それをFolangから呼ぶ、というのを基本的なスタイルとしています。

Folangにとっては、Golangとの共同作業は例外では無くて基本となります。
ある程度以上のFolangプログラムではGoで書かれる部分が含まれる事でしょう。

ここではGolangとの共同作業をどう行うかを見ていきます。

## Goのパッケージを直接呼ぶ例

第二回の[2. frtパッケージを使う (pkg_all.foiの説明)](2_UseFrtPackage_ja.md)でも少し触れましたが、
package_infoというものを書く事で、
fcに外で定義された関数の存在を知らせる事が出来ます。

基本的にFolangはフリースタンディング関数が多いライブラリはそのまま使える事が多く、
オブジェクト的なライブラリはラップして使う必要があります。

まずは良く使うものとして、Golangのfilepathパッケージを使ってみます。
ここでは

- `filepath.Dir`
- `filepath.Join`
- `filepath.Base`

を使ってみます。

```
package main

import frt
import "path/filepath"

package_info filepath =
  let Dir: string->string
  // 本当の型はvarargだが、2引数として使う
  let Join: string->string->string
  let Base: string->string

let main () =
  let test = "/home/karino2/src/folang/README.md"
  filepath.Dir test |> frt.Println
  filepath.Base test |> frt.Println
  let dir = filepath.Dir test
  filepath.Join dir "hello.txt" |> frt.Println

```

これは以下のように出力されます。

```
/home/karino2/src/folang
README.md
/home/karino2/src/folang/hello.txt
```

### package_infoの書き方

package_infoの最初の部分を抜き出すと以下のようになります。

```
package_info filepath =
  let Dir: string->string
```

まず `package_info` キーワードのあとにパッケージ名を書きます。
これは生成されるGoコードで例えば filepath.Dirとなってほしければfilepath、と書くという具合です。

次の行は関数のシグニチャを書きます。
`let`のあとに名前を書き、`:`を置いたあとに関数の型を書きます。

関数の型は引数の型を`->`でつなげて、最後にreturnの型を`->`でつなげます。

例えば、Golangで言う所の、以下のような型は、

```golang
func add(a int, b int) int {...}
```

`int->int->int` と書きます。
引数とreturnの型の区別が分かりにくく慣れが必要ですが、
慣れれば部分適用と相性の良い記述方法です。

もう一つ例を挙げておくと、以下の関数は

```golang
func toStr(str string, a int) string {...}
```

以下のようになります。

```
string->int->string
```

読む時は、最後の矢印のあとをreturnの型と解釈して、
それ以外を`,`で区切られていると思って読むとGolangの関数の型に翻訳出来ます。

## wrapper.goに関数を書いてそれを呼ぶ

Golangのパッケージをそのまま呼ぶのは、ちょっとした関数を使いたい時には便利ですが、
多くの場合には困る事になります。
Folangはポインタやメソッド呼び出し、可変長引数、多値のreturnなど、多くの機能をサポートしていません。
これらが使われている場合、そのままでは使う事は出来ません。

また、使う事が出来たとしても、Golangの関数は普通カリー化とパイプラインを前提とした引数の順番になっていません。
（詳しく知りたい方は[Partial application - F# for fun and profit](https://fsharpforfunandprofit.com/posts/partial-application/#designing-functions-for-partial-application)のDesigning functions for partial applicationあたりを参照）。

そこで通常は、wrapper.goという名前のファイルをパッケージ内に含めて、
そこにFolang向けに関数や型を用意し、
それをFolangから使います。

ここでは先ほど見たfilepathのJoinの引数の順番を変えただけのJoinTailをwrapper.goで作り、それを呼んでみたいと思います。

### 新規パッケージの作成

複数ファイルを使ったGoのコマンドを作るので、通常のGolangのスタイルに則り、
新しいディレクトリを作ってgo mod initしましょう。

ここでは、

- call_goというディレクトリを作る
- `go mod init call_go` する
- `fc`と`pkg_all.foi`をコピーする

という事をやってください。

### wrapper.goの内容

Go用の新たなコマンドを作る準備が出来たら、次はwrapper.goファイルの用意です。

今回は、2引数のJoinTailという関数を作ります。
これは以下のような関数です。

```golang
package main
import "path/filepath"

func JoinTail(tail string, head string) string {
  return filepath.Join(head, tail)
}
```

filepath.Joinに渡す順番と引数の順番が逆になっている事に注目してください。

なぜこれがJoinHeadでは無くてJoinTailなのかは、パイプラインを前提とした順番と名前となります。
実際にパイプラインで使ってみるとどういう事かわかるので、
これをFolangから呼んでみましょう。

### Folangから同じパッケージ内のGo関数を呼ぶ

call_go.foというファイルを以下のように作ります。

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

出力は以下。

```
/home/karino2/src/folang/samples/README.txt
```

では順番にコードを見ていきましょう。

### package_infoにアンダースコアを指定する

package_infoは以下のような内容になっています。

```
package_info _ =
  let JoinTail: string->string->string
```

前回と違う所として、パッケージの名前を書く所にアンダースコアを指定しています。
こうすると、JoinTailは同じ名前スペースに追加されるようになります。
つまり同じpackage内のファイルはこうやって呼ぶ事が出来るようになります。

基本的にはwrapper.goにGoの関数をいろいろ書き、Folangの方にはこうやって使う関数だけをpackage_infoのアンダースコアで書くのが一般的です。

package_infoは同じファイルに何回書いてもいいので、
それらを使う関数群のそばにpackage_infoを置いたりなどは一般的です。
（[fcのparse_state.fo](../../fc/parse_state.fo)などを参照）

この、同じパッケージ内のGo関数を簡単に呼ぶ事が出来る、
というのはFolangの重要な特徴で、
積極的に活用していきたい機能です。

### パイプラインのコードを見てみる

次にパイプラインのコードを見てましょう。

```fsharp
  "/home/karino2"
  |> JoinTail "src/folang"
  |> JoinTail "samples"
  |> JoinTail "README.txt"
  |> frt.Println
```

後ろに "src/folang"を追加したい場合は、それを最初の引数にして、追加する先を二番目にすると、
二番目の引数をパイプラインで受け取る事になるのでそれを連結していける訳です。
これがパイプラインを前提とした引数の順番となります。

AddPrefixでもTrimSuffixでもこの辺は変わりません。
要素をスライスに追加する場合も要素が先でスライスがあとになります。
この順番が普通の言語では逆になるので、
ある程度大きなプログラムを書くなら、結局はどのライブラリもラップする事になります。

### wrapper.goを使うケースについてもう少し

wrapper.goは引数の順番を変えるだけでは無く、いろいろなケースで必要となります。

Flangはポインタをサポートしていませんが、
参照渡しとして扱いたいものは、wrapper.goでポインタ型をtypeエイリアスを作って参照する事で、
参照型として扱う事が出来ます（Folangは何も関知しない）。

例えばGolangのbytes.Bufferを複数の関数に渡してやりとりする為にはbytes.Bufferのポインタを渡す事になりますが、
それは以下のようにwrapper.goに書き、

```golang
type Buffer = *bytes.Buffer
```

それを.foファイルから以下のように使えば良い。

```
package_info _ = 
  type Buffer
  let BufferNew: ()->Buffer
  let Write: Buffer->string->()
  let String: Buffer->string
```

まさにこういう事を行っているのが[pkg/buf](../../pkg/buf)になります。

また、このpackage_infoだけを集めたファイルを.foiという拡張子で書く決まりになっています。
`fc`コマンドは、拡張子が`.foi`の時には対応する.goファイルを生成しません。
それ以外は実は`.fo`ファイルと全く同じ事をしています。

これがfcコマンドにいつも渡すpkg_all.foiの正体です。

## まとめ

- package_infoに関数の型を書くと呼べるようになる
- package_infoのパッケージ名にアンダースコアを指定すると、現在のパッケージの名前空間に追加される
  - 同じパッケージ内のGoファイルの関数を呼べるようになる
- ある程度以上の規模のプログラムなら、パイプライン向きな形に.goファイルでラップしたものを.foファイルから使う
- .foのNYIや足りない機能は、.goファイルを使って回避する

## 終わりに

以上で本チュートリアルはおしまいです。

本文書執筆時点ではあまり他のドキュメントはありませんが、幾つか他へのリンクなどを挙げておきます。

[Folangとは何か？](../WhatIsFolang_ja.md)はコンセプトや哲学などが書かれています。

[サンプルのREADME.md](../../samples/README.md)は基本的な機能の検証に使ったコードを並べたものですが、
どういった機能があってどういうGoコードが生成されるのかを眺めるのに良いでしょう。

簡単なツールの例としては[cmd/build_sample_md](../../cmd/build_sample_md/)が、
サンプルのソースから前掲のマークダウンを生成するツールで、
こういうツールを書き捨てで書く場合の例として見る事が出来ます。

また、[fcコマンドのmain.fo](../../fc/main.fo)は、
コマンドラインのツールを書く時の参考になると思います。

F#やGolangとの違いで注意が必要な所は、[specノート](../specs/note_ja.md)に記述があります（開発中のメモくらいなので書き殴りですが）。

実装済みの機能を一番正しく知る方法は、[fc/fc_parser_test.go](../../fc/fc_parser_test.go)にある、
TestTranspileContainとTestTranspileContainsMultiにあるコードです。
あまり見やすいものではありませんが、基本的には実装済みの機能は全てここに入っています。

ソースコードとして一番本格的なものは[fcコマンドそのもの](../../fc/)です。
これはFolangで数千行規模のコードを作る場合に必要な事の全てが詰まっています。