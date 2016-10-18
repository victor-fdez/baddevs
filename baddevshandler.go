package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

//API related types
type APIFunc func(w http.ResponseWriter, r *http.Request)

type APIEndPoint struct {
	path    string
	method  string
	name    string
	handler http.HandlerFunc
}

var apiCalls []APIEndPoint
var apiMap map[string]http.HandlerFunc

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
			method:  "GET",
			name:    "Domains",
			handler: badDevsDomains,
		},
		APIEndPoint{
			path:    "/domains/{id}",
			method:  "GET",
			name:    "Domain",
			handler: badDevsDomain,
		},
		APIEndPoint{
			path:    "/domains/{id}",
			method:  "DELETE",
			name:    "DomainDelete",
			handler: badDevsDomainDelete,
		},
		APIEndPoint{
			path:    "/domains/add",
			method:  "PUT",
			name:    "DomainAdd",
			handler: badDevsDomainCategories,
		},
	}
	//setup api calls map for fast lookup
	apiMap = make(map[string]http.HandlerFunc)
	for _, call := range apiCalls {
		apiMap[call.path] = call.handler
	}
}

func badDevsIndex(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/dist/index.html")
}

func badDevsDomain(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/dist/index.html")
}

func badDevsDomains(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/dist/index.html")
}

func badDevsDomainDelete(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/dist/index.html")
}

func badDevsDomainAdd(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/dist/index.html")
}

func badDevsDomainCategories(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/dist/index.html")
}

func badDevsAPI(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	http.ServeFile(w, req, "index.html")
}

func badDevsHandler(r *mux.Router, s string) {
	//setup reachable urls
	badDevsInfo("Setting up routes for %v\n", s)
	r = r.Host("baddevs.io").Subrouter()
	r.HandleFunc("/", badDevsIndex)
	//setup api functions
	for _, call := range apiCalls {
		r.
			Methods(call.method).
			Path(call.path).
			Name(call.name).
			Handler(call.handler)
	}
}
