# Discriminated Unionの実装

例えば以下の簡単なケースを考えてみる。

```
type IntOrBool =
  |  I of int
  | B of bool
```

こうするとIとBという関数が出来て、結果はIntOrBool型で、実行時にどちらかがパターンマッチで判定出来る。
intとかboolのところは同じ型が来る場合もあるっぽい（上記のMicrosoft LearnのEquilateralTriangleとSquareの例参照）。
だから単なるtype assertionでは区別出来ない。

そこで、IntOrBoolをinterfaceとして、IntOrBool_I, IntOrBool_Bというstructを作る事にする。
以下のような実装。

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

## of無しのケース

```
type AorB =
  | A
  | B
```

のような事も出来る。この場合、Aは関数ではなく変数になる（引数無し関数と変数の区別がfsharpは無く、Unit引数の関数とは区別される）。

とりあえずGolang側は以下のようにvarにしてみる。

```golang
var New_AorB_A AorB = AorB_A{}
```

変数名にNewがついているのはおかしいが、あんまりofがある時と無い時でコードを変えたくないのでこうしておく。
どうせFolang上ではこの名前は出てこないしね。

## F#の関連リンク

そもそもDiscriminated Unionについての解説。

自分が馴染んでいるF#の実装の解説文書のリンクを貼っておく。

- [Discriminated Unions - F# for fun and profit](https://fsharpforfunandprofit.com/posts/discriminated-unions/) F#の機能としての説明
- [Discriminated Unions - F# - Microsoft Learn](https://learn.microsoft.com/en-us/dotnet/fsharp/language-reference/discriminated-unions)

## Genericsの実装

F# とは細かい所が違うので、比較の為にF# の挙動も載せておく。

### F# での挙動

F# などで以下のようなUnionのGenericsがあった時、

```fsharp
type Option<T> =
| Some of T
| None
```

Noneはグローバル変数のようなものになるのだが、Tの指定は無い。
以下のように変数にNoneを入れた場合、

```fsharp
let a = None
```

このaは`Option<'a>` という型になり、これは`Option<int>`型の変数にも`Option<string>`型の変数にも代入出来る。
aは全てになれる何かとなる。

### Folangでの実装

複数の具体的な型に同じ変数をassign出来る、というのはGolangで実現は難しいので、
FolangではもっとGolang的に自然な仕様にする。

GenericなUnion型の場合は値が無くても関数になる事にし、
type argumentであらわに指定する事にする。

```fsharp
type Option<T> =
| Some of T
| None

let a = None<int> () // Noneは関数

type AorB =
| A
| B

let b = B // Bはグローバル変数
```

なお、型推論が働くので型指定を明示的にしなくても解決出来る場合は指定の必要は無い。
