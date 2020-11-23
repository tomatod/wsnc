package main

import (
	"fmt"
	"net/http"
	"strings"
)

func getHeaderStr(h http.Header) string {
	if h == nil {
		return ""
	}

	var output string
	for k, v := range h {
		output += k + ": "
		for i, vv := range v {
			output += vv
			if i != len(v)-1 {
				output += ", "
			}
		}
		output += "\n"
	}

	return output
}

func makeHttpHeader(hs []string) (http.Header, error) {
	h := http.Header{}
	for _, v := range hs {
		tmp := strings.Split(v, ":")
		if len(tmp) != 2 {
			return nil, fmt.Errorf("Specified HTTP header is invalid.")
		}
		if v := strings.Trim(tmp[0], " "); v == "" {
			return nil, fmt.Errorf("Specified HTTP header is invalid.")
		}
		if v := strings.Trim(tmp[1], " "); v == "" {
			return nil, fmt.Errorf("%s header's parameter is empty.", tmp[0])
		}
		params := strings.Split(tmp[1], ",")
		for i, param := range params {
			params[i] = strings.Trim(param, " ")
		}
		h[tmp[0]] = params
	}
	return h, nil
}
