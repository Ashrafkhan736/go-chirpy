package main

import (
	// "fmt"
	"net/http"
)

func main() {
	handler := http.NewServeMux()
	sever := http.Server{}
	sever.Addr = ":8080"
	sever.Handler = handler
	handler.Handle("/", http.FileServer(http.Dir(".")))
	sever.ListenAndServe()
}
