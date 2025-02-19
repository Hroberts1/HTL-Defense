package api

import (
    "net/http"
    "encoding/json"
)

// Response represents a standard API response structure
type Response struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{
        Status:  "success",
        Message: "API is up and running",
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

// ExampleHandler handles example requests
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
    // Example data to return
    exampleData := map[string]string{"example": "data"}
    
    response := Response{
        Status:  "success",
        Message: "Example data retrieved",
        Data:    exampleData,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}