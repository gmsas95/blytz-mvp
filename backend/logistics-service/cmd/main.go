package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("🚚 Logistics service stub running on port 8087")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"logistics"}`))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Logistics service - under construction"}`))
	})

	log.Fatal(http.ListenAndServe(":8087", nil))
}