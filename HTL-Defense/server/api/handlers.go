package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Response represents a standard API response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var jwtSecret = []byte("your_secret_key")

// JWTAuthenticationMiddleware authenticates requests using JWT.
func JWTAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Missing auth token", http.StatusForbidden)
			return
		}

		parts := strings.Split(tokenHeader, " ")
		if len(parts) != 2 {
			http.Error(w, "Invalid/Malformed auth token", http.StatusForbidden)
			return
		}

		tokenPart := parts[1]
		token, err := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimitingMiddleware limits the rate of incoming requests.
func RateLimitingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add actual rate limiting logic here.
		next.ServeHTTP(w, r)
	})
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

// Example secured endpoint.
func securedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Secure data accessed at %s", time.Now().String())
}

// SetupRoutes defines API endpoints.
func SetupRoutes() {
	// Secured endpoint with both JWT authentication & rate limiting.
	http.Handle("/secure", RateLimitingMiddleware(JWTAuthenticationMiddleware(http.HandlerFunc(securedEndpoint))))
	// Public endpoint example.
	http.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Public data accessed at %s", time.Now().String())
	})
}
