package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	api := &API{store: NewStore()}
	mux.HandleFunc("POST /applicants", api.createApplicant)
	mux.HandleFunc("GET /applicants", api.listApplicants)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
