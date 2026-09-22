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

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(config.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.Timeout) * time.Second,
	}

	log.Printf("serving RSS for %s on http://%s (dir=%s)", config.Master.AuthorUsername, addr, config.Master.TargetDir)
	log.Fatal(srv.ListenAndServe())
}
