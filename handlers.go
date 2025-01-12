package main

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	headers := w.Header()
	headers.Set("content-type", "text/plain; charset=utf-8")
	w.Write([]byte("Hits: " + strconv.Itoa(int(cfg.fileServerHits.Load()))))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	headers := w.Header()
	headers.Set("content-type", "text/plain; charset=utf-8")
	cfg.fileServerHits.Store(0)
	w.Write([]byte("Hit reset to 0"))
}
