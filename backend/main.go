package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"minibox/backend/auth"
	"minibox/backend/db"
	"minibox/backend/files"
	"minibox/backend/graph"
	"minibox/backend/middleware"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	if err := db.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}
	defer db.Pool.Close()

	mux := http.NewServeMux()

	// --- Rotas públicas (sem autenticação) ---
	mux.HandleFunc("POST /api/auth/login", auth.LoginHandler)
	mux.HandleFunc("POST /api/auth/register", auth.RegisterHandler)

	// --- Rotas protegidas (requerem JWT) ---
	// auth.RequireAuth é o middleware — envolve o handler como HOF
	mux.HandleFunc("GET /api/graph/v1.0/me", auth.RequireAuth(graph.MeHandler))

	mux.HandleFunc("GET /api/files", auth.RequireAuth(files.ListHandler))
	mux.HandleFunc("POST /api/files/upload", auth.RequireAuth(files.UploadHandler))
	mux.HandleFunc("GET /api/files/download", auth.RequireAuth(files.DownloadHandler))

	// Aplica o middleware CORS em TODAS as rotas
	handler := middleware.CORS(mux)

	fmt.Println("Backend rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
