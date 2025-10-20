package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rafaelspotto/dlocal/cmd/server/internal/models"
)

type EventRepository interface {
	Create(ctx context.Context, event models.Event) (models.Event, error)
	Get(ctx context.Context, id uuid.UUID) (models.Event, error)
	List(ctx context.Context) ([]models.Event, error)
	Update(ctx context.Context, id uuid.UUID, event models.Event) (models.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type eventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event models.Event) (models.Event, error) {
	const q = `INSERT INTO events (id, title, description, start_time, end_time, created_at)
VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(ctx, q, event.ID, event.Title, event.Description, event.StartDate, event.EndDate, event.CreateAt)
	return event, err
}

func (r *eventRepository) Get(ctx context.Context, id uuid.UUID) (models.Event, error) {
	const q = `SELECT id, title, description, start_time, end_time, created_at FROM events WHERE id=$1`
	row := r.db.QueryRow(ctx, q, id)
	var e models.Event
	var desc *string
	if err := row.Scan(&e.ID, &e.Title, &desc, &e.StartDate, &e.EndDate, &e.CreateAt); err != nil {
		return models.Event{}, err
	}
	if desc != nil {
		e.Description = *desc
	}
	return e, nil
}

func (r *eventRepository) List(ctx context.Context) ([]models.Event, error) {
	const q = `SELECT id, title, description, start_time, end_time, created_at FROM events ORDER BY start_time`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Event
	for rows.Next() {
		var e models.Event
		var desc *string
		if err := rows.Scan(&e.ID, &e.Title, &desc, &e.StartDate, &e.EndDate, &e.CreateAt); err != nil {
			return nil, err
		}
		if desc != nil {
			e.Description = *desc
		}
		out = append(out, e)
	}
	return out, nil
}

func (r *eventRepository) Update(ctx context.Context, id uuid.UUID, event models.Event) (models.Event, error) {
	const q = `UPDATE events SET title=$2, description=$3, start_time=$4, end_time=$5 WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id, event.Title, event.Description, event.StartDate, event.EndDate)
	if err != nil {
		return models.Event{}, err
	}
	event.ID = id
	return event, nil
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM events WHERE id=$1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}
