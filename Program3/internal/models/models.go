package models

// Reservation holds reservation data (eventually to be stored in a database)
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}
