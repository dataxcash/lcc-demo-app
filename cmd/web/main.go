package main

import (
	"log"
	"net/http"

	"demo-app/internal/web"
)

func main() {
	// Start minimal Web UI + API server
	srv := web.NewServer()

	addr := ":9144" // default web ui port
	log.Printf("LCC Demo Web UI listening on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, srv.Router()); err != nil {
		log.Fatalf("web server error: %v", err)
	}
}
