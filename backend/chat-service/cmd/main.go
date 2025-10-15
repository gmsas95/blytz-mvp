package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("💬 Chat service stub running on port 8084")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"chat"}`))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Chat service - under construction"}`))
	})

	log.Fatal(http.ListenAndServe(":8084", nil))
}