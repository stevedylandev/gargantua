package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var portEnv = os.Getenv("PORT")

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "hello world\n")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.Method, r.URL, r.RemoteAddr)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	if portEnv == "" {
		portEnv = "8090"
	}

	http.HandleFunc("/", indexHandler)

	port := fmt.Sprintf(":%s", portEnv)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, logRequest(http.DefaultServeMux)))
}
