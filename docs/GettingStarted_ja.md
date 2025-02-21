# Getting Started Folang

## セットアップ

現在の所、手動でビルドをして試す事しかしていません。

以下の手順でfcコマンドが作れます。

```
$ git clone https://github.com/karino2/folang.git
$ cd folang/fc
$ go build
```

また、folangの下の pkg/pkg_all.foi にあるファイルが必要になります。

当面はこの２つを開発ディレクトリにコピーして作業しています。

```
$ cp fc ../../your/target/dir/
$ cp ../pkg/pkg_all.foi ../../your/target/dir/
```

## ハロー Folang

まずは一番簡単なHello Worldから始めましょう。

### go modの実行

folangはgo言語にトランスパイルしてgo言語として実行するため、
goを実行するための環境を整える必要があります。
また、後述するfrtなどの外部パッケージを使う為には、
通常のgoと同様にgo modのセットアップが必要となります。

とりあえず以下のようにしてみましょう。

```
$ mkdir hello_fo
$ cd hello_fo
$ go mod init hello_fo
go: creating new go.mod: module hello_fo
go: to add module requirements and sums:
	go mod tidy
$ go mod tidy
```

また、セットアップで作ったfcコマンドと、あとで説明するpkg_all.foiファイルをコピーします。
適当に各自の環境に読み替えて以下を実行してください。

```
$ cp ../folang/fc/fc ./
$ cp ../folang/pkg/pkg_all.foi ./
```

これで前準備は完了です。

### hello.foファイルの作成とトランスパイル

適当なディレクトリを作り、そこにhello.foという名前のテキストファイルで以下のような内容を作ります。

```
$ cat hello.fo
package main
import "fmt"

let main () =
  GoEval "fmt.Println(\"Hello World\")"

```

folangはインデントが重要な言語です。最後の行は空白を入れず空行になるようにしましょう。
また、バックスラッシュも注意してください。ちなみにこんな事をするのは最初の例だけなので安心してください。

GoEvalは引数の文字列をそのままgoファイルに素通しする、という特殊な関数です。

これをセットアップで作ったfcコマンドでトランスパイルします。
同じディレクトリにコピーしたとすると以下のように実行します。

```
$ ./fc hello.fo
transpile: hello.fo
$ ls
fc
gen_hello.go
hello.fo
```

これでgen_hello.goというファイルが出来ているはずです。
中を見ると以下のようになっています。

```
$ cat gen_hello.fo
package main

import "fmt"

func main() {
fmt.Println("Hello World")
}
```

これは通常のgoファイルとなるので、
以下のように実行出来ます。

```
$ go run gen_hello.go
Hello World
$
```

fcはインデントをしないので普通はgo fmtします。
私は後述するpkg_all.foiを含めるものと合わせて、シェルスクリプトにしていますが、
まずは手動で実行しましょう。

```
$ go fmt gen_hello.go
gen_hello.go

$ cat gen_hello.go
package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}
```

これで綺麗になりました。

## frtライブラリを使う (pkg_all.foiの説明)

GoEvalは特殊な関数で、
一切の外部ライブラリを使わずに実行出来ますが、
通常はfrtという標準のライブラリはほぼ確実に使います。

frtはFolang RunTimeの略で、ほとんどのfolangプログラムに必要とする基本機能を提供します。

ここでもその例を見てみましょう。

### hello_frt.fo ファイルの作成

frtを使うhello_frt.foファイルを作りましょう。
中は以下のようにします。

```
$ cat hello_frt.fo
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

基本的には folang/pkg の下にはfolangから使う事を前提としたライブラリが置かれています。

### hello_frt.foのトランスパイルとpkg_all.foiファイル

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

```
$ cat gen_hello_frt.go
package main

import "github.com/karino2/folang/pkg/frt"

func main() {
	frt.Println("Hello World")
}
```

これで通常のgo言語の実行と同じようにgo runを実行すると、

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

- folangは普通frtをimportする
- frtをimportする場合は pkg_all.foiというファイルをfcファイルの前に置く
- 外部パッケージの使用には通常のgo言語と同様にgo getなどが必要

となります。

### pkg_all.foiの内容について少しだけ

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
$ cat hello_frt2.fo
package main
import frt

// この2行ｗ追加
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
と思って型の解決をしてgoのコードを生成するだけです。
実際に生成したコードが正しく実行されるようにするのはプログラマ側の責任になります。

これは自分でラッパを書いて実行する時に重要になるので心に留めておきましょう。
