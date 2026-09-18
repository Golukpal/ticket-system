package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Golukpal/ticket-system/internal/model"
	"github.com/Golukpal/ticket-system/internal/repository"
)

var (
	ErrInvalidTicketTitle       = errors.New("title is required")
	ErrInvalidTicketDescription = errors.New("description is required")
)

type TicketService struct {
	tickets *repository.TicketRepository
}

func NewTicketService(
	tickets *repository.TicketRepository,
) *TicketService {
	return &TicketService{
		tickets: tickets,
	}
}

func (s *TicketService) Create(
	ctx context.Context,
	userID string,
	title string,
	description string,
) (model.Ticket, error) {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		return model.Ticket{}, ErrInvalidTicketTitle
	}

	if description == "" {
		return model.Ticket{}, ErrInvalidTicketDescription
	}

	return s.tickets.Create(
		ctx,
		userID,
		title,
		description,
	)
}