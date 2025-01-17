package nobindingshere

import "github.com/wailsapp/wails/v3/internal/generator/testcases/no_bindings_here/other"

// Person is not a number.
type Person struct {
	// They have a name.
	Name    string
	Friends [4]Impersonator // Exactly 4 sketchy friends.
}

// Impersonator gets their fields from other people.
type Impersonator other.OtherPerson[int]

// HowDifferent is a curious kind of person
// that lets other people decide how they are different.
type HowDifferent[How any] other.OtherPerson[map[string]How]

// PrivatePerson gets their fields from hidden sources.
type PrivatePerson personImpl

type personImpl struct {
	// Nickname conceals a person's identity.
	Nickname string
	Person
}
