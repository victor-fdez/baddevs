package main

import (
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func some(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "index.html")
}

func main() {
	var dir, port, host string
	var config BADDevsConfig

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&port, "port", "8000", "the port to bind to.")
	flag.StringVar(&host, "host", "0.0.0.0", "the ip address to bind to.")
	flag.StringVar(&config.badDevsDomain, "baddevs-domain", "baddevs.io", "the domain used by server to control dns.")
	flag.StringVar(&config.redisHost, "redis-host", "0.0.0.0", "the ip of the Redis server.")
	flag.StringVar(&config.redisPort, "redis-port", "6379", "the port of the Redis server.")
	flag.StringVar(&config.blackListDir, "black-list-dir", "blacklists/", "The directory containing the black list category folders.")
	flag.Parse()

	r := mux.NewRouter()

	// Setup logs
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Handler: loggedRouter,
		Addr:    host + ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Setup handler
	badDevsHandler(r, config)
	badHandler(r)

	// This will serve files under http://<host>:<port>/static/<filename>
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./client/dist/assets/"))))
	r.PathPrefix("/{filename}.js").Handler(http.FileServer(http.Dir("./client/dist/")))
	r.PathPrefix("/{filename}.map").Handler(http.FileServer(http.Dir("./client/dist/")))

	badDevsInfo("BADdevs starting @ %v:%v\n", host, port)

	log.Fatal(srv.ListenAndServe())
}
