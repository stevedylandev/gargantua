package main 

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var portEnv = envString("PORT", "8090")

var (
	dripDelayMinMs = envInt("DRIP_DELAY_MIN_MS", 50)
	dripDelayMaxMs = envInt("DRIP_DELAY_MAX_MS", 150)
	dripChunkMax   = envInt("DRIP_CHUNK_MAX_BYTES", 60)
)

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

var headEndTag = []byte("</head>")

func splitPage(page []byte) (head, body []byte) {
	if i := bytes.Index(page, headEndTag); i >= 0 {
		return page[:i+len(headEndTag)], page[i+len(headEndTag):]
	}
	return nil, page
}

func TarpitWrite(writer http.ResponseWriter, page []byte) {
	head, body := splitPage(page)
	if len(head) > 0 {
		if _, err := writer.Write(head); err != nil {
			return
		}
		if flusher, ok := writer.(http.Flusher); ok {
			flusher.Flush()
		}
	}
	drip(writer, body)
}

func drip(writer http.ResponseWriter, data []byte) {
	flusher, canFlush := writer.(http.Flusher)
	for len(data) > 0 {
		delayMs := dripDelayMinMs
		if spread := dripDelayMaxMs - dripDelayMinMs; spread > 0 {
			delayMs += rand.Intn(spread)
		}
		time.Sleep(time.Duration(delayMs) * time.Millisecond)

		chunk := 1 + rand.Intn(dripChunkMax)
		if chunk > len(data) {
			chunk = len(data)
		}
		if _, err := writer.Write(data[:chunk]); err != nil {
			return 
		}
		if canFlush {
			flusher.Flush()
		}
		data = data[chunk:]
	}
}

func IndexHandler(writer http.ResponseWriter, req *http.Request) {
	babbler := sharedBabbler()
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
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "index.html", data); err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	TarpitWrite(writer, buf.Bytes())
}

func RequestMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %s\n", req.Method, req.URL, req.UserAgent())
		writer.Header().Set("X-Robots-Tag", "noindex, nofollow")
		handler.ServeHTTP(writer, req)
	})
}

