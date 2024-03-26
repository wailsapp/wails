package data

// Person holds someone's most important attributes
type Person struct {
	// Name is the person's name
	Name string `json:"name"`

	// Counts tracks the number of time the person
	// has been greeted in various ways
	Counts []int `json:"counts"`
}
