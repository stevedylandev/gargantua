package main 

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var portEnv = os.Getenv("PORT")

//go:embed templates/*.html
var templatesFS embed.FS
var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.html"))

//go:embed robots.txt
var robotsTxt string

type PageData struct {
	Title   string
	Content string
	Links   []string
}

func RobotsHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(writer, robotsTxt)
}

func IndexHandler(writer http.ResponseWriter, req *http.Request) {
	babbler := NewBabbler()
	var title string
	if req.PathValue("slug") != "" {
		title = req.PathValue("slug")
	} else {
		title = babbler.Generate(2)
	}
	content := babbler.Generate(500)
	links := babbler.Links(5)
	data := PageData{
		Title:   title,
		Content: content,
		Links:   links,
	}
	if err := tmpl.ExecuteTemplate(writer, "index.html", data); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
	}
}

func RequestMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %s\n", req.Method, req.URL, req.UserAgent())
		writer.Header().Set("X-Robots-Tag", "noindex, nofollow")
		handler.ServeHTTP(writer, req)
	})
}

