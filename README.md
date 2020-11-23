# wsnc
wsnc is simple cli-tool made by go language for websocket. you can use wsnc as both client or server. So to speak, wsnc is the WebSocket version of netcat

![Gif](https://raw.githubusercontent.com/wiki/tomatod/wsnc/demo.gif)

# Manual for languages
- [日本語のマニュアル](./README_JP.md)

# Let's start
If you installed go environment.
``` sh
$ go get github.com/tomatod/wsnc
$ wsnc -h
...
```

If you want to get execution file, please download that from the following page.   
https://github.com/tomatod/wsnc/releases/tag/0.8.0

# Examples
```sh
# show help
wsnc -h

# connect to test server of websocket.org .
wsnc wss://echo.websocket.org

# send only one message and close connection. (one-shot mode)
wsnc -o -m hoge wss://echo.websocket.org

# open websocket server. (port: 8080 and path: "/")
wsnc -s 

# open websocket server. May need Sudo. (port: 80 and path: "/bar/")
wsnc -s -p 80 -u /bar/ 

# output time (-t) and loglevel (-l) in server usage.
wsnc -s -t -l
```

# Options
```
--server, -s
Run websocket server mode. Default is client mode (= false).

--path string, -u string
Specify standby path of server. This option is ignored in client mode. (default: "/")

--port integer, -p integer
Specify server listen port. In client mode, this option is ignored and argument (ws://... or wss://...) is refered to. (default: 8080)

--oneshot, -o
Run client with one shot mode (send only one message and close connection). This option is ignored in server mode. (default: false)

--broadcast, -b
If you broadcast from server, specify this option. (default: false)

--message string, -m string  
Specify message sent by server in text or binary message or by one-shot mode client.

--logtime, -t
enable timestamp of server logs. (default: false)

--loglevel, -l
enable log level of server logs. (default: false)

--debug, -d
Active additional debug logs. Now, not much different from when disabled (default: false)

--help, -h
Show help (default: false)
```

# Similar tools
- vi/websocat: https://github.com/vi/websocat
