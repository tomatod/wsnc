package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	dialer      = &websocket.Dialer{}
	clientClose = &ClientClose{}
)

// To check whether client connection is closed.
type ClientClose struct {
	mu  sync.RWMutex
	Yes bool
}

func (c *ClientClose) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Yes = true
}

func (c *ClientClose) check() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Yes
}

// Entry point of client mode.
func startClient() error {
	var headers http.Header = nil
	var err error
	if appConfig.Headers != nil {
		headers, err = makeHttpHeader(appConfig.Headers)
		if err != nil {
			return err
		}
	}
	// Connect to websocket server.
	conn, _, err := dialer.Dial(appConfig.Url, headers)
	if err != nil {
		return err
	}
	defer closeClientConn(conn)

	// Set control message (Ping/Pong/Close) handler
	setClientHandlers(conn)

	// Read text or binary message loop
	// control messages are read by dedicated handlers
	go readMessageFromServer(conn)

	// if one-shot mode
	if appConfig.IsOneShot {
		return clientOneShot(conn)
	}

	return clientWebSocketLoop(conn)
}

func clientWebSocketLoop(conn *websocket.Conn) error {
	rlogf("connected to %s ...", appConfig.Url)
	prompt()

	// Read user input.
	// TODO: Supports special inputs such as arrow keys.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		// Judge which wsnc command is called.
		cmdName, arg := getWscCmdNameAndArg(input)
		cmd, exist := wscCmds[cmdName]
		if !exist {
			if input != "" {
				rlogf("\"%s\" is invalid command. Show help by \"help\" or \"h\" command.", input)
			}
			prompt()
			continue
		}

		// Judge whether loop continue depneding on result of command.
		isContinue, err := cmd(conn, arg)
		if !isContinue {
			clientClose.close()
			return err
		}
		if err != nil {
			rlogf(err.Error())
		}

		prompt()
	}
	return nil
}

// Loop for recieving text or binary messages.
func readMessageFromServer(conn *websocket.Conn) {
	for {
		mt, recvMsg, err := conn.ReadMessage()
		if clientClose.check() {
			// one-shot mode.
			if appConfig.IsOneShot && err == nil {
				fmt.Println(string(recvMsg))
			}
			return
		}
		if err != nil {
			rlogf(err.Error())
		}
		if mt == -1 {
			rlogf("\nRecieved Unexpected packet and connection may be discarded from server.")
			return
		}
		msgFromServer(mt, string(recvMsg), 0)
	}
}

// Set Handler for recieve control message (Ping/Pong/Close)
func setClientHandlers(conn *websocket.Conn) {
	conn.SetCloseHandler(func(code int, recvMsg string) error {
		clientClose.close()
		if !appConfig.IsOneShot {
			msgFromServer(websocket.CloseMessage, recvMsg, code)
		}
		return nil
	})

	conn.SetPingHandler(func(recvMsg string) error {
		msgFromServer(websocket.PingMessage, recvMsg, 0)
		return nil
	})

	conn.SetPongHandler(func(recvMsg string) error {
		msgFromServer(websocket.PongMessage, recvMsg, 0)
		return nil
	})
}

// If client is running in one-shot mode, this function is called.
func clientOneShot(conn *websocket.Conn) error {
	err := conn.WriteMessage(websocket.TextMessage, []byte(appConfig.Message))
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.CloseMessage, []byte{0x03, 0xe8})
}

// Print prompt for wsnc client.
func prompt() {
	fmt.Print(">> ")
}

// Print message from server.
func msgFromServer(mtype int, msg string, code int) {
	if code != 0 {
		fmt.Printf("\r< Close Code %d\n", code)
		return
	}
	mtypeStr := strMsgType[mtype]
	fmt.Printf("\r< %s (%s)\n", msg, mtypeStr)
}

// When websocket connection loop is end, this function is always called.
func closeClientConn(conn *websocket.Conn) {
	clientClose.close()
	// wait close message from server.
	time.Sleep(time.Millisecond * 500)
	conn.Close()
	if !appConfig.IsOneShot {
		rlogf("Closed connection.")
	}
}

// Now messagetype for echo command.
var clientNowMessageType = websocket.TextMessage

// Client command map
var wscCmds = map[string]func(*websocket.Conn, string) (bool, error){
	// full name
	"echo": wscCmdEcho,
	"ping": wscCmdPing,
	"quit": wscCmdClose,
	"type": wscCmdMessageType,
	"help": wscCmdHelp,

	// alias
	"e": wscCmdEcho,
	"p": wscCmdPing,
	"q": wscCmdClose,
	"t": wscCmdMessageType,
	"h": wscCmdHelp,
}

// Command for send message.
func wscCmdEcho(conn *websocket.Conn, arg string) (bool, error) {
	if strings.Trim(arg, " ") == "" {
		return true, errors.New("Message is empty.")
	}
	if clientNowMessageType == websocket.CloseMessage {
		clientClose.close()
		code, err := strconv.Atoi(arg)
		if err != nil {
			e := conn.WriteMessage(clientNowMessageType, []byte(arg))
			return false, e
		}
		codeUint16 := uint16(code)
		bytes := []byte{byte(codeUint16 >> 8), byte(codeUint16)}
		err = conn.WriteMessage(clientNowMessageType, bytes)
		return false, err
	}
	err := conn.WriteMessage(clientNowMessageType, []byte(arg))
	return true, err
}

// Command for sending ping of websocket.
func wscCmdPing(conn *websocket.Conn, arg string) (bool, error) {
	err := conn.WriteMessage(websocket.PingMessage, []byte(arg))
	return true, err
}

// Command for close of websocket.
func wscCmdClose(conn *websocket.Conn, arg string) (bool, error) {
	rdebugf("Client select closing connection.")
	// send close code 1000
	err := conn.WriteMessage(websocket.CloseMessage, []byte{0x03, 0xe8})
	return false, err
}

// Command for set message type of websocket.
func wscCmdMessageType(conn *websocket.Conn, arg string) (bool, error) {
	if t, exsit := messageTypes[arg]; exsit {
		clientNowMessageType = t
		return true, nil
	}
	return true, errors.New("specified message type is invalid.")
}

// Help command
func wscCmdHelp(conn *websocket.Conn, arg string) (bool, error) {
	fmt.Println(`COMMANDS:
   echo, e  Send message to server. Message type depend on type command parameter (default: text)
   ping, p  Send ping message to server.
   quit, q  Send close message (code: 1000) to server and finish wsnc.
   type, t  Change echo message type (text|binary|ping|close).
   help, h  Display command help.`)
	return true, nil
}

// Parser for wsnc client command and arg.
func getWscCmdNameAndArg(input string) (string, string) {
	inputTrimed := strings.Trim(input, " ")
	args := strings.SplitAfterN(inputTrimed, " ", 2)
	if len(args) < 2 {
		return strings.Trim(args[0], " "), ""
	}
	return strings.Trim(args[0], " "), strings.Trim(args[1], " ")
}
