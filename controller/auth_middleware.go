package controller

import (
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret_key") // ganti dengan key rahasia

// AuthMiddleware memastikan hanya user dengan JWT valid yang bisa akses endpoint
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid Authorization Header", http.StatusUnauthorized)
            return
        }

        tokenStr := parts[1]

        // Verifikasi JWT
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Invalid Token", http.StatusUnauthorized)
            return
        }

        // Kalau token valid, teruskan ke handler berikutnya
        next.ServeHTTP(w, r)
    })
}
