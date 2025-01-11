package main

import (
	"log"
	"net/http"
)

func main() {
	const port = ":8080"
	const dir = "."
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(dir))))
	mux.HandleFunc("/healthz", healthCheck)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	headers := w.Header()
	headers.Set("content-type", "text/plain; charset=utf-8")
	i, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("Error writing response", err)
	}
	log.Println("Bytes written", i)
}
