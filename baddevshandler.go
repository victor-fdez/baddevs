package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"net/http"
	"regexp"
)

//BADDevs config
type BADDevsConfig struct {
	redisHost     string
	redisPort     string
	blackListDir  string
	badDevsDomain string
}

type APIEndPoint struct {
	path    string
	method  string
	name    string
	handler http.HandlerFunc
}

var apiCalls []APIEndPoint
var apiMap map[string]http.HandlerFunc
var client *redis.Client
var badDevsConfig BADDevsConfig

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
			handler: badDevsAddDomain,
		},
		APIEndPoint{
			path:    "/domain-categories/",
			method:  "GET",
			name:    "DomainCategories",
			handler: badDevsDomainCategories,
		},
		APIEndPoint{
			path:    "/domain-categories/{name:[a-z_-]+}",
			method:  "GET",
			name:    "DomainCategories",
			handler: badDevsDomainCategory,
		},
	}
	//setup api calls map for fast lookup
	apiMap = make(map[string]http.HandlerFunc)
	for _, call := range apiCalls {
		apiMap[call.path] = call.handler
	}
}

type APIError struct {
	ErrorMsg    string `json:"error"`
	Description string `json:"description"`
}

func badDevsJsonError(w http.ResponseWriter, e string, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	jsError := &APIError{
		ErrorMsg:    e,
		Description: s,
	}
	js, _ := json.Marshal(jsError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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

func badDevsAddDomain(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/dist/index.html")
}

func badDevsDomainCategory(w http.ResponseWriter, req *http.Request) {

}

func badDevsDomainCategories(w http.ResponseWriter, req *http.Request) {
	categoryDirs, err := ioutil.ReadDir(badDevsConfig.blackListDir)
	if err != nil {
		badDevsError("API /categories/ could not read black list directory @ %v\n",
			badDevsConfig.blackListDir)
		badDevsJsonError(w, "no-categories", "Could not get categories from server")
		return
	}
	categories := make([]string, 0, len(categoryDirs))
	for _, categoryDir := range categoryDirs {
		name := categoryDir.Name()
		if matched, _ := regexp.MatchString("^CATEGORIES", name); !matched {
			categories = append(categories, categoryDir.Name())
		}
	}
	js, err := json.Marshal(categories)
	if err != nil {
		badDevsError("API /categories/ could not generate json\n")
		badDevsJsonError(w, "no-categories", "Could not get categories from server")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func badDevsAPI(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, req, "index.html")
}

func badDevsHandler(r *mux.Router, config BADDevsConfig) {
	//setup reachable urls
	badDevsInfo("Setting up routes for %v\n", config.badDevsDomain)
	r = r.Host(config.badDevsDomain).Subrouter()
	r.HandleFunc("/", badDevsIndex)
	//setup api functions
	for _, call := range apiCalls {
		r.
			Methods(call.method).
			Path(call.path).
			Name(call.name).
			Handler(call.handler)
	}
	//setup Redis client
	client = redis.NewClient(&redis.Options{
		Addr:     config.redisHost + ":" + config.redisPort,
		Password: "",
		DB:       0,
	})
	//ping Redis
	_, err := client.Ping().Result()
	if err != nil {
		badDevsError("Redis did not respond to ping", config.badDevsDomain)
		panic(err)
	}
	badDevsInfo("Redis connected %v:%v\n", config.redisHost, config.redisPort)
	//store config for later
	badDevsConfig = config
}
