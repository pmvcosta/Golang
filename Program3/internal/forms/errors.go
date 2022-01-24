package forms

// map of strings, with content being a slice of strings
// Slice of strings since we might have more than one error in a form
type errors map[string][]string

// Add incorporates a message and its related form field into an errors map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	// return first index item of the slice of strings
	return es[0]

}
