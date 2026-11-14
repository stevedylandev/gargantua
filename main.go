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
// - [ ] Serve robots.txt with noindex
// - [/] Serve root path html with babbler and links to new paths
// - [ ] Serve dynamic path that takes in path param to serve babbler and more links
// - [ ] Slow down responses

var portEnv = os.Getenv("PORT")

//go:embed templates/*.html
var templatesFS embed.FS

type PageData struct {
	Title string
	Items []string
}

var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.html"))

func indexHandler(writer http.ResponseWriter, req *http.Request) {
	data := PageData{
		Title: "Babbler",
		Items: []string{"First link", "second link", "third link"},
	}
	if err := tmpl.ExecuteTemplate(writer, "index.html", data); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}
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
