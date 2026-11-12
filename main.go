package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	api := &API{store: NewStore()}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("POST /applicants", api.createApplicant)
	mux.HandleFunc("GET /applicants", api.listApplicants)
	mux.HandleFunc("PATCH /applicants/{id}/stage", api.updateStage)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "route not found")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
