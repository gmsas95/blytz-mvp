package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("🚀 Gateway service stub running on port 8080")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"gateway"}`))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Gateway service - under construction"}`))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}