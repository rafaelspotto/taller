package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rafaelspotto/dlocal/cmd/server/internal/models"
	"github.com/rafaelspotto/dlocal/cmd/server/internal/repository"
)

type EventHandler struct {
	repo repository.EventRepository
}

func NewEventHandler(repo repository.EventRepository) *EventHandler {
	return &EventHandler{repo: repo}
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.EventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if input.StartDate.IsZero() {
		http.Error(w, "Start date is required", http.StatusBadRequest)
		return
	}

	if input.EndDate.IsZero() {
		http.Error(w, "End date is required", http.StatusBadRequest)
		return
	}

	if input.EndDate.Before(input.StartDate) {
		http.Error(w, "End date must be after start date", http.StatusBadRequest)
		return
	}

	event := models.Event{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		CreateAt:    time.Now(),
	}

	createdEvent, err := h.repo.Create(r.Context(), event)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdEvent)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := h.repo.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	events, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(events) {
		events = []models.Event{}
	} else {
		if end > len(events) {
			end = len(events)
		}
		events = events[start:end]
	}

	response := map[string]interface{}{
		"events": events,
		"page":   page,
		"limit":  limit,
		"total":  len(events),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
