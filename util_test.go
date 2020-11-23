package main

import (
	"testing"
)

func TestMakeHttpHeader(t *testing.T) {
	hs := []string{
		"Host: host",
		"Test:test",
		"Cookie: hoo=var, bon=bar",
	}

	h, e := makeHttpHeader(hs)

	if e != nil {
		t.Error(e)
		t.Error(h)
	}
}
