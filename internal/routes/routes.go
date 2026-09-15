package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Golukpal/ticket-system/internal/handler"
)

func Setup(authHandler *handler.AuthHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc(
		"POST /auth/register",
		authHandler.Register,
	)

	mux.HandleFunc(
		"POST /auth/login",
		authHandler.Login,
	)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}