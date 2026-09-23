package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Golukpal/ticket-system/internal/handler"
	"github.com/Golukpal/ticket-system/internal/middleware"
)

func Setup(
	authHandler *handler.AuthHandler,
	ticketHandler *handler.TicketHandler,
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health",
		healthHandler,
	)

	mux.HandleFunc(
		"POST /auth/register",
		authHandler.Register,
	)

	mux.HandleFunc(
		"POST /auth/login",
		authHandler.Login,
	)

	mux.Handle(
		"POST /tickets",
		authMiddleware.RequireAuth(
			http.HandlerFunc(ticketHandler.Create),
		),
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