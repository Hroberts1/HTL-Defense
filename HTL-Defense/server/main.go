package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    // Set up routes
    r.HandleFunc("/api/endpoint", YourHandlerFunction).Methods("GET")

    // Start the server
    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}

func YourHandlerFunction(w http.ResponseWriter, r *http.Request) {
    // Handle your request here
    w.Write([]byte("Hello, HTL-Defense!"))
}