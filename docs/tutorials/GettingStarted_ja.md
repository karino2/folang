# Getting Started Folang

セットアップと最初のHello Worldが動くまでを解説します。

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

### hello.foの内容を少しだけ解説

golangとF#やOcaml系の言語の両方の知識がある人なら、見ただけでほぼ意味がわかる内容と思いますが、
簡単に内容も解説しておきます。

最初のpackage文とimport文はそのままgolangと同様なのでいいでしょう。生成した結果も同じ内容となっています。

次の文が関数を定義しているように見えます。
folangでは以下のように関数を定義します。

```
let main () =
  // ここに本体を書く
```

関数の定義は `let` です。そして引数はスペース区切りで書くのですが、
何も引数が無い場合は `()` という特殊なものを置きます。

これは引数がある例が出てくるまではそういうものと思ってください。

そのあとに「=」があるのもgolangと違う所です。

なお、関数の戻りの型は普通は書きません。最後の式の型が関数の型となります。

そして関数のbody部はインデントします。インデントが終わる所がブロックの終わりと解釈されます。

GoEvalなどは少しトリッキーなのでここでは深入りせずに先に進みましょう。

簡単に以上の内容をまとめておきます。

- 関数の定義は `let`
- 引数無しは　`()` で表す
- 引数のあとには `=` が必要
- 本体はインデントする(Pythonとかと同じ)

## 次回: frtパッケージを使う (pkg_all.foiの説明)

[frtパッケージを使う (pkg_all.foiの説明)](UseFrtPackage_ja.md)