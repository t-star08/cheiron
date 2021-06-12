# Cheiron (ケイロン, ケイローン)

コードを Markdown などのファイルに挿入する  
`git commit` 時にコードには自動フォーマットがかかるが Markdown の pre タグにはかからない、といった場合に有用かもしれない


## v1.1.0 (2021/06/12)

主な変更点

- 挿入するコードを指定できるように改善
- 出力の分かりやすさ改善


## v1.0.0 (2021/05/14)

主な変更点

- 強制的に `$PROJECT_ROOT/$M` を対象にすることを廃止して、`$BRANCH` を使って、`$PROJCET_ROOT/$BRANCH` のように柔軟に指定できるように変更
- 正規表現ベース
- `cheiron project quiver` を追加


## 目次
- [Cheiron (ケイロン, ケイローン)](#cheiron-ケイロン-ケイローン)
  - [v1.1.0 (2021/06/12)](#v110-20210612)
  - [v1.0.0 (2021/05/14)](#v100-20210514)
  - [目次](#目次)
  - [インストール](#インストール)
  - [使い方](#使い方)
  - [前書き](#前書き)
  - [出来ること](#出来ること)
  - [cheiron insert について](#cheiron-insert-について)
  - [cheiron project について](#cheiron-project-について)
  - [cheiron project init について](#cheiron-project-init-について)
    - [cheiron project init -f](#cheiron-project-init--f)
  - [cheiron project arrow について](#cheiron-project-arrow-について)
  - [設定ファイル 【 arrow.json 】](#設定ファイル--arrowjson-)
    - [project_root](#project_root)
    - [branch](#branch)
    - [sources](#sources)
    - [insert_target](#insert_target)
    - [$CODE_DIRE の拡張](#code_dire-の拡張)
  - [cheiron project quiver について](#cheiron-project-quiver-について)
    - [設定ファイル 【 quiver.json 】](#設定ファイル--quiverjson-)
  - [挿入について](#挿入について)
    - [挿入するコードについて](#挿入するコードについて)
    - [挿入箇所](#挿入箇所)
    - [挿入時の安全性](#挿入時の安全性)
    - [挿入するコードの指定](#挿入するコードの指定)
  - [正規表現について](#正規表現について)
    - [正規表現の注意点](#正規表現の注意点)
    - [正規表現のショートカット](#正規表現のショートカット)
  - [フラグまとめ](#フラグまとめ)
    - [--keyword](#--keyword)
    - [-f, --force](#-f---force)
    - [-k, --key](#-k---key)
    - [-s, --simple](#-s---simple)
  - [ライブラリ](#ライブラリ)


## インストール

```sh
# $GOPATH/bin にインストールされる
$ go get -u github.com/t-star08/cheiron
```


## 使い方

```sh
# キーワードに Java を、コピー元に ./sample/Main.java を、挿入先に ./sample/sample.md を指定して実行
$ cheiron insert --keyword Java ./sample/Main.java ./sample/sample.md

# プロジェクトに対しての処理
# 設定ファイルの準備
# .cheiron が作成される
$ cheiron project init

# arrow.json を自分に合わせて編集する
$ vim .cheiron/arrow.json

# 挿入を実行する
$ cheiron project arrow hook1 hook2 ...
```


## 前書き

この README で使う表現についてのまとめ
- $PROJECT_ROOT : `arrow.json` の project_root で指定するパス
- $BRANCH : `arrow.json` の branch で指定するパス
- $CODE_DIRE : `arrow.json` の code_dire で指定するパス
- $CODE_FILE : `arrow.json` の code_file で指定するファイル
- $KEYWORD : `arrow.json` の keyword で指定する文字列
- $INSERT_TARGET : `arrow.json` の insert_target で指定するパス
- $M : 以下のいずれかに該当する文字が 1 字以上連続する名前のディレクトリ
  - アルファベット
  - アンダーバー
  - 数字
  - 垂直タブ以外の空白文字


## 出来ること

挿入先のファイルに pre タグで囲まれたコードを挿入すること  
- [cheiron insert](#cheiron-insert-について)
  - 1 つのコピー元のファイルと 1 つの挿入先のファイルを指定して挿入を実行する
- [cheiron project arrow](#cheiron-project-arrow-について)
  - `arrow.json` に基づいて挿入を実行する
- [cheiron project quiver](#cheiron-project-quiver-について)
  - `quiver.json` に基づいて複数のフォルダで `cheiron project arrow` を実行する


## cheiron insert について

`cheiron insert` は以下の 3 つを指定して使う
- $KEYWORD (詳しくは[こちら](#挿入箇所))
- コピー元のファイルへのパス (*.java など)
- 挿入先のファイルへのパス (*.md など)


## cheiron project について

`cheiron project` は以下のように使う
- `cheiron project init` (詳しくは[こちら](#cheiron-project-init-について))
- `cheiron project arrow hook1 hook2 ...` (詳しくは[こちら](#cheiron-project-arrow-について))
- `cheiron project quiver hook1 hook2 ...` (詳しくは[こちら](#cheiron-project-quiver-について))


## cheiron project init について

[.cheiron] とその中に設定ファイル [arrow.json] と [quiver.json] を作成する  
`arrow.json` については[こちら](#設定ファイル-【-arrow.json-】)  
`quiver.json` については[こちら](#設定ファイル-【-quiver.json-】)

```sh
$ cheiron project init
Created directory [./.cheiron]
By Editing json, you can add hook
$ cat .cheiron/arrow.json
{
    "project_root": "./",
    "branch": ".*",
    "sources": {
        "python3": {
            "code_dire": "code_python3",
            "code_file": "main.py",
            "keyword": "Python3"
        },
        "java": {
            "code_dire": "code_java",
            "code_file": "Main.java",
            "keyword": "Java"
        },
        "cpp": {
            "code_dire": "code_c-plus-plus",
            "code_file": "main.cpp",
            "keyword": "C++"
        },
        "cc": {
            "code_dire": "code_c-plus-plus",
            "code_file": "main.cc",
            "keyword": "C++"
        }
    },
    "insert_target": "DEFAULT.md"
}

$ cat .cheiron/quiver.json
{ "specify": [], "ignore": [] }

```
※ `arrow.json` はデフォルトでは改行なし 1 行の json ファイルだが、見やすさの観点から整形している


### cheiron project init -f

`.cheiron` を上書きする


## cheiron project arrow について

- `cheiron project arrow hook1 hook2 ...` のように使う
  - hook は `arrow.json` の sources 下にあるオブジェクトのキーを指定する
  - 例えば、[この json](#cheiron-project-init-について) では python3 や java などを hook として指定する
- `$PROJECT_ROOT/$BRANCH/$CODE_DIRE/CODE_FILE` のコードを `$PROJECT_ROOT/$BRANCH/$INSERT_TARGET` 内に挿入する


## 設定ファイル 【 arrow.json 】

`arrow.json` を編集することで [cheiron project arrow](#cheiron-project-arrow-について) の処理対象を設定できる

`arrow.json` は主に以下の 4 要素に分けられる
- [project_root](#project_root)
- [branch](#branch)
- [sources](#sources)
- [insert_target](#insert_target)


### project_root

ここで指定するディレクトリについて[処理](#挿入について)を行う  
パスは相対パスでも絶対パスでもよい


### branch

挿入を実行する単位となるパスを指定する  

挿入の単位とは少し詳しく書くと、  
`$PROJECT_ROOT/$BRANCH` 内で `$CODE_DIRE/$CODE_FILE` を探して、  
`$PROJECT_ROOT/$BRANCH` 内の `INSERT_TARGET` に挿入を実行する

正規表現として処理される (正規表現についての注意点は[こちら](#正規表現についての注意点))


### sources

sources は任意の文字列をキーとして以下の構造を持つオブジェクトを任意の数持っている  
[このオブジェクトのキーを `cheiron project arrow` の引数として使うことができる](#cheiron-project-arrow-について)
- code_dire
  - コピー元のファイルがあるディレクトリ名
  - `$PROJECT_ROOT/$BRANCH/` 下にあるディレクトリへのパスを指定する
  - 正規表現として処理される (正規表現についての注意点は[こちら](#正規表現についての注意点))
- code_file
  - コピー元のファイル名
  - `$PROJECT_ROOT/$BRANCH/$CODE_DIRE` 下にあるファイル名を指定する
  - 正規表現として処理される (正規表現についての注意点は[こちら](#正規表現についての注意点))
- keyword
  - 挿入先のファイルのどの行に挿入するかを決める(詳しく[こちら](#挿入について))


### insert_target

insert_target はコードを挿入するファイルへのパスを指定する  

正規表現として処理される (正規表現についての注意点は[こちら](#正規表現についての注意点))


### $CODE_DIRE の拡張

`cheiron project arrow` を実行したとき探すディレクトリは以下の 2 通り

- `$PROJECT_ROOT/$BRANCH/$CODE_DIRE`
- `$PROJECT_ROOT/$BRANCH/$CODE_DIRE_2` 
- `$PROJECT_ROOT/$BRANCH/$CODE_DIRE_i` (i は 3 以上の整数)
   - `$PROJECT_ROOT/$BRANCH/$CODE_DIRE_(i-1)` が存在すれば  `$PROJECT_ROOT/$BRANCH/$CODE_DIRE_i` を探す
   - つまり、 `$PROJECT_ROOT/$BRANCH/$CODE_DIRE_3` は  `$PROJECT_ROOT/$BRANCH/$CODE_DIRE_2` があれば探す

※ [フラグ](#-s---simple)で `$PROJECT_ROOT/$BRANCH/$CODE_DIRE_2` 以降を探さないこともできる


## cheiron project quiver について

`cheiron project quiver java python3` のように使うと、複数フォルダで `cheiron project arrow java python3` を実行できる

対象となるフォルダはデフォルトでは、実行ディレクトリ直下の `$M` に該当するディレクトリ  
※ [`quiver.json` でディレクトリ名を指定して実行することもできる](#設定ファイル-【-quiver.json-】)


### 設定ファイル 【 quiver.json 】

`quiver.json` を変更することで、`cheiron project quiver` の対象となるフォルダを変更できる

quiver.json の要素は以下の 2 つ
- specify
  - ディレクトリ名を指定する
  - この配列を空にすると、`$M` にマッチするすべてのディレクトリが対象となる
  - 正規表現に対応していない
- ignore
  - 処理の対象にしないディレクトリを指定する
  - 正規表現に対応していない


## 挿入について


### 挿入するコードについて

挿入時はコードに以下の処理を施す (※ コピー元のファイルの改変は一切行わない)

- コードの先頭に `<pre lang="$KEYWORD">` を付ける
- コードの末尾に `</pre>` を付ける
- 「&lt;」を `&lt;` に変換する
  - ただし、&lt; の後ろの文字が記号の場合は以下の場合を除いて変換しない
    - 「!」
    - 「?」


### 挿入箇所

挿入箇所は以下のルールによって決められる  
`arrow.json` で設定する keyword が大事になってくる

- pre タグのエリア
  - `<pre lang="$KEYWORD">` から `</pre>` まで
    - 例えば `<pre lang="Java">` から `</pre>` まで
  - 必ず上記の形式でなければならない
  - `<pre lang="$KEYWORD">` から `</pre>` までを削除して挿入を行う
  - つまり、pre タグのエリアは上書きされる
- `\$KEYWORD` となっている行
  - 例えば `\Java` と書かれている行
  - `\$KEYWORD` の行を削除して挿入を行う

2 つのルールの優先度はデフォルトでは以下の通り
1. pre タグのエリア
2. `\$KEYWORD` となっている行

「1.」で挿入箇所が 1 つも見つからなかったとき、「2.」で挿入箇所を探す  

※ [フラグ](#-k---key)で優先度を逆転させることが可能


### 挿入時の安全性

以下の条項に当てはまるとき、挿入は実行されない

- コピー元のファイルの数と挿入箇所の数が合わない

※ [フラグ](#-f---force)で強制実行させることも可能


### 挿入するコードの指定

`\$KEYWORD` で挿入箇所を指定する際、挿入するコードを指定できる

- 指定方法
  - `\$KEYWORD` に続けて以下の 2 通りの方法で記述する
    - `\$CODE_FILE があるディレクトリ\$CODE_FILE`
    - `\$CODE_FILE があるディレクトリ/$CODE_FILE`
- 挿入
  - `$CODE_FILE があるディレクトリ/CODE_FILE` に該当するコードが挿入される
  - ただ、同一ブランチ内に `$CODE_FILE があるディレクトリ/CODE_FILE` に該当するコードが 2 つ以上ある場合、動作の保証はされない
- 注意
  - 「挿入するコードの指定」を行った場合、`\$KEYWORD` で決まる挿入箇所は指定がある箇所のみとなる
  - つまり、単に `\$KEYWORD` となっている箇所は挿入対象にならない
- 例
  - `\Java\code_java\Main.java`
    - `$PROJECT_ROOT/$BRANCH/.../code_java/Main.java` が `\Java\code_java\Main.java` が書かれている行に挿入される
  - `\Python3\code_python3\main.py`
    - `$PROJECT_ROOT/$BRANCH/.../code_python3/main.py` が `\Python3\code_python3\main.py` が書かれている行に挿入される


## 正規表現について


### 正規表現の注意点

このコマンドでは以下の点に注意が必要
1. 正規表現は完全一致ではなく、部分一致
1. `path/.+/to` のような正規表現をした際、"path" 下に文字列としての ".+" という名前のディレクトリがある場合、".+" は正規表現として認識されない


### 正規表現のショートカット

よく使いそうな正規表現のショートカットを用意している

- \$A : ^[\w\s]+$
- \$B : ^[\w]+$
- \$C : .*
- \$D : ^[0-9]
- \$E : .

`path/$A/to` のように指定することで、`$A` の箇所は正規表現 `^[\w\s]+$` と解釈して探索を行う  

※ ただし、"$A" というディレクトリ名やファイル名がある場合、そちらが優先される


## フラグまとめ

`cheiron insert` と `cheiron project arrow` で使えるフラグのまとめ


### --keyword

- 概要 : `cheiron project arrow` を行うときの $KEYWORD に当たるものを指定する
- 注意 : これは `cheiron insert` を実行する際に必ず必要
- 対応しているコマンド
  - `cheiron insert`

### -f, --force

- 概要 : [挿入時の安全性](#挿入時の安全性)の条項に引っかかっていても強制的に挿入を実行する  
なお、この場合は以下のように場合分けされる
  - コピー元のコードの数 > 挿入箇所
    - 挿入箇所の数だけ挿入する
    - 挿入の順番は見つかった順
  - コピー元のコードの数 < 挿入箇所
    - コピー元のコードの数だけ挿入する
    - 挿入の順番は見つかった順
- 対応しているコマンド
  - `cheiron insert`
  - `cheiron project arrow`


### -k, --key

- 概要 : [挿入箇所](#挿入箇所)を決定するときに `$\KEYWORD` となっている行優先で行う
- 対応しているコマンド
  - `cheiron insert`
  - `cheiron project arrow`


### -s, --simple

- 概要 : [$CODE_DIRE の拡張](#code_dire-の拡張)を行わない
- 対応しているコマンド
  - `cheiron project arrow`


## ライブラリ

- [cobra](https://github.com/spf13/cobra)
