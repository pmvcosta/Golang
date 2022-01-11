package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pmvcosta/golangPractice1/pkg/config"
	"github.com/pmvcosta/golangPractice1/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	//Use 3rd party packages for routing
	//To get rid of installed 3rd party packages that are not being used
	// use command "go mod tidy"

	// //Create a multiplex -- an http handler (pat)
	// mux := pat.New()
	//
	// //Introduce routes (pat)
	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	// mux.Get("/divide", http.HandlerFunc(handlers.Repo.Divide))

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
