package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mperuch/quote-engine/internal/quote"
)

func main() {
	path := os.Getenv("TABLES_PATH")
	if path == "" {
		path = "data/tables.json"
	}

	eng := quote.New(path, 30*time.Second)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /v1/quotes", func(w http.ResponseWriter, r *http.Request) {
		var req quote.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		res, err := eng.Quote(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})

	addr := ":8080"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}
	log.Printf("quote-engine listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
