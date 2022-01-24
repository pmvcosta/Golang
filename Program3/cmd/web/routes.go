package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pmvcosta/bookings/internal/config"
	"github.com/pmvcosta/bookings/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	//Example of middleware usage (chi)
	mux.Use(middleware.Recoverer)

	//Example of middleware usage (created by us)
	mux.Use(WriteToConsole)

	//Example of middleware usage (nosurf)
	mux.Use(NoSurf)

	//Use session permanence middleware
	mux.Use(SessionLoad)

	//Introduce routes (chi)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/divide", handlers.Repo.Divide)
	mux.Get("/lakeside", handlers.Repo.LakeSide)
	mux.Get("/mountains", handlers.Repo.Mountains)
	mux.Get("/sky", handlers.Repo.FlorenceSky)
	mux.Get("/reservation", handlers.Repo.Reservation)
	// Adding a post method to the reservations page
	mux.Post("/reservation", handlers.Repo.PostAvailability)
	mux.Post("/reservation-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/contact", handlers.Repo.Contact)

	// Create a file server (place to go get static files from )
	fileServer := http.FileServer(http.Dir("./static/"))

	//Use fileServer
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
