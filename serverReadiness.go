package main

import (
	"net/http"
)

func statusHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Status-URI", "200")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
