
# wsnc
wsnc は Go 言語で開発されたシンプルな WebSocket 用の CLI ツールです。wsnc 単体でクライアントとサーバーの両方に使うことができます。

# 開始方法
``` sh
$ go get github.com/tomatod/wsnc
$ wsnc -h
...
```

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

--debug, -d
追加のデバッグログを出力したい場合には指定してください。現状、あまり変わりません・・・。(デフォルト: false)

--help, -h
ヘルプを表示します。(default: false)
```

# 類似ツール
- vi/websocat: https://github.com/vi/websocat
