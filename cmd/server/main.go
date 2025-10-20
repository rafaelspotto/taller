package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rafaelspotto/dlocal/cmd/server/internal/db"
	"github.com/rafaelspotto/dlocal/cmd/server/internal/handlers"
	"github.com/rafaelspotto/dlocal/cmd/server/internal/repository"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal("DB connect:", err)
	}
	defer pool.Close()

	repo := repository.NewEventRepository(pool)

	h := handlers.NewEventHandler(repo)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/events", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
	})

	adress := ":8080"

	log.Printf("Listening on %s", adress)

	if err := http.ListenAndServe(adress, r); err != nil {
		log.Fatalf("Server error %v", err)
	}

}
