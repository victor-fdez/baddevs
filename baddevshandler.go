package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v5"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

//BADDevs config
type BADDevsConfig struct {
	redisHost     string
	redisPort     string
	blackListDir  string
	badDevsDomain string
	badDevsIP     string
	goDnsHash     string
	verbose       bool
}

type BADDevsCategory struct {
	Name       string `json:"name"`
	Set        uint   `json:"set"`
	NumDomains int    `json:"num_domains"`
}

type APIError struct {
	ErrorMsg    string `json:"error"`
	Description string `json:"description"`
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
var badDevsCategoriesMap map[string]*BADDevsCategory
var badDevsCategories []*BADDevsCategory

//Domain related types
type Domain struct {
	name string
	ip   string
}

func init() {
	//intialize api calls, add your api calls here
	apiCalls = []APIEndPoint{
		APIEndPoint{
			path:    "/domain-categories/",
			method:  "GET",
			name:    "DomainCategories",
			handler: badDevsDomainCategories,
		},
		APIEndPoint{
			path:    "/domain-categories/{name:[a-z_-]+}",
			method:  "GET",
			name:    "DomainCategory",
			handler: badDevsDomainCategory,
		},
		APIEndPoint{
			path:    "/domain-categories/{name:[a-z_-]+}/set",
			method:  "PUT",
			name:    "DomainCategorySet",
			handler: badDevsDomainCategorySet,
		},
		APIEndPoint{
			path:    "/domain-categories/{name:[a-z_-]+}/unset",
			method:  "PUT",
			name:    "DomainCategoryUnset",
			handler: badDevsDomainCategoryUnset,
		},
	}
	//setup api calls map for fast lookup
	apiMap = make(map[string]http.HandlerFunc)
	for _, call := range apiCalls {
		apiMap[call.path] = call.handler
	}
}

func badDevsJsonError(w http.ResponseWriter, status int, e string, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	jsError := &APIError{
		ErrorMsg:    e,
		Description: s,
	}
	js, _ := json.Marshal(jsError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func badDevsIndex(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "client/dist/index.html")
}

func badDevsDomainCategories(w http.ResponseWriter, req *http.Request) {
	js, err := json.Marshal(badDevsCategories)
	if err != nil {
		badDevsError("API /categories/ could not generate json\n")
		badDevsJsonError(w, http.StatusBadRequest, "no-categories", "Could not get categories from server")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func badDevsDomainCategory(w http.ResponseWriter, req *http.Request) {
	category, ok := badDevsUtilGetCategory(w, req)
	if !ok {
		return
	}
	badDevsUtilReturnCategory(category, w, req)
}

func badDevsDomainCategorySet(w http.ResponseWriter, req *http.Request) {
	category, ok := badDevsUtilGetCategory(w, req)
	if !ok {
		return
	}
	ok = badDevsRedisAddDomains(category, w, req)
	if !ok {
		return
	}
	category.Set = 1
	badDevsUtilReturnCategory(category, w, req)
}

func badDevsDomainCategoryUnset(w http.ResponseWriter, req *http.Request) {
	category, ok := badDevsUtilGetCategory(w, req)
	if !ok {
		return
	}
	ok = badDevsRedisRemoveDomains(category, w, req)
	if !ok {
		return
	}
	category.Set = 0
	badDevsUtilReturnCategory(category, w, req)
}

func badDevsRedisAddDomains(c *BADDevsCategory, w http.ResponseWriter, req *http.Request) bool {
	domains, ok := badDevsGetDomains(c, w, req)
	if !ok {
		return false
	}
	_, err := client.HMSet(badDevsConfig.goDnsHash, *domains).Result()
	if err != nil {
		badDevsError("Redis: Failed to HMSet map\n%v\nmap with %v keys", err, len(*domains))
		//for key, value := range *domains {
		//fmt.Printf("%v -> %v", key, value)
		//}
		badDevsJsonError(w, http.StatusExpectationFailed, "server-error", "Couldn't set/unset the given category")
		return false
	}
	return true
}

func badDevsRedisRemoveDomains(c *BADDevsCategory, w http.ResponseWriter, req *http.Request) bool {
	domainsM, ok := badDevsGetDomains(c, w, req)
	if !ok {
		return false
	}
	domains := make([]string, 0, len(*domainsM))
	for k := range *domainsM {
		domains = append(domains, k)
	}
	_, err := client.HDel(badDevsConfig.goDnsHash, domains...).Result()
	if err != nil {
		badDevsError("Redis: Failed to HDel map\n")
		badDevsJsonError(w, http.StatusExpectationFailed, "server-error", "Couldn't set/unset the given category")
		return false
	}
	return true
}

func badDevsGetDomains(c *BADDevsCategory, w http.ResponseWriter, req *http.Request) (*map[string]string, bool) {
	filePath := path.Join(badDevsConfig.blackListDir, c.Name, "domains")
	file, err := os.Open(filePath)
	if err != nil {
		badDevsError(" %v unable to read domain file\n", filePath)
		badDevsJsonError(w, http.StatusExpectationFailed, "server-error", "Couldn't set/unset the given category")
		return nil, false
	}
	domains := make(map[string]string, c.NumDomains)
	//loop thru each domain in the file, checking and adding it to the map
	domainRegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
	reader := bufio.NewReader(file)
	line, prefix, err := reader.ReadLine()
	for line != nil && !prefix && err == nil {
		trimmedLine := strings.TrimSpace(string(line))
		//fmt.Printf(trimmedLine)
		if domainRegExp.MatchString(trimmedLine) {
			domains[trimmedLine] = badDevsConfig.badDevsIP
		} else {
			badDevsError("Domain name is incorrectly formatted %v\n", trimmedLine)
		}
		line, prefix, err = reader.ReadLine()
	}
	if prefix {
		if prefix {
			badDevsError("Line read is to long!\n")
		} else {
			badDevsError("%v\n", err)
		}
		badDevsJsonError(w, http.StatusExpectationFailed, "server-error", "Couldn't set/unset the given category")
		return nil, false
	}
	return &domains, true
}

func badDevsUtilGetCategory(w http.ResponseWriter, req *http.Request) (*BADDevsCategory, bool) {
	vars := mux.Vars(req)
	name := vars["name"]
	category, ok := badDevsCategoriesMap[name]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		badDevsError("API GET /domain-categories/{name} could not find that category name\n")
		badDevsJsonError(w, http.StatusBadRequest, "no-categories", "Could not find that Category")
		return nil, false
	}
	return category, true
}

func badDevsUtilReturnCategory(c *BADDevsCategory, w http.ResponseWriter, req *http.Request) {
	js, err := json.Marshal(c)
	if err != nil {
		badDevsJsonError(w, http.StatusExpectationFailed, "server-error", "Internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func badDevsInitCategories() bool {
	// read list of directories in black list dir, each
	// corresponds to a category of urls
	categoryDirs, err := ioutil.ReadDir(badDevsConfig.blackListDir)
	if err != nil {
		badDevsError("Unable to open %v", badDevsConfig.blackListDir)
		return false
	}
	// create category map
	badDevsCategories = make([]*BADDevsCategory, 0, len(categoryDirs))
	badDevsCategoriesMap = make(map[string]*BADDevsCategory)
	for _, categoryDir := range categoryDirs {
		name := categoryDir.Name()
		if matched, _ := regexp.MatchString("^CATEGORIES|jstor", name); !matched {
			// open files
			filePath := path.Join(badDevsConfig.blackListDir, name, "domains")
			file, err := os.Open(filePath)
			if err != nil {
				badDevsError(" %v unable to read domain file\n", name)
				continue
			}
			// get num domains
			numDomains, err := lineCounter(file)
			if err != nil {
				numDomains = 0
			}
			// init category
			category := &BADDevsCategory{
				Name:       name,
				Set:        0,
				NumDomains: numDomains,
			}
			// store the category
			badDevsCategories = append(badDevsCategories, category)
			badDevsCategoriesMap[name] = category
			badDevsDebug("\t%v has %v domains\n", name, numDomains)
		}
	}
	return true
}

func lineCounter(r io.Reader) (int, error) {
	count := 0
	buf := make([]byte, 32*1024)
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
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
	// initialize categories
	if !badDevsInitCategories() {
		badDevsError("Unable to initialize Categories")
		panic(nil)
	}
}
