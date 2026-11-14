package main

import (
	"embed"
	"fmt"
	"log"
	"html/template"
	"net/http"
	"os"
)

// Plan
// - [x] Serve robots.txt with noindex
// - [/] Serve root path html with babbler and links to new paths
// - [ ] Serve dynamic path that takes in path param to serve babbler and more links
// - [ ] Slow down responses

var portEnv = os.Getenv("PORT")

//go:embed templates/*.html
var templatesFS embed.FS
var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.html"))

//go:embed robots.txt
var robotsTxt string

type PageData struct {
	Title string
	Items []string
}

func robotsHandler(writer http.ResponseWriter, req *http.Request){
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(writer, robotsTxt)
}

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	data := PageData{
		Title: "Babbler",
		Items: []string{"First link", "second link", "third link"},
	}
	if err := tmpl.ExecuteTemplate(writer, "index.html", data); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}
}

func requestMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %s\n", req.Method, req.URL, req.RemoteAddr)
		writer.Header().Set("X-Robots-Tag", "noindex, nofollow")
		handler.ServeHTTP(writer, req)
	})
}

func main() {
	if portEnv == "" {
		portEnv = "8090"
	}

	http.HandleFunc("/robots.txt", robotsHandler)
	http.HandleFunc("/", indexHandler)

	port := fmt.Sprintf(":%s", portEnv)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, requestMiddleware(http.DefaultServeMux)))
}
