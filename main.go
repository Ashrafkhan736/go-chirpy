package main

import (
	"log"
	"net/http"
)

func main() {
	const port = ":8080"
	const dir = "."
	apiMertic := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiMertic.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(dir)))))
	mux.HandleFunc("/healthz", healthCheck)
	mux.HandleFunc("/metrics", apiMertic.getMetrics)
	mux.HandleFunc("/reset", apiMertic.resetMetrics)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
