package main

import (
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 24 * time.Hour

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET env var is required")
	}
	return []byte(secret)
}

func generateToken(username string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func verifyToken(tokenStr string) bool {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return jwtSecret(), nil
	})
	return err == nil && token.Valid
}

// POST /api/v1/login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	wantUser := os.Getenv("ADMIN_USERNAME")
	wantPass := os.Getenv("ADMIN_PASSWORD")

	validUser := subtle.ConstantTimeCompare([]byte(req.Username), []byte(wantUser)) == 1
	validPass := subtle.ConstantTimeCompare([]byte(req.Password), []byte(wantPass)) == 1

	if !validUser || !validPass {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := generateToken(req.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// requireAuth wraps a handler so it only runs for requests with a valid Bearer token.
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		if tokenStr == "" || tokenStr == header || !verifyToken(tokenStr) {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next(w, r)
	}
}
