package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func BadIndex(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/bad.html")
}

func badHandler(r *mux.Router) {
	//setup reachable urls
	r.Host("")
	r.HandleFunc("/", BadIndex)
}
