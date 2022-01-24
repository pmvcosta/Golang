package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New creates a new form and returns a pointer to it
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
	// since we are declaring it to be empty, we must use string{}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		// Need to add errors at some point, e.g. here:
		f.Errors.Add(field, "This field cannot be blank.")
		return false
	}
	return true
}

// Valid checks the validity of the data of a Form
func (f *Form) Valid() bool {
	// returns true if there are no errors
	return len(f.Errors) == 0
}
