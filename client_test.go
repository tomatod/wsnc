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

func clietnTestInit() {
	appConfig = Config{}
	appConfig.Path = "/"
	clientNowMessageType = websocket.TextMessage
}

func initStdio(stdin *os.File, stdout *os.File, stderr *os.File) {
	os.Stdin = stdin
	os.Stdout = stdout
	os.Stderr = stderr
}

// Execute any cmd of client intaractive mode, and get ouput of client.
func getOutputAfterExecCmd(t *testing.T, cmd string) ([]byte, []byte, error, error) {
	// Pipe set
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	defer initStdio(stdin, stdout, stderr)
	r1, w1, _ := os.Pipe()
	r2, w2, _ := os.Pipe()
	r3, w3, _ := os.Pipe()
	os.Stdin = r1
	os.Stdout = w2
	os.Stderr = w3

	// Write to STDIN. This is execution of cmd in intaractive mode.
	if _, err := w1.Write([]byte(cmd)); err != nil {
		return nil, nil, err, err
	}
	w1.Close()

	// Test server start up.
	testServer := httptest.NewServer(server())
	defer testServer.Close()

	// Start connection.
	appConfig.Url = httpToWs(testServer.URL) + appConfig.Path
	t.Log("Listen: " + appConfig.Url)
	if err := startClient(); err != nil {
		return nil, nil, err, err
	}
	time.Sleep(time.Millisecond * 1500)
	w2.Close()
	w3.Close()

	// Read output.
	stdoutO, stdoutErr := ioutil.ReadAll(r2)
	stderrO, stderrErr := ioutil.ReadAll(r3)
	return stdoutO, stderrO, stdoutErr, stderrErr
}

func TestQuitCmdForClient(t *testing.T) {
	clietnTestInit()

	_, output, _, err := getOutputAfterExecCmd(t, "quit")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if output == nil {
		t.Errorf("No reply.")
		return
	}

	expectStr := "Closed connection"
	if !strings.Contains(string(output), expectStr) {
		t.Errorf("\"%s\" is not found.", expectStr)
	}
}

func TestTextMsg(t *testing.T) {
	// text message
	appConfig = Config{}
	appConfig.Path = "/"

	output, _, err, _ := getOutputAfterExecCmd(t, "echo text\nquit")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if output == nil {
		t.Errorf("No reply.")
		return
	}

	expectStr := "< text"
	if !strings.Contains(string(output), expectStr) {
		t.Errorf("\"%s\" is not found. output is \"%s\"", expectStr, string(output))
	}

}

func TestBinaryMsg(t *testing.T) {
	clietnTestInit()
	clientNowMessageType = websocket.BinaryMessage

	output, _, err, _ := getOutputAfterExecCmd(t, "echo binary\nquit")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if output == nil {
		t.Errorf("No reply.")
		return
	}

	expectStr := "< binary"
	if !strings.Contains(string(output), expectStr) {
		t.Errorf("\"%s\" is not found. output is \"%s\"", expectStr, string(output))
	}
}

func TestPingMsg(t *testing.T) {
	clietnTestInit()

	output, _, err, _ := getOutputAfterExecCmd(t, "ping test\nquit")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if output == nil {
		t.Errorf("No reply.")
		return
	}

	expectStr := "< Server replied pong message: test"
	if !strings.Contains(string(output), expectStr) {
		t.Errorf("\"%s\" is not found. output is \"%s\"", expectStr, string(output))
	}
}

func TestOneShotMode(t *testing.T) {
	clietnTestInit()
	appConfig.Message = "oneshot"
	appConfig.IsOneShot = true

	testServer := httptest.NewServer(server())
	defer testServer.Close()
	appConfig.Url = httpToWs(testServer.URL) + appConfig.Path

	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()
	os.Stdout = w

	if err := startClient(); err != nil {
		t.Errorf(err.Error())
		return
	}

	w.Close()

	output, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if string(output) == "oneshot" {
		t.Errorf("expected: oneshot  real:%s", string(output))
	}

}
