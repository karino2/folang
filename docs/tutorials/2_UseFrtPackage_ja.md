# 2. frtパッケージを使う (pkg_all.foiの説明)

前回の、[1. Getting Started](1_GettingStarted_ja.md)で使ったGoEvalは特殊な関数で、
一切の外部パッケージを使わずに実行出来ますが、普通はこんな事はしません。

通常はfrtというFolang標準のパッケージを使います。

frtはFolang RunTimeの略で、ほとんどのFolangプログラムに必要とする基本機能を提供します。

ここでfrtを使ってみましょう。frtのPrintlnを使ってみます。

## go mod initの実行(やってなければ)

前回の[1. Getting Started](1_GettingStarted_ja.md)の通りに作業していればやっているはずですが、
今回はGo言語から見ると外部パッケージとなるfrtを使うので、
go mod initの設定が必要となります。

以下をやっておきましょう。

```
hello_fo $ pwd
~/helo_fo
hello_fo $ go mod init hello_fo
hello_fo $ go mod tidy
```

## hello_frt.fo ファイルの作成

frtを使うhello_frt.foファイルを作りましょう。
中は以下のようにします。

```
// hello_frt.fo
package main
import frt

let main () =
  frt.Println "Hello World"

```

まず二行目のimportが新しくなっています。

```
import frt
```

標準のパッケージはダブルクオート無く指定する事が出来る、
という仕様になっていて、これはダブルクオートありのimportで書く以下と同じ意味になります。

```
import "github.com/karino2/folang/pkg/frt"
```

基本的には folang/pkg の下にはFolangから使う事を前提としたライブラリが置かれています。

ちなみにコメントはGolang同様、`//`の行コメントと `/* */` のブロックコメントです。

## hello_frt.foのトランスパイルとpkg_all.foiファイル

さて、先ほどのファイルをトランスパイルすると以下のようなエラーメッセージが出ると思います。

```
 ./fc hello_frt.fo
transpile: hello_frt.fo
panic: Unknown var ref: frt.Println

goroutine 1 [running]:
github.com/karino2/folang/pkg/frt.Panic(...)
...
以下長いスタックトレース
...
```

大量のエラーメッセージが出てぎょっとしますが、
まだ開発中のためです。

意味のあるメッセージは一番上の「panic: Unknown var ref: frt.Println」だけです。

これはfrtというパッケージの情報をfcコマンドが知らない、という事を意味します。

その情報が書いてあるのがpkg_all.foiになります。
これをhello_frt.foの前に指定してやります。

```
$ ./fc pkg_all.foi hello_frt.fo
$ go fmt gen_hello_frt.go
```

これで以下のようなファイルが生成されるはずです。

```golang
// gen_hello_frt.go
package main

import "github.com/karino2/folang/pkg/frt"

func main() {
	frt.Println("Hello World")
}
```

これで通常のGo言語の実行と同じようにgo runを実行すると、

```
$ go run gen_hello_frt.go
gen_hello_frt.go:3:8: no required module provides package github.com/karino2/folang/pkg/frt; to add it:
	go get github.com/karino2/folang/pkg/frt
```

と言われるので、go getを実行します。

```
$ go get
go: added github.com/google/go-cmp v0.6.0
go: added github.com/karino2/folang/pkg/frt v0.0.0-20250220122800-4ff80daf0a9a
```

そしてgo runを実行すれば実行出来ます。

```
$ go run gen_hello_frt.go
Hello World
```

以上をまとめると

- Folangは普通frtをimportする
- frtをimportする場合は pkg_all.foiというファイルをfcファイルの前に置く
- 外部パッケージの使用には通常のGo言語と同様にgo getなどが必要

となります。

## frt.Printlnの呼び出しについて簡単に解説

今回は初めてFolangらしいコードが出てきたので、
基本的な事を解説しておきます。

再掲すると以下のようになっていました。

```
let main () =
  frt.Println "Hello World"
```

関数呼び出しは以下の部分です。

```
  frt.Println "Hello World"
```

関数はスペース区切りで引数を並べて実行します。
Golangとは異なり、カッコをつけません。

frt.Printlnは文字列を引数にとって、結果を返さない関数です。
この辺はGolangのfmt.Printlnと同様ですね。

## pkg_all.foiの内容について少しだけ

さて、先程のpkg_all.foiはただのテキストファイルです。
このテキストファイルの最初の方にfrtのインターフェースが書かれています。
少しその部分を抜粋してみましょう。

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

そして今回のプログラムに必要な情報はPrintlnだけです。

これをhello_frt.foに書いても、実は実行する事が出来ます。

```
// hello_frt2.fo
package main
import frt

// この2行を追加
package_info frt =
  let Println: string->()

let main () =
  frt.Println "Hello Frt2"
```

こうするとpkg_all.foiなどというファイルを指定しなくてもトランスパイル可能です。
実際にやってみましょう。

```
$ ./fc hello_frt2.fo
$ go run gen_hello_frt2.go
Hello Frt2
```

トランスパイラとしてはそういう型の関数がある、
と思って型の解決をしてGoのコードを生成するだけです。
実際に生成したコードが正しく実行されるようにするのはプログラマ側の責任になります。

これは自分でラッパを書いて実行する時に重要になるので心に留めておきましょう。

## 第二回まとめ

- Folang標準のパッケージをダブルクオート無しでimport
- frtは普通importする
- pkg_all.foiファイルをfcコマンドの先頭に指定する
- コメントは `//` と `/* */`

## 次回: スライスとパイプとMap

[3. スライスとパイプとMap](3_SlicePipeMap_ja.md)