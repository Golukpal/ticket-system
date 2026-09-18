package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Golukpal/ticket-system/internal/middleware"
	"github.com/Golukpal/ticket-system/internal/service"
)

type TicketHandler struct {
	ticketService *service.TicketService
}

func NewTicketHandler(
	ticketService *service.TicketService,
) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
	}
}

type createTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())

	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var req createTicketRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	ticket, err := h.ticketService.Create(
		r.Context(),
		userID,
		req.Title,
		req.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTicketTitle),
			errors.Is(err, service.ErrInvalidTicketDescription):

			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})

		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
		}

		return
	}

	writeJSON(w, http.StatusCreated, ticket)
}
