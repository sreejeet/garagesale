package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(Echo)
	const serveURL string = "localhost:8000"

	log.Printf("Serving at %s\n", serveURL)

	if err := http.ListenAndServe(serveURL, handler); err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
}

// Echo is a basic http handler
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Performed %s %s\n", r.Method, r.URL.Path)
}
