package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/api", apiHandler)

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to HTL-Defense API"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API endpoint"))
}