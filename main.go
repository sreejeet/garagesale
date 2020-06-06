package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(Echo)

	if err := http.ListenAndServe("localhost:8000", handler); err != nil {
		log.Fatal("ERROR: %s\n", err)
	}
}

// Echo is a basic http handler
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Performed %s %s\n", r.Method, r.URL.Path)
}
