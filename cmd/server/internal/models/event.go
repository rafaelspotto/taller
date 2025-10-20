package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreateAt    time.Time `json:"created_at"`
}

type EventInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreateAt    time.Time `json:"created_at"`
}

type EventOutput struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
