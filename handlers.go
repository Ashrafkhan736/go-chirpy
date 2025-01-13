package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	response := fmt.Sprintf(`
    <html>
      <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
      </body>
    </html>`, cfg.fileServerHits.Load())
	w.Write([]byte(response))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	headers := w.Header()
	headers.Set("content-type", "text/plain; charset=utf-8")
	cfg.fileServerHits.Store(0)
	w.Write([]byte("Hit reset to 0"))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	data, err := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: msg})
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error in parsing json %s \n", err)
		return
	}
	w.Write(data)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error in parsing json %s \n", err)
		return
	}
	w.Write(data)
}

func sliceToMap(slice []string) map[string]bool {
	result := make(map[string]bool)
	for _, item := range slice {
		result[item] = true
	}
	return result
}

func cleanChirp(msg string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	wordMap := sliceToMap(profaneWords)
	result := []string{}
	for _, word := range strings.Split(msg, " ") {
		if wordMap[strings.ToLower(word)] {
			result = append(result, "****")
		} else {
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Set("content-type", "application/json")

	chirp := struct {
		Body string `json:"body"`
	}{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&chirp)
	if err != nil {
		log.Println("Error decoding response body " + err.Error())
		respondWithError(w, 500, err.Error())
		return
	}

	if len([]rune(chirp.Body)) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	payload := struct {
		CleanedBody string `json:"cleaned_body"`
	}{}

	payload.CleanedBody = cleanChirp(chirp.Body)
	respondWithJson(w, 200, payload)
}
