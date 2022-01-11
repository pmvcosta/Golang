package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pmvcosta/golangPractice1/pkg/config"
	"github.com/pmvcosta/golangPractice1/pkg/models"
	"github.com/pmvcosta/golangPractice1/pkg/render"
)

//TemplateData is used by both render and handlers, thus it cannot be placed here
// It has been moved to models to avoid an import cycle error between render and handlers
/*
// TemplateData serves to hold any and all types of data sent to the templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} //this interface{} is used when you do not know the type of data that will be held
	CSRFToken string
	Flash     string //a flash message, indicating submission success, failure, etc
	Warning   string
	Error     string
}
*/

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Since "handlers" and "main" are separate .go fils, we must run them
// simultaneously, by doing:
// Unix:  go run *.go
// Windows: go run .   <--- emphasis on the "."

//In order for a function to respond to a request from a web browser
// it has to handle a ResponseWriter and a Request
//(m *Repository) is a receiver with type Repository that gives access to everything inside m
// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	//Test the implementation of sessions
	//Grab the remote IP of the person visiting the site and store in the session
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	//render.RenderTemplate(w, "home.page.tmpl") //Render the home page template, no data passed
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})

	n, err := fmt.Fprintf(w, "This is the home page!")
	fmt.Println("Bytes written:", n)
	//To print with placeholders:
	//fmt.Println(fmt.Sprintf("Number of of Bytes written: %d", n))

	if err != nil {
		fmt.Println("There was an error:", err)
	}
}

// About is the about page handler (the first thing to appear in a comment pertaining
// to a function is the function name)
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//Perform some logic with the data
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	//m.App.Session.Cookie

	//Pull remote IP from session
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//render.RenderTemplate(w, "about.page.tmpl") //No data passed
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap, //Assign created stringMap to struct element StringMap
	})
	n, err := fmt.Fprintf(w, "This is the about page!")
	fmt.Println("Bytes written:", n)

	//To print with placeholders:
	// %d for int
	// %s for string
	//fmt.Println(fmt.Sprintf("Number of of Bytes written: %d", n))

	if err != nil {
		fmt.Println("There was an error:", err)
	}

	sum := addValues(2, 7)

	fmt.Fprintf(w, fmt.Sprintf("The sum is: %d", sum))
}

// addValues simply adds two integers
func addValues(x, y int) int {
	return x + y
}

func (m *Repository) Divide(w http.ResponseWriter, r *http.Request) {
	var x float32 = 100.
	var y float32 = 10.
	f, err := divideValues(x, y)
	if err != nil {
		//Careful with format specifier used in error printing
		fmt.Fprintf(w, fmt.Sprintf("An error ocurred: %s", err))
		return
	}
	fmt.Fprintf(w, "Dividing %f by %f yields %f", x, y, f)
}

func divideValues(x, y float32) (float32, error) {

	if y == 0 {
		err := errors.New("Cannot divide by 0.")
		return 0, err
	}
	result := x / y
	return result, nil

}
