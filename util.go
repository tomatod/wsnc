package main

import (
	"net/http"
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
