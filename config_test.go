package main

import (
	"testing"
	"time"
)

func TestClientModeConfig(t *testing.T) {
	url := "ws://test.internal/"
	msg := "hello wsnc"
	// Ignore whether client is running collectly
	go app.Run([]string{"wsnc", "-o", "-m", msg, url})
	time.Sleep(time.Millisecond * 500)
	if appConfig.Url != url {
		t.Errorf("appConfig.Url is expected: %s  real: %s", url, appConfig.Url)
	}
	if appConfig.Message != msg {
		t.Errorf("appConfig.Message is expected: %s  real: %s", msg, appConfig.Message)
	}
	if !appConfig.IsOneShot {
		t.Errorf("appConfig.IsOneShot should be enable.")
	}
	if dialer == nil {
		t.Errorf("Client doesn't run.")
	}
}

func TestServerModeConfig(t *testing.T) {
	path := "/test/"
	port := 8081
	portStr := "8081"
	// Ignore whether server is running collectly
	go app.Run([]string{"wsnc", "-s", "-t", "-l", "-u", "/test/", "-d", "-p", portStr, path})
	time.Sleep(time.Millisecond * 1000)
	if upgrader == nil {
		t.Errorf("Server doesn't start up.")
	}
	if !appConfig.IsTime {
		t.Errorf("appConfig.IsTime should be enabled.")
	}
	if !appConfig.IsLevel {
		t.Errorf("appConfig.IsLevel should be enabled.")
	}
	if !appConfig.IsDebug {
		t.Errorf("appConfig.IsDebug should be enabled.")
	}
	if appConfig.Path != path {
		t.Errorf("appConfig.Path is expected: %s  real: %s", path, appConfig.Path)
	}
	if !(appConfig.Port == port && appConfig.PortStr == portStr) {
		t.Errorf("appConfig.Port and appConfig.PortStr is expected: %s  real: %d", portStr, appConfig.Port)
	}
}
