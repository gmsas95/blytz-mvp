package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("💳 Payment service stub running on port 8086")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"payment"}`))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Payment service - under construction"}`))
	})

	log.Fatal(http.ListenAndServe(":8086", nil))
}