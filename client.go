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

var dialer = &websocket.Dialer{}

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
	clientController.setHandlers(conn, nil)
	clientController.IsConnect = true

	// Read text or binary message loop
	// control messages are read by dedicated handlers
	go readMessageFromServer(conn, clientController)

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
		if !clientController.isConnectWithServer() {
			return clientController.Error
		}
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
			return err
		}
		if err != nil {
			rlogf(err.Error())
			if cmdName == "e" || cmdName == "echo" {
				prompt()
				continue
			}
		}

		// Some command don't require to await message, so "continue".
		if cmdName == "t" || cmdName == "type" || cmdName == "h" || cmdName == "help" {
			prompt()
			continue
		}

		// Wait to recieve message from server.
		for !clientController.checkMsgFromServer() {
			time.Sleep(time.Millisecond * 50)
		}

		// print message for server depend on message types
		if err != nil {
			return err
		}

		err, brk := handleRecvMsg()
		if brk {
			return err
		}
		if err != nil {
			rlogf(err.Error())
		}

		prompt()
	}
	return nil
}

func handleRecvMsg() (error, bool) {
	clientController.mu.Lock()
	defer clientController.mu.Unlock()
	clientController.IsMessageRecv = false
	switch clientController.MsgType {
	case websocket.TextMessage, websocket.BinaryMessage:
		msgFromSererf(clientController.Message)
	case websocket.PingMessage:
		msgFromSererf("Server replied ping message: %s", clientController.Message)
	case websocket.PongMessage:
		msgFromSererf("Server replied pong message: %s", clientController.Message)
	case websocket.CloseMessage:
		msgFromSererf("Server request closing connection (code:%d). Message is %s", clientController.Code, clientController.Message)
		return clientController.Error, true
	}
	return clientController.Error, false
}

// Await message from server. This is called by goroutine and continuously run until the wsnc client end.
func readMessageFromServer(conn *websocket.Conn, ctrl *ClientController) {
	for {
		mt, recvMsg, err := conn.ReadMessage()
		if !ctrl.isConnectWithServer() {
			return
		}
		if err != nil {
			rlogf(err.Error())
			if mt == -1 {
				rlogf("\nRecieved Unexpected packet and connection may be discarded from server.")
				ctrl.connectionClose()
				return
			}
		}
		if !ctrl.isAwaitMsg() {
			msgFromSererf("%s", recvMsg)
			continue
		}
		ctrl.setRecvMsgFromServer(mt, string(recvMsg), 0, err)
	}
}

// Print prompt for wsnc client.
func prompt() {
	fmt.Print(">> ")
}

func msgFromSererf(format string, param ...interface{}) (int, error) {
	return fmt.Printf("< "+format+"\n", param...)
}

// If client is running in one-shot mode, this function is called.
func clientOneShot(conn *websocket.Conn) error {
	clientController.startAwaitMessage()
	err := conn.WriteMessage(websocket.TextMessage, []byte(appConfig.Message))
	if err != nil {
		return err
	}
	for !clientController.checkMsgFromServer() {
		time.Sleep(time.Millisecond * 50)
	}
	conn.WriteMessage(websocket.CloseMessage, []byte{0x03, 0xe8})
	clientController.connectionClose()
	fmt.Println(clientController.Message)
	return clientController.Error
}

// When websocket connection loop is end, this function is always called.
func closeClientConn(conn *websocket.Conn) {
	conn.Close()
	if !appConfig.IsOneShot {
		rlogf("Closed connection.")
	}
}

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

var clientNowMessageType = websocket.TextMessage

// Command for send text or binary message of websocket.
func wscCmdEcho(conn *websocket.Conn, arg string) (bool, error) {
	if strings.Trim(arg, " ") == "" {
		return true, errors.New("Message is empty.")
	}
	if clientNowMessageType == websocket.CloseMessage {
		code, err := strconv.Atoi(arg)
		if err != nil {
			e := conn.WriteMessage(clientNowMessageType, []byte(arg))
			return false, e
		}
		codeUint16 := uint16(code)
		bytes := []byte{byte(codeUint16 >> 8), byte(codeUint16)}
		err = conn.WriteMessage(clientNowMessageType, bytes)
		clientController.connectionClose()
		return false, err
	}
	clientController.startAwaitMessage()
	err := conn.WriteMessage(clientNowMessageType, []byte(arg))
	return true, err
}

// Command for sending ping of websocket.
func wscCmdPing(conn *websocket.Conn, arg string) (bool, error) {
	err := conn.WriteMessage(websocket.PingMessage, []byte(arg))
	clientController.startAwaitMessage()
	return true, err
}

// Command for close of websocket.
func wscCmdClose(conn *websocket.Conn, arg string) (bool, error) {
	rdebugf("Client select closing connection.")
	clientController.connectionClose()
	// send close code 1000
	err := conn.WriteMessage(websocket.CloseMessage, []byte{0x03, 0xe8})
	return false, err
}

var messageTypes = map[string]int{
	"text":   websocket.TextMessage,
	"binary": websocket.BinaryMessage,
	"close":  websocket.CloseMessage,
	"ping":   websocket.PingMessage,
	"pong":   websocket.PongMessage,
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

// This is flag set for client mode and, have information on whether connection is enable
// and whether client have recieved message from server, and recieved messsage and so on.
type ClientController struct {
	mu            sync.RWMutex
	IsMessageRecv bool
	IsConnect     bool
	IsAwaitMsg    bool
	MsgType       int
	Code          int
	Message       string
	Error         error
}

var clientController = &ClientController{}

// Notice that send message from client, and await reply message from server.
func (c *ClientController) startAwaitMessage() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.IsAwaitMsg = true
}

// Check whether the client have recieved message from server.
func (c *ClientController) checkMsgFromServer() bool {
	if !c.isConnectWithServer() {
		return true
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.IsMessageRecv
}

// When reciving message, handlers call this function.
func (c *ClientController) setRecvMsgFromServer(mtype int, msg string, code int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if mtype == websocket.CloseMessage {
		c.IsConnect = false
	}
	c.IsAwaitMsg = false
	c.Code = code
	c.Message = string(msg)
	c.MsgType = mtype
	c.Error = err
	c.IsMessageRecv = true
}

// Check whether the connection with server is still enable.
func (c *ClientController) isConnectWithServer() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.IsConnect
}

func (c *ClientController) isAwaitMsg() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.IsAwaitMsg
}

func (c *ClientController) connectionClose() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.IsConnect = false
}

// hdls: 1. PingHandler 2. PongHandler
func (c *ClientController) setHandlers(conn *websocket.Conn, closeHdl func(int, string) error, hdls ...func(string) error) {
	conn.SetCloseHandler(func(code int, recvMsg string) error {
		var err error
		if closeHdl != nil {
			err = closeHdl(code, recvMsg)
		}
		if !clientController.isAwaitMsg() {
			msgFromSererf("< Server sent close code: %d", recvMsg, code)
		}
		c.setRecvMsgFromServer(websocket.CloseMessage, recvMsg, code, err)
		return err
	})

	conn.SetPingHandler(func(recvMsg string) error {
		var err error
		if len(hdls) >= 1 {
			err = hdls[0](recvMsg)
		}
		if !clientController.isAwaitMsg() {
			msgFromSererf("< Server sent ping message: %d", recvMsg)
		}
		c.setRecvMsgFromServer(websocket.PingMessage, recvMsg, 0, err)
		return err
	})

	conn.SetPongHandler(func(recvMsg string) error {
		var err error
		if len(hdls) >= 2 {
			err = hdls[1](recvMsg)
		}
		if !clientController.isAwaitMsg() {
			msgFromSererf("< Server sent pong message: %d", recvMsg)
		}
		c.setRecvMsgFromServer(websocket.PongMessage, recvMsg, 0, err)
		return err
	})
}
