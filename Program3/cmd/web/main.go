package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pmvcosta/bookings/internal/config"
	"github.com/pmvcosta/bookings/internal/handlers"
	"github.com/pmvcosta/bookings/internal/render"
)

const portNumber = ":8080"

// This way any function in the main package can access this variable
var app config.AppConfig

var session *scs.SessionManager

// main is the main application function
func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// Assing the configured session to the app-wide configuration
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache.")
	}

	app.TemplateCache = tc
	app.UseCache = false

	// Create new repo with AppConfig app
	repo := handlers.NewRepo(&app)

	// Set the repository for the handlers
	handlers.NewHandlers(repo)

	// Set config for the template package
	render.NewTemplates(&app)

	//Using 3rd party routing package (pat)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	// Start a webserver to listen for requests, so that the previously defined
	// function can be called

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	//Start server using 3rd party routing package (pat)
	err = srv.ListenAndServe()
	log.Fatal(err)

}
