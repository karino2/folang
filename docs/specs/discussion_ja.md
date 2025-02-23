# 仕様検討時のメモ

仕様を検討してた時のメモなどをまとめておく。
決定した仕様は[note_ja.md](note_ja.md)の方に書く。

あくまでメモなのであまり体裁は気にしない。

## スライス expressionの区切り

ReScriptはカンマ、FSharpはセミコロン。

- [Array & List - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/array-and-list)
   - [Record - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/record) レコードもカンマ
   - [Tuple - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/tuple) タプルはカッコをつけてカンマ
- FSharpはセミコロン [Lists - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/lists)
   - レコードもセミコロン
   - タプルはカッコをつけてカンマ [Tuples - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/tuples)

レコードの区切りをセミコロンにしているので、セミコロンに統一する事に。

F＃久しぶりに触るとカンマにしてしまいがちだが、レコードのフィールドがタプルの時はやっぱりややこしいのでセミコロンの方がいい気もするな。

スライスの区切りはセミコロンに決定。

## 関数呼び出しの仕様検討

とりあえずhello world的なものを考えて仕様を考えていきたい。

まずgolangの関数呼び出しはF# のdotnetの関数呼び出しのようにしたい気がする。

```
import "fmt"

let main () =
   fmt.Println("hoge")
```

これは以下に展開されて欲しい。

```golang
import "fmt"

func main() {
   fmt.Println("hoge")
}
```

次にfolangの関数定義を考える。
最初はtype inferenceは無い状態から始めよう。すると関数定義は以下か。

```
import "fmt"

let hello (msg: string) =
    fmt.Println(msg)

let main() =
    hello "hoge"
```

folangとしては関数はカッコ無しで呼び出し、部分適用されていくFSharp的な実行でいいだろう。
このカッコで呼び出すかそのまま呼び出すかでどちらの呼び出しかを分ける感じにしたい。

これはどういうコードに展開されるかはちょっと現時点ではよくわからないな。

```golang
import "fmt"

func hello(msg string) {
    fmt.Println(msg)
}

func main() {
   hello("hoge")
}
```

このhello("hoge")に展開されるのか部分適用されるのかはコンパイル時にたぶん決定出来るよな。

複数引数だと以下みたいな感じになるか。

```
let hello (msg: string) (num: int) = 
   fmt.Printf(msg, num)

let main () =
   let temp = hello "hoge%d"
   temp 123
```

部分適用すると以下みたいか？

```golang
func hello(msg string, num int) {
    fmt.Printf(msg, num)
}

func main() {
   temp := func(num int) { hello("hello%d", num) }
   temp(123)
}
```

とりあえずこのくらいを生成出来るようにする所から始めるか。

いや、関数定義はもっとgolang的でいいのではないか？

```
import "fmt"

func hello (msg string, num int) = 
   fmt.Printf(msg, num)

func main () =
   let temp = hello "hoge%d"
   temp 123
```

いや、やはり関数定義のシンタックスが面倒なのは良くないな。



## Genericsのシンタックス（タイプパラメータ）

goは大括弧だがFSharpは角括弧だ。

なるべくならgolangに揃えたいが、FSharpはシンタックス的に変数と関数の区別が曖昧になるようになっているので、
インデックスアクセスとややこしい事になる。

```
let Length[T any] (args: []T) =
   ...

Length[int] listOfList[3]
```

やはりこれは厳しいな。角括弧にするか。

```
let Length<T any> (args: []T) =
   ...

Length<int> listOfList[3]
```

これはこれでかなりパースがトリッキーになるんだが、仕方なし。

## 外部の型情報

パッケージアクセスをそろそろ考えたい。型情報ファイルを手作業で作って加えるのでいいのだが、どういうシンタックスにするか。

### 関連情報

FSharpのシグニチャファイルは似たような話だ。

- [Signature files - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/signature-files)

関数はvalになる。うーん。このためだけにvalというのもなぁ。moduleとかnamespaceは通常の定義と同じ。

ReScriptではletで定義している。

- [Module - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/module#signatures)

ただシンタックスはコロンになる。気持ち悪さはあるが、変数が定義されていると思えばletが正しい気もする。
module typeが先頭に来る所が通常のmoduleとは違う。

Borgoは関数定義がキーワードであるので普通に関数定義のようになっている。

- [Borgo Programming Language](https://borgo-lang.github.io/#package-definitions)

パーサーは同じのが使えるというメリットはあるが、まぁこれは言語のシンタックスが違いすぎるのであまり参考にはならないか。

### folangでの外部の型情報

さて、folangではどうしよう？

モジュールという名前にしておくか？いや、変えておく方がいいか。
golang的にはpackageなんだよなぁ。

package_infoにしよう。
そしてReScriptを真似してletでコロンとしてみるか。

```
package_info slice =
   let Length<T>: []T -> int
   let Take<T> : int->[]T->[]T 
```

sliceは開幕からgenericsが要るやん... LengthはanyでもいいがTakeは要るよなぁ。仕方ない、諦めて対応しよう。
別にしばらくは明示的に呼び出し時に指定するでもいいので。

type parameterを`<T>`にするか`[T]`にするかは悩ましいが、インデックスアクセスが大括弧なのでFSharpに揃える。

このままBufferも書いてみよう。

```
package_info buf =
    type Buffer
    let New: ()->Buffer
    let Write: Buffer->()
    let String: ()->string
```

こっちはいい気がするな。typeは右に書くものが無いな。type aliasを中で定義するようなのはサポートしなくていいだろう。
コンストラクタとかはどうせラップするのだからNew関数を作らせる事にする。

ファイルの拡張子はとりあえずfoiにしておくかな。

## コメント

コメントはgolangと同様、CスタイルとC++の一行コメントの２つをサポートする。（ただし現状はCスタイルのみ実装）

## Discriminated Unionの実装方針

まずは自分が馴染んでいるF#の実装の解説文書のリンクから。

- [Discriminated Unions - F# for fun and profit](https://fsharpforfunandprofit.com/posts/discriminated-unions/) F#の機能としての説明
- [Discriminated Unions - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/discriminated-unions)

例えば以下の簡単なケースを考えてみる。

```
type IntOrBool =
  |  I of int
  | B of bool
```

こうするとIとBという関数が出来て、結果はIntOrBool型で、実行時にどちらかがパターンマッチで判定出来る。
intとかboolのところは同じ型が来る場合もあるっぽい（上記のMicrosoft LearnのEquilateralTriangleとSquareの例参照）。
だから単なるtype assertionでは区別出来ない。

以下のような実装だとどうだろう？

```golang
type IntOrBool interface {
  IntOrBool_Union()
}

func (IntOrBool_I) IntOrBool_Union(){}
func (IntOrBool_B) IntOrBool_Union(){}

type IntOrBool_I struct {
   Value int
}

type IntOrBool_B struct {
   Value bool
}

func New_IntOrBool_I(v int) IntOrBool { return IntOrBool_I{v} }
func New_IntOrBool_B(v bool) IntOrBool { return IntOrBool_B{v} }
```

IかBはNewXXXの関数呼び出しにマップすれば良さそう。
最初はポインタにしていたが、interfaceとstructを内部で区別するのが定義があとに来るケースでは困難なので全部実体に統一。

これならIntOrBoolはtype assertionで実行時にIかBは区別出来るんじゃないか？
試してみよう。

```fsharp
match iob with
| I ival -> printfn "i=%d" ival
| B bval -> printfn "b=%v" bval
```

この単純なケースなら単なるtype assertで実現出来そうだな。

```golang
switch iob.(type) {
case IntOrBool_I:
   ival := iob.Value
   fmt.Printf("i=%d", ival)
case IntOrBool_B:
   bval := iob.Value
   fmt.Printf("b=%v", bval)
}
```

もちろん実際はもっと複雑なパターンがありうるのでtype switchで書けるのか、という問題はあるが、たぶんcaseの中にさらなる条件で全部書けるはずか？
まぁ複雑なパターンはしばらく使わないので、まずはこの単純なケースが動くようにすべきか。

### of無しのケース

```
type AorB =
  | A
  | B
```

のような事も出来る。この場合、Aは関数ではなく変数になる（引数無し関数と変数の区別がfsharpは無く、Unit引数の関数とは区別される）。

とりあえずgolang側は以下のようにvarにしてみる。

```golang
var New_AorB_A AorB = AorB_A{}
```

変数名にNewがついているのはおかしいが、あんまりofがある時と無い時でコードを変えたくないのでこうしておく。
どうせfolang上ではこの名前は出てこないしね。

### UnionのaltJSの実装を見てみる。

- [Fable · Features](https://fable.io/docs/typescript/features.html#f-unions) FableのUnionの実装
- [Pattern Matching / Destructuring - ReScript Language Manual](https://rescript-lang.org/docs/manual/v10.0.0/pattern-matching-destructuring) ReScriptのUnion実装、payloadのあたりが参考になる。

fableでもReScriptでも、まず各caseを表すタグを用意してここに文字列かインデックスを入れている。
JSだと実行時に型情報が無くなるので必要な気がするが、goならいらないのでは？

Borgoという言語にはRustのenumみたいなのがあるので見てみる。 [borgo/compiler/test/snapshot/codegen-emit/enums.exp at main · borgo-lang/borgo](https://github.com/borgo-lang/borgo/blob/main/compiler/test/snapshot/codegen-emit/enums.exp)

タグを使っているなぁ。


