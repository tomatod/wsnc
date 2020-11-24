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
		errLogf("Upgrade error : \"%s\".", err.Error())
		return
	}
	infoLogf("Connected from %s ...", r.RemoteAddr)
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
			warnLogf("Invalid message from %s : \"%s\"", r.RemoteAddr, err.Error())
			if mt == -1 {
				deleteConnectionFromClient(r.RemoteAddr)
				warnLogf("Connection may be closed from %s.", r.RemoteAddr)
				return
			}
			continue
		}

		strMsgT, _ := strMsgType[mt]
		infoLogf("%s > %s (%s)", r.RemoteAddr, string(recvMsg), strMsgT)

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
		err := conn.WriteMessage(mt, []byte(recvMsg))
		if err != nil {
			warnLogf("Send error to %s: \"%s\"\n", key, err.Error())
		}
	}
}

// Output log, when server start up.
func outputServerInfo() {
	infoLogf("Server start up :%s%s <<<", appConfig.PortStr, appConfig.Path)
}

func setServerHandlers(conn *websocket.Conn, req *http.Request) {
	conn.SetPingHandler(func(recvMsg string) error {
		t, _ := strMsgType[websocket.PingMessage]
		infoLogf("%s > %s (%s)", req.RemoteAddr, recvMsg, t)
		// reply same messsage from client.
		return conn.WriteMessage(websocket.PongMessage, []byte(recvMsg))
	})
	conn.SetCloseHandler(func(code int, recvMsg string) error {
		deleteConnectionFromClient(req.RemoteAddr)
		infoLogf("%s > Close Code %d", req.RemoteAddr, code)
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
	infoLogf("Connection closed from %s", req.RemoteAddr)
}

// Connection from client list and this is required to be thread safe.
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
