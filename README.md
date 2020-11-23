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

# Client mode commands
```
# Please run help command
>> help
COMMANDS:
   echo, e  Send message to server. Message type depend on type command parameter (default: text)
   ping, p  Send ping message to server.
   quit, q  Send close message (code: 1000) to server and finish wsnc.
   type, t  Change echo message type (text|binary|ping|close).
   help, h  Display command help.
   
# example
## send text message "hoge"
>> e "hoge"

## send binary message "bar"
>> t binary
>> e bar

## send close message and quit
>> q

## send ping message "ping"
>> p "ping"

## send specify close code (example is 1011) and quit
>> t close
>> e 1011
```
Note: Close code is the following   
https://tools.ietf.org/html/rfc6455#section-7.4.1


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
Enable timestamp of server logs. (default: false)

--loglevel, -l
Enable log level of server logs. (default: false)

--header string, -H string   
In client mode, specify any HTTP headers. If you specify multiple headers, please use multiple "-H". Example: -H "hoo:var" -H "bon:bar".

--debug, -d
Active additional debug logs. Now, not much different from when disabled (default: false)

--help, -h
Show help (default: false)
```

# Similar tools
- vi/websocat: https://github.com/vi/websocat
