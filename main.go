package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	config := GetConfig()

	mux := http.NewServeMux()
	registerRoutes(mux, &config)

	const host = "0.0.0.0"
	const port = 8066
	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(config.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.Timeout) * time.Second,
	}

	log.Printf("serving RSS for %s on http://%s (dir=%s)", config.Master.AuthorName, addr, config.Master.TargetDir)
	log.Fatal(srv.ListenAndServe())
}
