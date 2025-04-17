package service

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: Check for a valid Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !isValidToken(authHeader) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call the next handler if authentication is successful
		next.ServeHTTP(w, r)
	})
}

// isValidToken is a placeholder for your token validation logic
func isValidToken(token string) bool {
	// Implement your token validation logic here
	return token == "valid-token"
}
