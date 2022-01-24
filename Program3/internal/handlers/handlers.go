package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/pmvcosta/bookings/internal/config"
	"github.com/pmvcosta/bookings/internal/forms"
	"github.com/pmvcosta/bookings/internal/models"
	"github.com/pmvcosta/bookings/internal/render"
)

//TemplateData is used by both render and handlers, thus it cannot be placed here
// It has been moved to models to avoid an import cycle error between render and handlers

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

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	//Test the implementation of sessions
	//Grab the remote IP of the person visiting the site and store in the session
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	//render.RenderTemplate(w, "home.page.tmpl") //Render the home page template, no data passed
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})

	//n, err := fmt.Fprintf(w, "This is the home page!")
	/*fmt.Println("Bytes written:", n)
	if err != nil {
	fmt.Println("There was an error:", err)
	}*/
}

func (m *Repository) LakeSide(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "lakeside.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Mountains(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "mountains.page.tmpl", &models.TemplateData{})
}

func (m *Repository) FlorenceSky(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "sky.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// Retrieving form input values!
	// Default format of retrived values is string
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte("You have successfully booked this room from " + start + " to " + end))
	//Can also use placeholders:
	w.Write([]byte(fmt.Sprintf("\nStart date is %s and end date is %s.", start, end)))

}

func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	// Create an empty instance of Reservation, to then hold form values
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	// Render page with an empty Form and Data (ready to receive info)
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	// Good practice when parsing form data
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	// Use the Reservation object created in models.go
	reservation := models.Reservation{
		FirstName: r.Form.Get("firstName"),
		LastName:  r.Form.Get("lastName"),
		Phone:     r.Form.Get("phoneNum"),
		Email:     r.Form.Get("email"),
	}

	// Create new form object (r.PostForm is of type url.values)
	form := forms.New(r.PostForm)

	// Check integrity of form data
	// Following checks add errors to Form if fields are empty...
	form.Has("firstName", r) // Check if form has a non-nil value for firstName
	form.Has("lastName", r)
	form.Has("phoneNum", r)
	form.Has("email", r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		// Re-render page with form information/data typed in by user
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}

}

//generate json interface
// variable typeOfVar name-of-field-in-JSON
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles the request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// About is the about page handler (the first thing to appear in a comment pertaining
// to a function is the function name)
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//Perform some logic with the data
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	//Pull remote IP from session
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//render.RenderTemplate(w, "about.page.tmpl") //No data passed
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap, //Assign created stringMap to struct element StringMap
	})

	/*n, err := fmt.Fprintf(w, "This is the about page!")
		fmt.Println("Bytes written:", n)
		if err != nil {
		fmt.Println("There was an error:", err)
	}

	sum := addValues(2, 7)

	fmt.Fprintf(w, fmt.Sprintf("The sum is: %d", sum))*/
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
