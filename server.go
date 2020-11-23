package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = &websocket.Upgrader{}

// Entry point of server mode.
func startServer() error {
	outputServerInfo()
	if err := http.ListenAndServe(":"+appConfig.PortStr, server()); err != nil {
		return errLogf(err.Error())
	}
	return nil
}

func server() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(appConfig.Path, mainHandler)
	return mux
}

// HTTP request handler for upgrade to websocket
func mainHandler(w http.ResponseWriter, r *http.Request) {
	debugLogf("=== Received request header ===\n%s", getHeaderStr(r.Header))
	// Upgrade http to websocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errLogf("Upgrade error. Message is \"%s\".", err.Error())
		return
	}
	infoLogf("Connected from client(%s)...", r.RemoteAddr)
	defer closeServerConn(r, conn)
	addConnectionFromClient(r.RemoteAddr, conn)
	setServerHandlers(conn, r)

	// If upgrade is ok, websocket loop start for each client.
	for {
		// Read client message.
		mt, recvMsg, err := conn.ReadMessage()
		if !existConnectionFromClient(r.RemoteAddr) {
			return
		}
		if err != nil {
			warnLogf("Invalid message from client(%s). Error is \"%s\"", r.RemoteAddr, err.Error())
			if mt == -1 {
				deleteConnectionFromClient(r.RemoteAddr)
				warnLogf("Connection from client(%s) may be closed.", r.RemoteAddr)
				return
			}
			continue
		}

		infoLogf("Client(%s): %s", r.RemoteAddr, string(recvMsg))

		if appConfig.IsBroadcast {
			sendBroadcast(r.RemoteAddr, mt, string(recvMsg))
			continue
		}
		if appConfig.Message != "" {
			if err = conn.WriteMessage(mt, []byte(appConfig.Message)); err != nil {
				warnLogf(err.Error())
				return
			}
		}
		if err = conn.WriteMessage(mt, []byte(recvMsg)); err != nil {
			warnLogf(err.Error())
			return
		}
	}
}

// Send broadcast if "-b" is selected.
func sendBroadcast(addr string, mt int, recvMsg string) {
	connectionsFromClients.mu.RLock()
	defer connectionsFromClients.mu.RUnlock()
	for key, conn := range connectionsFromClients.m {
		err := conn.WriteMessage(mt, []byte("\""+recvMsg+"\" by Client("+addr+")"))
		if err != nil {
			warnLogf("Send to client(%s) is error: %s\n", key, err.Error())
		}
	}
}

// Output log, when server start up.
func outputServerInfo() {
	var msg string = "repeat"
	if appConfig.Message != "" {
		msg = "\"" + appConfig.Message + "\""
	}
	infoLogf("Server start up => port:%s | path:%s | message:%s", appConfig.PortStr, appConfig.Path, msg)
}

func getReplyMessageType(req *http.Request, conn *websocket.Conn, messageType int, recvMsg []byte) int {
	switch messageType {
	case websocket.TextMessage:
		infoLogf("Client(%s): %s", req.RemoteAddr, string(recvMsg))
		return messageType
	case websocket.BinaryMessage:
		infoLogf("Client(%s): %s", req.RemoteAddr, string(recvMsg))
		return messageType
	default:
		warnLogf("Client(%s): request invalid message type.", req.RemoteAddr)
		return websocket.CloseMessage
	}
}

func setServerHandlers(conn *websocket.Conn, req *http.Request) {
	conn.SetPingHandler(func(recvMsg string) error {
		infoLogf("Client(%s): request ping message \"%s\".", req.RemoteAddr, recvMsg)
		// reply same messsage from client.
		return conn.WriteMessage(websocket.PongMessage, []byte(recvMsg))
	})
	conn.SetCloseHandler(func(code int, recvMsg string) error {
		deleteConnectionFromClient(req.RemoteAddr)
		infoLogf("Client(%s): request closing connection (code: %d). Message is \"%s\"", req.RemoteAddr, code, recvMsg)
		// reply same close code from client.
		codeUint16 := uint16(code)
		replyByte := []byte{byte(codeUint16 >> 8), byte(codeUint16)}
		return conn.WriteMessage(websocket.CloseMessage, replyByte)
	})
}

// When websocket connection loop is end, this function is always called.
func closeServerConn(req *http.Request, conn *websocket.Conn) {
	conn.Close()
	deleteConnectionFromClient(req.RemoteAddr)
	infoLogf("Connection closed on client(%s)", req.RemoteAddr)
}

// Connection from client list.
type ConnectionsFromClients struct {
	m  map[string]*websocket.Conn
	mu sync.RWMutex
}

var connectionsFromClients = ConnectionsFromClients{
	m: map[string]*websocket.Conn{},
}

func addConnectionFromClient(addr string, conn *websocket.Conn) {
	connectionsFromClients.mu.Lock()
	defer connectionsFromClients.mu.Unlock()
	if _, exist := connectionsFromClients.m[addr]; !exist {
		connectionsFromClients.m[addr] = conn
	}
}

func deleteConnectionFromClient(addr string) {
	connectionsFromClients.mu.Lock()
	defer connectionsFromClients.mu.Unlock()
	if _, exist := connectionsFromClients.m[addr]; exist {
		delete(connectionsFromClients.m, addr)
	}
}

func existConnectionFromClient(addr string) bool {
	connectionsFromClients.mu.RLock()
	defer connectionsFromClients.mu.RUnlock()
	_, exist := connectionsFromClients.m[addr]
	return exist
}
