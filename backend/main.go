package main

import (
	"log"
	"net/http"
	"os"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db := connectDB()
	srv := &Server{db: db}

	mux := http.NewServeMux()

	// Diakses tanpa login: firmware ESP8266 + endpoint login itu sendiri.
	mux.HandleFunc("GET /healthz", srv.handleHealth)
	mux.HandleFunc("POST /api/v1/access", srv.handleAccess)
	mux.HandleFunc("GET /api/v1/unlock/poll", srv.handlePollUnlock)
	mux.HandleFunc("POST /api/v1/login", handleLogin)

	// Butuh login (dashboard).
	mux.HandleFunc("GET /api/v1/cards", requireAuth(srv.handleListCards))
	mux.HandleFunc("POST /api/v1/cards", requireAuth(srv.handleCreateCard))
	mux.HandleFunc("DELETE /api/v1/cards/{id}", requireAuth(srv.handleDeleteCard))
	mux.HandleFunc("GET /api/v1/logs", requireAuth(srv.handleListLogs))
	mux.HandleFunc("POST /api/v1/unlock", requireAuth(srv.handleUnlock))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
