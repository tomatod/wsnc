package main

import (
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"strings"
	"testing"
)

var testDialer = &websocket.Dialer{}

func httpToWs(http string) string {
	return strings.Replace(http, "http", "ws", -1)
}

func getConnAndCtrl(wsUrl string) (*websocket.Conn, error) {
	conn, _, err := testDialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, err
	}

	go readMessageFromServer(conn)

	return conn, nil
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
	conn, err := getConnAndCtrl(wsUrl)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer closeClientConn(conn)

	// send message
	sendMsg := "test"
	err = conn.WriteMessage(websocket.TextMessage, []byte(sendMsg))
	if err != nil {
		t.Errorf(err.Error())
	}

	// TODO: get output and check.
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
	conn, err := getConnAndCtrl(wsUrl)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer closeClientConn(conn)

	sendMsg := "hello"
	err = conn.WriteMessage(websocket.TextMessage, []byte(sendMsg))
	if err != nil {
		t.Errorf(err.Error())
	}

	// TODO: get output and check.
}
