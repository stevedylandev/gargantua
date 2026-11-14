package main

import (
	"fmt"
	"log"
	"net/http"
)

// Plan
// - [x] Serve robots.txt with noindex
// - [x] Serve root path html with babbler and links to new paths
// - [x] Serve dynamic path that takes in path param to serve babbler and more links
// - [ ] Slow down responses

func main() {
	if portEnv == "" {
		portEnv = "8090"
	}

	http.HandleFunc("/robots.txt", RobotsHandler)
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/{slug}", IndexHandler)

	port := fmt.Sprintf(":%s", portEnv)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, RequestMiddleware(http.DefaultServeMux)))
}
