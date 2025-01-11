package main

import (
	"log"
	"net/http"
)

func main() {
	const port = ":8080"
	const dir = "."
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(dir)))

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
