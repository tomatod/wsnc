package main

import (
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var testDialer = &websocket.Dialer{}

func httpToWs(http string) string {
	return strings.Replace(http, "http", "ws", -1)
}

func getConnAndCtrl(wsUrl string) (*websocket.Conn, *ClientController, error) {
	conn, _, err := testDialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, nil, err
	}

	ctrl := &ClientController{}
	ctrl.IsConnect = true
	ctrl.IsAwaitMsg = true
	ctrl.setHandlers(conn, nil)
	go readMessageFromServer(conn, ctrl)

	return conn, ctrl, nil
}

func TestServerSimpleStartUp(t *testing.T) {
	appConfig = Config{}
	appConfig.IsServer = true
	appConfig.Path = "/"

	// server open
	testServer := httptest.NewServer(server())
	defer testServer.Close()

	// client connect
	wsUrl := httpToWs(testServer.URL) + "/"
	conn, ctrl, err := getConnAndCtrl(wsUrl)
	defer closeClientConn(conn)

	// send message
	sendMsg := "test"
	err = conn.WriteMessage(websocket.TextMessage, []byte(sendMsg))
	if err != nil {
		t.Errorf(err.Error())
	}

	for !ctrl.checkMsgFromServer() {
		time.Sleep(time.Millisecond * 50)
	}

	if ctrl.Message != sendMsg {
		t.Errorf("Message => send: %s  reply: %s", sendMsg, ctrl.Message)
	}
	if ctrl.MsgType != websocket.TextMessage {
		t.Errorf("Message type => send: %d  reply: %d", websocket.TextMessage, ctrl.MsgType)
	}

}

func TestServerStaticMsgMode(t *testing.T) {
	// server setting
	appConfig = Config{}
	appConfig.IsServer = true
	appConfig.Path = "/test/"
	appConfig.Message = "static"

	testServer := httptest.NewServer(server())
	defer testServer.Close()

	wsUrl := httpToWs(testServer.URL) + "/test/"
	conn, ctrl, err := getConnAndCtrl(wsUrl)
	defer closeClientConn(conn)

	sendMsg := "test"
	err = conn.WriteMessage(websocket.TextMessage, []byte(sendMsg))
	if err != nil {
		t.Errorf(err.Error())
	}

	for !ctrl.checkMsgFromServer() {
		time.Sleep(time.Millisecond * 50)
	}

	if ctrl.Message != "static" {
		t.Errorf("Message => send: %s  reply: %s", sendMsg, ctrl.Message)
	}
	if ctrl.MsgType != websocket.TextMessage {
		t.Errorf("Message type => send: %d  reply: %d", websocket.TextMessage, ctrl.MsgType)
	}

}
