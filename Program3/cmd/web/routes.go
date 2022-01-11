package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pmvcosta/bookings/pkg/config"
	"github.com/pmvcosta/bookings/pkg/handlers"
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

	return mux
}
