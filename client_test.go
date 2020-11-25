package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var wsMsgTest struct {
	mtype int
	msg   string
	code  int
	err   error
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		wsMsgTest.err = err
		return
	}
	defer conn.Close()

	conn.SetPingHandler(func(msg string) error {
		if wsMsgTest.mtype != 0 {
			return nil
		}
		wsMsgTest.mtype = websocket.PingMessage
		wsMsgTest.msg = msg
		return nil
	})
	conn.SetCloseHandler(func(code int, msg string) error {
		if wsMsgTest.mtype != 0 {
			return nil
		}
		wsMsgTest.mtype = websocket.CloseMessage
		wsMsgTest.msg = msg
		wsMsgTest.code = code
		return nil
	})

	for {
		mt, msg, err := conn.ReadMessage()
		if wsMsgTest.mtype != 0 {
			return
		}
		wsMsgTest.mtype = mt
		wsMsgTest.msg = string(msg)
		wsMsgTest.err = err
		conn.Close()
		return
	}
}

func clietnTestInit() {
	wsMsgTest.mtype = 0
	wsMsgTest.msg = ""
	wsMsgTest.code = 0
	wsMsgTest.err = nil
	appConfig = Config{}
	appConfig.Path = "/"
	clientNowMessageType = websocket.TextMessage
}

func TestTextMsg(t *testing.T) {
	clietnTestInit()

	// start up test server
	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()
	appConfig.Url = httpToWs(testServer.URL) + appConfig.Path

	// input to STDIN
	r1, w1, _ := os.Pipe()
	stdin := os.Stdin
	os.Stdin = r1
	defer func() {
		os.Stdin = stdin
	}()
	if _, err := w1.Write([]byte("e test\nquit")); err != nil {
		t.Errorf(err.Error())
	}
	w1.Close()

	// start client
	if err := startClient(); err != nil {
		t.Errorf(err.Error())
	}

	// awiat processing
	time.Sleep(time.Millisecond * 500)

	// check message type and message reached the server
	if !(wsMsgTest.mtype == clientNowMessageType && wsMsgTest.msg == "test") {
		if wsMsgTest.err != nil {
			t.Errorf(wsMsgTest.err.Error())
		}
		t.Error(wsMsgTest)
	}
}

func TestBinaryMsg(t *testing.T) {
	clietnTestInit()

	// start up test server
	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()
	appConfig.Url = httpToWs(testServer.URL) + appConfig.Path

	// input to STDIN
	r1, w1, _ := os.Pipe()
	stdin := os.Stdin
	os.Stdin = r1
	defer func() {
		os.Stdin = stdin
	}()
	if _, err := w1.Write([]byte("type binary\ne hello\nquit")); err != nil {
		t.Errorf(err.Error())
	}
	w1.Close()

	// start client
	if err := startClient(); err != nil {
		t.Errorf(err.Error())
	}

	// awiat processing
	time.Sleep(time.Millisecond * 500)

	// check message type and message reached the server
	if !(wsMsgTest.mtype == websocket.BinaryMessage && wsMsgTest.msg == "hello") {
		if wsMsgTest.err != nil {
			t.Errorf(wsMsgTest.err.Error())
		}
		t.Error(wsMsgTest)
	}
}

func TestPingMsg(t *testing.T) {
	clietnTestInit()

	// start up test server
	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()
	appConfig.Url = httpToWs(testServer.URL) + appConfig.Path

	// input to STDIN
	r1, w1, _ := os.Pipe()
	stdin := os.Stdin
	os.Stdin = r1
	defer func() {
		os.Stdin = stdin
	}()
	if _, err := w1.Write([]byte("p ping\nquit")); err != nil {
		t.Errorf(err.Error())
	}
	w1.Close()

	// start client
	if err := startClient(); err != nil {
		t.Errorf(err.Error())
	}

	// awiat processing
	time.Sleep(time.Millisecond * 500)

	// check message type and message reached the server
	if !(wsMsgTest.mtype == websocket.PingMessage && wsMsgTest.msg == "ping") {
		if wsMsgTest.err != nil {
			t.Errorf(wsMsgTest.err.Error())
		}
		t.Error(wsMsgTest)
	}
}

func TestCloseMsg(t *testing.T) {
	clietnTestInit()

	// start up test server
	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()
	appConfig.Url = httpToWs(testServer.URL) + appConfig.Path

	// input to STDIN
	r1, w1, _ := os.Pipe()
	stdin := os.Stdin
	os.Stdin = r1
	defer func() {
		os.Stdin = stdin
	}()
	if _, err := w1.Write([]byte("q")); err != nil {
		t.Errorf(err.Error())
	}
	w1.Close()

	// start client
	if err := startClient(); err != nil {
		t.Errorf(err.Error())
	}

	// awiat processing
	time.Sleep(time.Millisecond * 500)

	// check message type and message reached the server
	if !(wsMsgTest.mtype == websocket.CloseMessage && wsMsgTest.code == 1000) {
		if wsMsgTest.err != nil {
			t.Errorf(wsMsgTest.err.Error())
		}
		t.Error(wsMsgTest)
	}
}

func TestSpecifiedCloseMsg(t *testing.T) {
	clietnTestInit()

	// start up test server
	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()
	appConfig.Url = httpToWs(testServer.URL) + appConfig.Path

	// input to STDIN
	r1, w1, _ := os.Pipe()
	stdin := os.Stdin
	os.Stdin = r1
	defer func() {
		os.Stdin = stdin
	}()
	if _, err := w1.Write([]byte("t close\ne 1011")); err != nil {
		t.Errorf(err.Error())
	}
	w1.Close()

	// start client
	if err := startClient(); err != nil {
		t.Errorf(err.Error())
	}

	// awiat processing
	time.Sleep(time.Millisecond * 500)

	// check message type and message reached the server
	if !(wsMsgTest.mtype == websocket.CloseMessage && wsMsgTest.code == 1011) {
		if wsMsgTest.err != nil {
			t.Errorf(wsMsgTest.err.Error())
		}
		t.Error(wsMsgTest)
	}
}
