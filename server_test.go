package main

import (
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var testDialer = &websocket.Dialer{}

func httpToWs(http string) string {
	return strings.Replace(http, "http", "ws", -1)
}

func getConn(wsUrl string) (*websocket.Conn, error) {
	conn, _, err := testDialer.Dial(wsUrl, nil)
	if err != nil {
		return nil, err
	}

	go readMessageFromServer(conn)

	return conn, nil
}

func initStdio(stdin *os.File, stdout *os.File, stderr *os.File) {
	os.Stdin = stdin
	os.Stdout = stdout
	os.Stderr = stderr
}

func TestServerSimpleStartUp(t *testing.T) {
	// server config
	clientClose = &ClientClose{}
	appConfig = Config{}
	appConfig.IsServer = true
	appConfig.Path = "/"

	// server open
	testServer := httptest.NewServer(server())
	defer testServer.Close()

	// client connect
	wsUrl := httpToWs(testServer.URL) + appConfig.Path
	conn, err := getConn(wsUrl)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer closeClientConn(conn)

	// get output of client (stdout) and server (stderr)
	cr, cw, _ := os.Pipe()
	sr, sw, _ := os.Pipe()
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	defer initStdio(stdin, stdout, stderr)
	os.Stdout = cw
	os.Stderr = sw

	// send message from client
	sendMsg := "test"
	err = conn.WriteMessage(websocket.TextMessage, []byte(sendMsg))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// await of writting
	time.Sleep(time.Millisecond * 100)
	cw.Close()
	sw.Close()

	// check client output (this is server reply)
	cout, err := ioutil.ReadAll(cr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	ex := "< test (Text)"
	if !strings.Contains(string(cout), ex) {
		t.Errorf("expected: %s / real: %s", ex, string(cout))
	}

	// check server output (this is client message)
	sout, err := ioutil.ReadAll(sr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	ex = "> test (Text)"
	if !strings.Contains(string(sout), ex) {
		t.Errorf("expected: %s / real: %s", ex, string(sout))
	}
}

func TestServerStaticMsgMode(t *testing.T) {
	// server config
	clientClose = &ClientClose{}
	appConfig = Config{}
	appConfig.IsServer = true
	appConfig.Path = "/path/"
	appConfig.Message = "static"

	// server open
	testServer := httptest.NewServer(server())
	defer testServer.Close()

	// client connect
	wsUrl := httpToWs(testServer.URL) + appConfig.Path
	conn, err := getConn(wsUrl)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer closeClientConn(conn)

	// get output of client (stdout) and server (stderr)
	cr, cw, _ := os.Pipe()
	sr, sw, _ := os.Pipe()
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	defer initStdio(stdin, stdout, stderr)
	os.Stdout = cw
	os.Stderr = sw

	// send message from client
	sendMsg := "hello"
	err = conn.WriteMessage(websocket.BinaryMessage, []byte(sendMsg))
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// await of writting
	time.Sleep(time.Millisecond * 100)
	cw.Close()
	sw.Close()

	// check client output (this is server reply)
	cout, err := ioutil.ReadAll(cr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	ex := "< static (Binary)"
	if !strings.Contains(string(cout), ex) {
		t.Errorf("expected: %s / real: %s", ex, string(cout))
	}

	// check server output (this is client message)
	sout, err := ioutil.ReadAll(sr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	ex = "> hello (Binary)"
	if !strings.Contains(string(sout), ex) {
		t.Errorf("expected: %s / real: %s", ex, string(sout))
	}
}
