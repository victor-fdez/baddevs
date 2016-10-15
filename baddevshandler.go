package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

//API related types
type APIFunc func(w http.ResponseWriter, r *http.Request)

type APIEndPoint struct {
	path    string
	apiFunc APIFunc
}

var apiCalls []APIEndPoint
var apiMap map[string]APIFunc

//Domain related types
type Domain struct {
	name string
	ip   string
}

func init() {
	//intialize api calls, add your api calls here
	apiCalls = []APIEndPoint{
		APIEndPoint{
			path:    "/domains",
			apiFunc: badDevsAPI,
		},
		APIEndPoint{
			path:    "/domains/{id}",
			apiFunc: badDevsAPI,
		},
		APIEndPoint{
			path:    "/domains/{id}/delete",
			apiFunc: badDevsAPI,
		},
		APIEndPoint{
			path:    "/domains/add",
			apiFunc: badDevsAPI,
		},
	}
	//setup api calls map for fast lookup
	apiMap = make(map[string]APIFunc)
	for _, call := range apiCalls {
		apiMap[call.path] = call.apiFunc
	}
}

func badDevsIndex(w http.ResponseWriter, req *http.Request) {
	//fmt.Printf("%+v\n", req)
	http.ServeFile(w, req, "client/index.html")
}

func badDevsAPI(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "index.html")
}

func badDevsHandler(r *mux.Router, s string) {
	//setup reachable urls
	badDevsInfo("Setting up routes for %v\n", s)
	r = r.Host("baddevs.io").Subrouter()
	r.HandleFunc("/", badDevsIndex)
	//setup api functions
	for _, call := range apiCalls {
		r.HandleFunc(call.path, badDevsAPI)
	}
}
