# 仕様に関するノート

仕様に関する決定をしたら上に追記していく。
順番が下からになるのでちょっとトリッキー。
文体などはまだ一致していない。

もう少し固まったら整理したい。

こうした決定に関する経緯は[discussion_ja.md](discussion_ja.md)に。

## string interpolationとbacktickのraw string

string interpolationはF# っぽくてraw stringはGolangっぽいので、
組み合わせ的には独自仕様となる。
またエスケープの仕方が独自。

### string interpolation

String interpolationはドル始まりの文字列で、中括弧で囲まれたシンボルが変数として中身に置き換わる。

```
let hoge () =
  let a = 123
  let b = "abc"
  $"a is :{a}, b is {b}"
```

この結果は"a is :123, b is \"abc\""となる。

ブレースのエスケープはバックスラッシュにした。
ここはF#と違う。

```
let hoge () =
  let a = 123
  $"a is :{a}, \{a\}"
```

結果は"a is :123, {a}"となる。


### rawstring

rawstringはgolangに合わせてbacktickにしてある。

```
let hoge () =
  let s = `This is
raw string, "Double quote" and backslash \ is not escaped.
`
  frt.Println s
```

### rawstringのinterpolation

rawstringでもinterpolationが出来る。
ただし、この中でブレースをエスケープする方法は無い。

```
let hoge () =
  let a = 123
  let s = $`This is
raw string, a is "{a}"
`
  frt.Println s
```


## スライスとタプルの型の優先順位

スライスとタプルの優先順位はスライスのsyntaxがGoなので起こるFolang特有の問題。
以下のように決める。

- `[]T*U` は `[](T*U)` とパースされる。
- `T*[]U` は合法

上はちょっと微妙だけれど、良く使うのでカッコをつけるのは嫌だった。

## スライス expression

大括弧で区切りはセミコロンで型名はつけない。（区切りがカンマじゃないのに注意）

```
let ika () =
  [1; 2; 3]
```

## イコールの比較

equalは`=`で、not equalは`<>`で行う。
equalはsliceも中を比較する（go-compのEqualを使う[cmp package - github.com/google/go-cmp/cmp - Go Packages](https://pkg.go.dev/github.com/google/go-cmp/cmp)）

## importのシステムライブラリ

ダブルクオートで括られてないimportはシステムのimportとみなし、内部的にはFolang/pkgへのパスがprefixでついているとみなす。
具体的には以下の２つは同じ意味になる。（生成されるGoコードはどちらも二行目になる）

```
import frt
import "github.com/karino2/folang/pkg/frt"
```

## Genericsのシンタックス（タイプパラメータ）

Goは大括弧だがFSharpは角括弧。Folangは角括弧を採用する。

```
package_info slice =
  let New<T>: ()->[]T

slice.New<int> ()
```


## 外部の型情報

外部のパッケージなどをアクセスするための言語要素。FSharpのシグニチャファイルとかと似たようなもの。

ファイルの拡張子は.foiに書く（ただし.foファイルの中に書く事も出来る、.foiに書いてあると対応するGoファイルが生成されないだけ）。
package_infoというもので定義し、関数はReScriptを真似してletでコロンとしてみる。

```
package_info slice =
   let Length<T>: []T -> int
   let Take<T> : int->[]T->[]T 
```

型はtypeで以下のように書く。

```
package_info buf =
    type Buffer
    let New: ()->Buffer
    let Write: Buffer->()
    let String: ()->string
```

ファイルの拡張子はfoi。（ただし.foファイルに書く事も出来る。foiだと対応するGoファイルが生成されないだけ）

### 自身のネームスペースに外部型情報を追加する

パッケージの名前をアンダースコアにする事で、プレフィクス無しで現在のネームスペースに追加される。

```
package_info _ =
   type wrappedType
   let New: ()->wrappedType
   let doWork: wrappedType -> ()
```

これは、同じパッケージ内にFolang向けのラッパーをwrapper.goなどで作ってその情報を参照する時などに使われる。
メソッドになっているものは普通の関数でラップして使う。
この時カリー化の事を考えてF#的に引数の順番を決めておくと良い。

## コメント

コメントはGolangと同様、CスタイルとC++の一行コメントの２つをサポートする。

```
package main

/*
これはコメントです
*/

let ika () =
  123 // これもコメントです。
```

## GoEval

Goのコードをそのまま文字列として書く、イメージとしてはインラインアセンブラに近い機能。
以下のようなコードは、

```
package main
import "fmt"

let main () =
  GoEval "fmt.Println(\"Hello World\")"
```

以下のコードに展開される。

```
package main
import "fmt"

func main() {
   fmt.Println("Hello World")
}
```

### 戻りの型指定

デフォルトではUnitとみなされる。戻りの型を指定したい場合はタイプパラメータで指定する。

```
   // このsはstring
   let s = GoEval<string> "fmt.Sprintf(\"hoge %d\", 123)"
```

展開されるコードは以下になる。

```
   s := fmt.Sprintf("hoge %d", 123)
```

Go側では型指定はないので、Folang上のsと同じ型になるように指定してやらないといけない。

### 引数の使用

引数のidentifierはGoでもそのまま持ち越されるので、以下のように書く事が出来る。aを使っている事に注意。

```
let ika (a:int) =
   GoEval<string> "fmt.Sprintf(\"hoge %d\", a)"
```

これは以下のように展開される。

```
func ika(a int) string {
   return fmt.Sprintf("hoge %d", a)
}
```

下に生成されるGoコードを意識しながら引数などを使ってやる。この辺はインラインアセンブラと同じ。

Folangで対応してない機能はGoEvalでラップしてやれば割と使える。

## Discriminated Unionの実装

長いので[union_ja.md](union_ja.md)へ。

## 関数定義

基本的な関数定義は以下のようになる。

```
let ika (a:string) (b:string) =
  a+b
```

これは以下に展開される。

```golang
func ika(a string, b string) string {
  return a+b
}
```

なお、引数無しは引数unitが一つと定義する。
GoEvalは引数のコードをそのままGoに流すexpression、型は型パラメータで指定するが指定無しだとUnit。

それを用いると以下のようなコードは、

```
package main
import "fmt"

let main () =
  GoEval "fmt.Println(\"Hello World\")"
```

以下のコードになる。

```golang
import "fmt"

func main() {
   fmt.Println("Hello World")
}
```

メソッドはサポートしない（手でラップする）

### 型推論

引数は推論され、決定出来ないものはgenericsになる。
以下のGoコードは

```fsharp
let add10 a = 
  a + 10
```

以下のコードになる。

```golang
func add10(a int) int {
  return a + 10
}
```

型が不明な場合はgenericsのtype parameterになる。

```fsharp
let secondHead s =
   slice.Item 2 s
```

この場合、sがスライスな事は確定するが要素は確定しないので、以下のようなコードが生成される。

```golang
func sceondHead[T0 any](s []T0) T0 {
  slice.Item(2, s)
}
```

今の所type constraintsはサポートしてないので全部anyになる。

だから以下のコードはGoのコンパイルエラーになる。

```
let ika a b =
  a+b
```

以下のように片方に指定があると、両方同じ型だとは解決されるので動く。

```
let ika (a:int) b =
  a+b
```


### 関数呼び出し

関数呼び出しはF# スタイルで引数が足りない時は部分適用となる。

まずは基本的な呼び出しから。以下のhello関数の呼び出しに注目。

```
import "fmt"

let hello (msg: string) =
    GoEval "fmt.Println(msg)"

let main() =
    hello "hoge"
```

Folangとしては関数はカッコ無しで呼び出す。`hello "hoge"`の所。

これは以下のようなコードに展開される。

```golang
import "fmt"

func hello(msg string) {
    fmt.Println(msg)
}

func main() {
   hello("hoge")
}
```

### 関数呼び出しの部分適用

複数引数で部分適用すると以下。

```
let hello (msg: string) (num: int) = 
   GoEval "fmt.Printf(msg, num)"

let main () =
   let temp = hello "hoge%d"
   temp 123
```

以下の行で、2引数のhelloに1引数だけ渡している。

```
let temp = hello "hoge%d"
```

生成されるコードは以下。

```
temp := func(num int) { hello("hello%d", num) }
```

