package main

import (
    "log"
    "net/http"
)

func main() {
    // Set up routes and handlers
    http.HandleFunc("/", homeHandler)

    // Start the server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to the HTL-Defense Cybersecurity Platform!"))
}