
# wsnc
wsnc は Go 言語で開発されたシンプルな WebSocket 用の CLI ツールです。wsnc 単体でクライアントとサーバーの両方に使うことができます。簡単に言うと netcat というコマンドの WebSocket 版です。

![Gif](https://raw.githubusercontent.com/wiki/tomatod/wsnc/demo.gif)

# 開始方法
Go の実行環境が構築済みの場合   
``` sh
$ go get github.com/tomatod/wsnc
$ wsnc -h
...
```
実行ファイルが欲しい場合は、以下のページからダウンロードしてください。  
https://github.com/tomatod/wsnc/releases/tag/0.8.0

# 使用例

```sh
# ヘルプを表示します。
wsnc -h

# websocket.org のテストサーバーにクライアントとして接続します。
wsnc wss://echo.websocket.org

# websocket.org のテストサーバーに hoge というメッセージを一度送信し、接続を終了します。(one-shot モードと呼んでいます)
wsnc -o -m hoge wss://echo.websocket.org

# WebSocket サーバーを起動します。(port: 8080 and path: "/")
wsnc -s

# WebSocket サーバーを起動します。sudo が必要かもしれません。(port: 80 and path: "/bar/")
wsnc -s -p 80 -u /bar/

# サーバーとして使用する際に、時間 (-t) とログレベル (-l) を表示します。
wsnc -s -t -l
```

# クライアントモードのコマンド
```
# 基本的にはヘルプコマンドを実行して確認してください。
>> help
COMMANDS:
   echo, e  Send message to server. Message type depend on type command parameter (default: text)
   ping, p  Send ping message to server.
   quit, q  Send close message (code: 1000) to server and finish wsnc.
   type, t  Change echo message type (text|binary|ping|close).
   help, h  Display command help.
   
# 例
## "hoge" という text メッセージを送信します。
>> e "hoge"

## "bar" という binary メッセージを送信します。
>> t binary
>> e bar

## close メッセージを送信ししてコマンドを終了します。
>> q

## "ping" という文字列で、ping メッセージを送信します。
>> p "ping"

## 特定の close code (example is 1011) を close メッセージで送信して、コマンドを終了します。
>> t close
>> e 1011
```
※ close code は以下参照   
https://tools.ietf.org/html/rfc6455#section-7.4.1

# オプション
```
--server, -s
サーバーモードとして起動します。デフォルトではクライアントモードです。(デフォルト: false)

--path string, -u string
サーバーが待ち受けるパスを指定します。クライアントモードでは無視されます。(デフォルト: "/")

--port integer, -p integer
サーバーが待ち受けるポートを指定します。クライアントモードではこのオプションは無視され、引数に渡される ws://... や wss://... といった文字列が参照されます。(デフォルト: 8080)

--oneshot, -o
クライアントとして one-shot モード (一つメッセージを送信して接続を終了するモード) で実行します。サーバーモードでは無視されます。(デフォルト: false)

--broadcast, -b
サーバーをブロードキャストモードで起動したい場合にはこのオプションを指定してください。クライアントモードでは無視されます。(デフォルト: false)

--message string, -m string
サーバーから応答される text または binary メッセージの文面、または one-shot モード時にクライアントが送信するメッセージの文面を指定します。

--logtime, -t
サーバーが出力するログにタイムスタンプを含めたい場合には指定してください。クライアントモードでは無視されます。(デフォルト: false)

--loglevel, -l
サーバーが出力するログにログレベルを含めたい場合には指定してください。クライアントモードでは無視されます。(デフォルト: false)

--header string, -H string   
クライアントモードで使用する際に HTTP ヘッダーを設定することができます。複数のヘッダーを設定する場合は、複数の "-H" を使用してください。例: -H "hoo:var" -H "bon:bar"

--debug, -d
追加のデバッグログを出力したい場合には指定してください。現状、あまり変わりません・・・。(デフォルト: false)

--help, -h
ヘルプを表示します。(default: false)
```

# 類似ツール
- vi/websocat: https://github.com/vi/websocat
