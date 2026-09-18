package repository

import (
	"context"
	"fmt"

	"github.com/Golukpal/ticket-system/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketRepository struct {
	db *pgxpool.Pool
}

func NewTicketRepository(db *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) Create(
	ctx context.Context,
	userID string,
	title string,
	description string,
) (model.Ticket, error) {
	var ticket model.Ticket

	query := `
		INSERT INTO tickets (
			user_id,
			title,
			description
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			user_id,
			title,
			description,
			status,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		title,
		description,
	).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.Title,
		&ticket.Description,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)

	if err != nil {
		return model.Ticket{}, fmt.Errorf("create ticket: %w", err)
	}

	return ticket, nil
}
