package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")

	w.WriteHeader(200)
	_, err := w.Write(IndexHTML)
	if err != nil {
		panic(err)
	}
}
