package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	log.Println("Listening on :8080")
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request URL: %s", r.URL)
		io.WriteString(w, "pong")
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
