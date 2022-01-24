package models

import (
	"github.com/pmvcosta/bookings/internal/forms"
)

// TemplateData serves to hold any and all types of data sent to the templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} //this interface{} is used when you do not know the type of data that will be held
	CSRFToken string                 //Cross Site Request Forgery Token
	Flash     string                 //a flash message, indicating submission success, failure, etc
	Warning   string
	Error     string
	Form      *forms.Form //Available to every page regardless of whether they have forms or not
}
