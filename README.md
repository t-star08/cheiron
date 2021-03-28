# Cheiron (ケイロン, ケイローン)

コードを Markdown などのファイルに挿入する  
`git commit` 時にコードには自動フォーマットがかかるが Markdown の pre タグにはかからない、といった場合に有用かもしれない


## 目次
- [Cheiron (ケイロン, ケイローン)](#cheiron-ケイロン-ケイローン)
  - [目次](#目次)
  - [インストール](#インストール)
  - [使い方](#使い方)
  - [前書き](#前書き)
  - [出来ること](#出来ること)
  - [cheiron insert について](#cheiron-insert-について)
  - [cheiron project について](#cheiron-project-について)
  - [cheiron project init について](#cheiron-project-init-について)
    - [cheiron project init -f](#cheiron-project-init--f)
  - [arrow.json について](#arrowjson-について)
    - [project_root](#project_root)
    - [sources](#sources)
  - [cheiron project arrow について](#cheiron-project-arrow-について)
  - [$CODE_DIRE の拡張](#code_dire-の拡張)
  - [挿入について](#挿入について)
    - [挿入箇所](#挿入箇所)
    - [挿入するコードについて](#挿入するコードについて)
    - [挿入時の安全性](#挿入時の安全性)
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
# ../cheiron_settings/arrow.json が作成される
$ cheiron project init

# arrow.json を自分に合わせて編集する
$ vim ../cheiron_settings/arrow.json

# 挿入を実行する
$ cheiron project arrow hook1 hook2 ...
```


## 前書き

この先使う表現についてのまとめ
- $PROJECT_ROOT : `arrow.json` の project_root で指定するパス
- $CODE_DIRE : `arrow.json` の code_dire で指定するパス
- $CODE_FILE : `arrow.json` の code_file で指定するファイル
- $KEYWORD : `arrow.json` の keyword で指定する文字列
- $INSERT_TARGET : `arrow.json` の insert_target で指定するファイル
- $M : 以下のいずれかに該当する文字が 1 字以上連続する名前のディレクトリ
  - アルファベット
  - アンダーバー
  - 数字
  - 垂直タブ以外の空白文字


## 出来ること

出来ることは挿入先のファイルに pre タグで囲まれたコードを挿入すること  
- [cheiron insert](#cheiron-insert-について)
  - 1 つのコピー元のファイルと 1 つの挿入先のファイルを指定して挿入を実行する
- [cheiron project arrow](#cheiron-project-arrow-について)
  - `arrow.json` に基づいて挿入を実行する
  - 該当するディレクト全てにおいて n 個のコピー元のコードを 1 つの挿入先のファイルに挿入する(詳しくは[ここ](#cheiron-project-arrow-について)で)


## cheiron insert について

`cheiron insert` は以下の 3 つを指定して使う
- $KEYWORD (詳しくは[こちら](#挿入箇所))
- コピー元のファイルへのパス (*.java など)
- 挿入先のファイルへのパス (*.md など)


## cheiron project について

`cheiron project` は以下のように使う
- `cheiron project init` (詳しくは[こちら](#cheiron-project-init-について))
- `cheiron project arrow lang1 lang2 ...` (詳しくは[こちら](#cheiron-project-arrow-について))


## cheiron project init について

設定ファイル (../cheiron_settings/arrow.json) を作成する  
`arrow.json` については[こちら](#arrow.json-について)

```sh
$ cheiron project init
Created directory [../cheiron_settings/arrow.json]
You can edit it to add target language
$ cat ../cheiron_settings/arrow.json
{
    "project_root": "./",
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
```
※ `arrow.json` はデフォルトでは改行なし 1 行の json ファイルだが、見やすさの観点から整形している


### cheiron project init -f

`../cheiron_settings/arrow.json` を上書きする


## arrow.json について

`arrow.json` を編集することで [cheiron project arrow](#cheiron-project-arrow-について) の処理対象を設定できる
`arrow.json` は主に以下の 3 要素に分けられる
- [project_root](#project_root)
- [sources](#sources)
- [insert_targt](#insert_target)


### project_root

ここで指定するディレクトリについて[処理](#挿入について)を行う  
パスは相対パスでも絶対パスでもよい


### sources

sources は任意の文字列をキーとして以下の構造を持つオブジェクトを任意の数持っている
- code_dire
  - コピー元のファイルがあるディレクトリ名
  - `$PROJECT_ROOT/$M/` 下にあるディレクトリ名を指定する
- code_file
  - コピー元のファイル名
  - `$PROJECT_ROOT/$M/$CODE_DIRE` 下にあるファイル名を指定する
- keyword
  - 挿入先のファイルのどの行に挿入するかを決める(詳しく[こちら](#挿入について))


## cheiron project arrow について

主な点は以下の 2 点
- `cheiron project arrow hook1 hook2 ...` のように使う
  - hook は `arrow.json` の sources 下にあるオブジェクトのキーを指定する
  - 例えば、[この](#cheiron-project-init-について) json では python3 や java などを hook として指定する
- `$PROJECT_ROOT/$M/$CODE_DIRE/CODE_FILE` のコードを `$PROJECT_ROOT/$M/$INSERT_TARGET` 内に挿入する


## $CODE_DIRE の拡張

`cheiron project arrow` を実行したとき探すディレクトリは以下の 2 通り

- `$PROJECT_ROOT/$M/$CODE_DIRE`
- `$PROJECT_ROOT/$M/$CODE_DIRE_2` 
- `$PROJECT_ROOT/$M/$CODE_DIRE_i` (i は 3 以上の整数)
   - `$PROJECT_ROOT/$M/$CODE_DIRE_(i-1)` が存在すれば  `$PROJECT_ROOT/$M/$CODE_DIRE_i` を探す
   - つまり、 `$PROJECT_ROOT/$M/$CODE_DIRE_3` は  `$PROJECT_ROOT/$M/$CODE_DIRE_2` があれば探す


## 挿入について

### 挿入箇所

挿入箇所は以下のルールによって決められる  
`arrow.json` で設定する keyword が大事になってくる

- pre タグのエリア
  - `<pre lang="$KEYWORD">` から `</pre>` まで
    - 例えば `<pre lang="Java">` から `</pre>` まで
  - 必ず上記の形式でなければならない
  - `<pre lang="$KEYWORD">` から `</pre>` までを一旦削除して挿入を行う
  - つまり、pre タグのエリアは上書きされる
- `\$KEYWORD` となっている行
  - 例えば `\Java` と書かれている行
  - 1 行に `\$KEYWORD` だけ書かれている行が対象
  - `\$KEYWORD` を削除して挿入を行う

2 つのルールの優先度はデフォルトでは以下の通り
1. pre タグのエリア
2. `\$KEYWORD` となっている行

「1.」で挿入箇所が 1 つも見つからなかったとき、「2.」で挿入箇所を探す  


### 挿入するコードについて

挿入時はコードに以下の処理を施す(※ コピー元のファイルの改変は一切行わない)

- コードの先頭に `<pre lang="$KEYWORD">` を付ける
- コードの末尾に `</pre>` を付ける
- 「&lt;」を `&lt;` に変換する
  - ただし、&lt; の後ろの文字が記号の場合は以下の場合を除いて変換しない
    - 「!」
    - 「?」

### 挿入時の安全性

以下の条項に当てはまるとき、挿入は実行されない

- コピー元のファイルの数と挿入箇所の数が合わない


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
