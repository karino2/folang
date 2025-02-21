# 4. goで書いたラッパーを呼ぶ

前回: [3. スライスとパイプとMap](3_SlicePipeMap_ja.md)

今回はgoで書いた自分のコードをfolangから呼ぶ例を見ていきます。

## Folangはgoと分業して開発するのが基本

Folangは全ての機能を含む事を目指していません。

例えばループや変数への破壊的代入、メソッド呼び出し、ポインタなどがありません。

これらの機能を使う部分はgolangで書き、
それをfolangから呼ぶ、というのを基本的なスタイルとしています。

Folangにとっては、golangとの共同作業は例外では無くて基本となります。
ある程度以上のfolangプログラムではgoで書かれる部分が含まれる事でしょう。

ここではgolangとの共同作業をどう行うかを見ていきます。

## goのパッケージを直接呼ぶ例

第二回の[2. frtパッケージを使う (pkg_all.foiの説明)](2_UseFrtPackage_ja.md)でも少し触れましたが、
package_infoというものを書く事で、
fcに外で定義された関数の存在を知らせる事が出来ます。

基本的にfolangはフリースタンディング関数が多いライブラリはそのまま使える事が多く、
オブジェクト的なライブラリはラップして使う必要があります。

まずは良く使うものとして、golangのfilepathパッケージを使ってみます。
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

続きを書く。
