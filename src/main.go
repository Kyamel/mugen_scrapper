package main

import (
	"log"
	"mugen_scraper/src/routes"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/animes", routes.GetAnimes)
	r.Get("/episodes", routes.GetEpisodes)
	r.Get("/platforms", routes.GetPlatforms)
	r.Get("/songs", routes.GetSongs)

	return r
}
