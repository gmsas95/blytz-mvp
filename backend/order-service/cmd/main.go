package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("📦 Order service stub running on port 8085")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"order"}`))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Order service - under construction"}`))
	})

	log.Fatal(http.ListenAndServe(":8085", nil))
}