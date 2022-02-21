# Cheiron (ケイロン)

テキストを挿入する  
- 分散させて保存したテキストを結合する
- 決まった形式のファイルを生成する

## v2.0.0 (2022/02/21)

主な変更点

- 挿入するテキストを、挿入を受けるファイルから指定するように変更

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
- [Cheiron (ケイロン)](#cheiron-ケイロン)
  - [v2.0.0 (2022/02/21)](#v200-20220221)
  - [v1.1.0 (2021/06/12)](#v110-20210612)
  - [v1.0.0 (2021/05/14)](#v100-20210514)
  - [目次](#目次)
  - [インストール](#インストール)
  - [使い方](#使い方)
  - [cheiron init](#cheiron-init)
    - [-f, --force](#-f---force)
    - [-a, --any](#-a---any)
  - [cheiron status](#cheiron-status)
  - [設定ファイル 【 cheiron.json 】](#設定ファイル--cheironjson-)
  - [挿入箇所の書式 (および cheiron arrow の処理対象)](#挿入箇所の書式-および-cheiron-arrow-の処理対象)
    - [cheiron arrow single](#cheiron-arrow-single)
    - [cheiron arrow multi](#cheiron-arrow-multi)
    - [cheiron arrow routine](#cheiron-arrow-routine)
  - [cheiron arrow ??? のオプション](#cheiron-arrow--のオプション)
    - [-o, --overwrite](#-o---overwrite)
    - [-p, --practice](#-p---practice)
    - [-q, --quiet](#-q---quiet)
    - [-q, --quiet](#-q---quiet-1)
  - [実行結果を記録する json について](#実行結果を記録する-json-について)
  - [ライブラリ](#ライブラリ)


## インストール

```sh
# $GOPATH/bin にインストールされる
$ go get -u github.com/t-star08/cheiron
```


## 使い方

```sh
# 設定ファイルおよび、それらを保存するディレクトリを作成
$ cheiron init

# cheiron.json を自分に合わせて編集する
$ vim .cheiron/cheiron.json

# template を元にファイルを生成
$ cheiron arrow routine strategy1 strategy2 ...

# 挿入ファイルを指定して挿入
$ cheiron arrow single --target path/to/target strategy1 strategy2 ...

# 複数のファイルに挿入 (cheiron.json による)
$ cheiron arrow multi strategy1 strategy2 ...

```


## cheiron init

```sh
$ cheiron project init
Created ".cheiron/cheiron.json"
By Editing json, you can set config

$ ls .cheiron
cheiron.json  history

```

- cheiron.json: 設定ファイル
- history: 実行時の詳細な情報の json ファイルが保存される

### -f, --force

`.cheiron` を上書きする

### -a, --any

生成するファイルのなかから足りないものを生成する


## cheiron status

```sh
cheiron status
CONFIG: path/to/cheiron.json
HISTORY: path/to/history
LAST USED: datetime when last used
```

- `.cheiron` を、ファイル階層を 5 つまで遡って探し、その情報を表示
- CONFIG: `cheiron.json` へのパス
- HISTORY: `history` へのパス
- LAST USED: 最後に使った日時


## 設定ファイル 【 cheiron.json 】

```sh
$ cat .cheiron/cheiron.json
{
  "projectRoot": ".",
  "branch": ".*",
  "ignore": [
    "branch path written here be ignored"
  ],
  "insertTarget": "DEFAULT.md",
  "strategies": {
    "strategy1": {
      "useEscapeOption": false,
      "usePreLangOption": false,
      "targetSuffixes": [
        ".py",
        ".ruby"
      ]
    },
    "strategy2": {
      "useEscapeOption": true,
      "usePreLangOption": false,
      "targetSuffixes": [
        ".java",
        ".cpp",
        ".cc"
      ]
    },
    "strategy3": {
      "useEscapeOption": true,
      "usePreLangOption": true,
      "targetSuffixes": [
        ".*"
      ]
    }
  },
  "strategyAliases": {
    "aliase": [
      "strategy comb"
    ],
    "aliase1": [
      "stratgey1",
      "strategy2"
    ],
    "aliase2": [
      "strategy1",
      "strategy3"
    ]
  },
  "preLangSuffixes": {
    "suffix": "language"
    ".c": "C",
    ".cc": "C++",
    ".cpp": "C++",
    ".go": "GO",
    ".java": "Java",
    ".py": "Python3",
    ".ruby": "Ruby",
  },
  "routine": [
    {
      "template": "Path/to/template",
      "priority": 0
    },
    {
      "template": "template.md",
      "priority": 1
    }
  ]
}
```

- `cheiron.json` を編集することで `cheiron arrow...` の処理対象を設定できる
- 【 projectRoot 】
  - コマンドを実行するディレクトリから `branch` までのパス
- 【 branch 】
  - 挿入する単位 (`projectRoot/branch/`)
  - 正規表現として解釈される
    - 正規表現なため、複数の単位で挿入できる
- 【 ignore 】
  - `branch` のなかで無視するディレクトリ (`projectRoot/ignore/`)
- 【 insertTarget 】
  - `projectRoot/branch/insertTarget` を挿入の対象にする
- 【 strategy 】
  - 挿入時のオプション
  - コマンド実行時に引数として指定
    - 複数可
    - 指定順に優先される
  - 挿入するテキストを...
    - `useEscapeOption`: `<` を `&lt;` で置換
    - `usePreLangOption`: `<pre lang="%s">` `</pre>` で挟む
  - `targetSuffixes`: 上記のオプションを適応する拡張子
    - 特に指定しない (すべてのファイルに適応する)とき、`.*` か `*` を書く
- 【 strategyAliases 】
  - `strategy` の組み合わせ
  - コマンド実行時に引数として `strategy` を複数指定する代わりに `strategyAlias` を指定できる
- 【 preLangSuffixes 】
  - `<pre lang="%s">` の `%s` 
  - 拡張子を key に、`%s` を value に
- 【 routine 】
  - `cheiron arrow routine` を実行するときに使う
  - `template`: 生成するファイルの型となるファイルへの `.cheiron` からのパス
  - `priority`: `template` の優先順位


## 挿入箇所の書式 (および cheiron arrow の処理対象)

path: `projectRoot/branch/insertTarget`

```txt
<<< source/sample.txt

<<< source/sample.txt?

<<< source/sample.txt? テスト

```

- `<<< source/sample.txt`
  - `projectRoot/branch/source/sample.txt` のテキストで置換する
  - `projectRoot/branch/source/sample.txt` がないとエラー
- `<<< source/sample.txt?`
  - `projectRoot/branch/source/sample.txt` のテキストで置換する
  - `projectRoot/branch/source/sample.txt` がなくてもエラーにならない
- `<<< source/sample.txt? テスト`
  - 挿入時 `テスト` を保護する
  - 次の行に `projectRoot/branch/source/sample.txt` のテキストを挿入する
  - `projectRoot/branch/source/sample.txt` がなくてもエラーにならない
  - `projectRoot/branch/source/sample.txt` がないとき、`テスト` は消される


### cheiron arrow single

```sh
$ cheiron arrow single --target path/to/target strategy...
```

- `.cheiron` を、ファイル階層を 5 つまで遡って探し、最初に見つかった `.cheiron/cheiron.json` の設定をもとに実行する
- `--target` または `-t`:
  - `cheiron.json` での `projectRoot/branch/insertTarget` を指定
  - `insertTarget` が属するディレクトリを `branch` として処理する
  - (`projectRoot` は `.` になる)
  - 正規表現不可
- もとの `projectRoot/branch/insertTarget` は `projectRoot/branch/.>>>[insertTarget` に移動させられる
- 他は `cheiron.json` に従う


### cheiron arrow multi

```sh
$ cheiron arrow multi strategy...
```

- `.cheiron` を、ファイル階層を 5 つまで遡って探し、最初に見つかった `.cheiron/cheiron.json` の設定をもとに実行する
- `cheiron.json` に従って、`projectRoot/branch` にマッチするディレクトリで挿入を実行する
- もとの `projectRoot/branch/insertTarget` は `projectRoot/branch/.>>>[insertTarget` に移動させられる


### cheiron arrow routine

```sh
$ cheiron arrow routine strategy...
```

- `.cheiron` を、ファイル階層を 5 つまで遡って探し、最初に見つかった `.cheiron/cheiron.json` の設定をもとに実行する
- `cheiron.json` に従って、`projectRoot/branch` にマッチするディレクトリで次の処理を実行する
- 【 処理 】
  - `cheiron.json` の `routine` で指定する `template` のなかから、各ブランチで挿入時にエラーが発生しない `template` を探す
  - 各 `branch` で `template` の挿入要件が満たされるように挿入をおこない、`projectRoot/branch/insertTarget` の場所に新しくファイルを生成する
  - すでに `projectRoot/branch/insertTarget` がある場合、そのファイルを `projectRoot/branch/.>>>[insertTarget` として移動する


## cheiron arrow ??? のオプション

`cheiron arrow` の `single`, `multi`, `routine` では以下のオプションを使うことができる (`--help` は省略)

### -o, --overwrite

`projectRoot/branch/insertTarget` を `projectRoot/branch/.>>>[insertTarget` として移動せず、上書きする


### -p, --practice

もし `cheiron arrow...` を実行したとき、どのディレクトリで挿入が実行されるか、またはなぜ実行されないか、などの情報を挿入せずに確認する


### -q, --quiet

実行結果のメッセージおよび、実行結果の詳細を記録する json を生成しない


### -q, --quiet

実行結果のメッセージは出力するが、実行結果の詳細を記録する json を生成しない


## 実行結果を記録する json について

```sh
$ cat .cheiron/history/datetime.json
{
  "performance": false,
  "recognized": [
    {
      "pathToBranch": "path/to/branch",
      "met": [
        {
          "pathToSource": "path/to/source",
          "Usedstrategy": "used strategy",
          "cuz": "-"
        }
      ],
      "unMet": [
        {
          "pathToSource": "path/to/source",
          "Usedstrategy": "-",
          "cuz": "why skipped"
        }
      ],
      "usedTemplate": "used template",
      "cuz": "-"
    }
  ],
  "ignored": [
    {
      "pathToBranch": "path/to/branch",
      "met": null,
      "unMet": null,
      "usedTemplate": "",
      "cuz": "why ignored"
    },
  ]
}

```

詳細な実行結果が記録された json
- performance: `practice` ではなかったかどうか (`practice == true` だったら `false`)
- recognized: 挿入対象になった `branch`
  - pathToBranch: `branch` へのパス
  - met: 挿入されたテキスト
    - pathToSource: テキストへのパス
    - usedStrategy: 適応された `strategy`
    - cuz: `met` では常に `-`
  - unMet: 挿入されなかったテキスト (`<<< path/to/source?` のように `?` をつけたテキストのうち、実体がなかったテキスト)
    - pathToSource: テキストへのパス
    - usedStrategy: `unMet` では常に `-`
    - cuz: なぜ対象にならなかったか
- ignored: 挿入対象にならなかったか `branch`
  - `recognized` と同様

## ライブラリ

- [cobra](https://github.com/spf13/cobra)
- [hand](https://github.com/t-star08/hand)
