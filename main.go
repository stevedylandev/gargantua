package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	dripDelayMinMs = envInt("GARGANTUA_DELAY_MIN_MS", 150)
	dripDelayMaxMs = envInt("GARGANTUA_DELAY_MAX_MS", 250)
	dripChunkMax   = envInt("GARGANTUA_CHUNK_MAX_BYTES", 45)
	portEnv        = envString("GARGANTUA_PORT", "8090")
	hostEnv        = envString("GARGANTUA_HOST", "127.0.0.1")
	sourceText     = envString("GARGANTUA_CORPUS", "moby.txt")
)

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/{slug}", IndexHandler)
	http.HandleFunc("/robots.txt", RobotsHandler)

	port := fmt.Sprintf("%s:%s", hostEnv, portEnv)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, RequestMiddleware(http.DefaultServeMux)))
}
