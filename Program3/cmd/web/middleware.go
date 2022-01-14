package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// WriteToConsole writes a message to the console. Mainly serves to show how to implement our own middleware
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Hit the Page.")
		next.ServeHTTP(w, r)
	})
}

//Web Servers, by their very nature are not state aware, need to add some middleware
// to remember state
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
	// LoadAndSave automatically loads and saves session data for the current
	// request, and communicates the session token to and from the client in a
	// cookie
}

//Using the 3rd party middleware nosurf to prevent csrf attacks
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next) //creates a handler

	//Creating a cookie "Path: "/" " refers to the entire site
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler

}
