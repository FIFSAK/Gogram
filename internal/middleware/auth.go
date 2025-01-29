package middleware

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}
	return string(hashedPassword), nil
}

// ComparePasswords checks if the provided password matches the hashed password
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// RequireAuth is a middleware that ensures a valid JWT token is present in the request
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header (Bearer token)
		tokenString, err := extractTokenFromHeader(r)
		if err != nil {
			http.Error(w, "Unauthorized: Token missing or malformed", http.StatusUnauthorized)
			return
		}

		// Verify the token
		claims, err := VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Attach claims (user information) to the request context for later use in handlers
		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))

		// Proceed to the next handler if token is valid
		next(w, r)
	}
}

// extractTokenFromHeader extracts the token from the "Authorization" header
// Expected format: "Bearer <token>"
func extractTokenFromHeader(r *http.Request) (string, error) {
	// Get the "Authorization" header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	// The token should be in the format "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", fmt.Errorf("authorization header format is incorrect")
	}

	// Return the token string
	return tokenParts[1], nil
}
