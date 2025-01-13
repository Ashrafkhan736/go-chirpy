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
	mux.HandleFunc("GET /api/healthz", healthCheck)
	mux.HandleFunc("GET /admin/metrics", apiMertic.getMetrics)
	mux.HandleFunc("POST /admin/reset", apiMertic.resetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
