package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pmvcosta/golangPractice1/pkg/config"
	"github.com/pmvcosta/golangPractice1/pkg/handlers"
	"github.com/pmvcosta/golangPractice1/pkg/render"
)

//enable go modules:
// -Go to root directory of current project
// - type "go mod init nameOfPackage", nameOfPackage should match name of a repository
//although repository doesn't have to exist necessarily
// -This creates a go.mod file, typically not edited by hand...

//When organised in folders, where main is in cmd/web/,
// use the command "go run cmd/web/*.go" to run the program

//unlike var, const cannot be changed at will
const portNumber = ":8080"

//This way any function in the main package can access this variable
var app config.AppConfig

var session *scs.SessionManager

// main is the main application function
func main() {
	//Keeping this variable declaration here limited the scope of the functions that could access it
	//Moved outisde main() to be more acessible
	//var app config.AppConfig

	app.InProduction = false

	//When using sessions package, you most commonly need to initialize it:
	//session := scs.New() //This "session" and the one defined above main() are not the same (problem: variable shadowing), BUT
	session = scs.New() //This way both "sessions" are made to be the same (we are referring to the variable defined outside of main() )
	//(Not strictly necessary) define session duration
	session.Lifetime = 24 * time.Hour
	//If we want session to persist agter user closes the window or the browser
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	//To ensure that connection is secure, and that request is coming from https instead of http
	// Also insists that cookies be encrypted
	// Set to false here, as we are hitting localhost which is not an encrypted connection
	session.Cookie.Secure = app.InProduction

	//Assing the configured session to the app-wide configuration
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache.")
	}

	app.TemplateCache = tc
	app.UseCache = false

	//Create new repo with AppConfig app
	repo := handlers.NewRepo(&app)

	//Set the repository for the handlers
	handlers.NewHandlers(repo)

	//Set config for the template package
	render.NewTemplates(&app)

	//fmt.Println("Hello, world!")
	//http.HandleFunc(pattern, func(arg1 http.ResponseWriter, arg2 *http.Request) {  })

	//Create handler function; Still needs to be called afterwards
	//One way of implementing a function, but it is to long for the main function
	// (i.e. it is not ideal for testing)
	//"/"refers to the root level of the application, or the home page
	/*http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
	  //Fprintf returns the number of bytes written, and any error encountered
	  //var n int
	  //var err error

	  n, err := fmt.Fprintf(w, "Hello, world!")
	  fmt.Println("Bytes written:",n)
	  //To print with placeholders:
	  //fmt.Println(fmt.Sprintf("Number of of Bytes written: %d", n))

	  if err != nil {
	    fmt.Println("There was an error:", err)
	  }
	})*/

	/*
		    //Without (m *Repository) receiver
				http.HandleFunc("/", handlers.Home)
				http.HandleFunc("/about", handlers.About)
				http.HandleFunc("/divide", handlers.Divide)
	*/

	//With the (m *Repository) receiver
	//Using Go-native route creating
	/*http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)
	http.HandleFunc("/divide", handlers.Repo.Divide)*/

	//Using 3rd party routing package (pat)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	//Start a webserver to listen for requests, so that the previously defined
	//function can be called

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	//First 1024 ports of any computer are privileged, and should generally be
	// avoided for tests, as they require you to be a super user to resort to them
	//Thus we use port 8080
	//Using _="" we dump the output of ListenAndServe, whic is an error (nil if
	// no errors occur)

	//Start server using Go-native routing
	//_ = http.ListenAndServe(portNumber, nil)

	//Start server using 3rd party routing package (pat)
	err = srv.ListenAndServe()
	log.Fatal(err)

}
