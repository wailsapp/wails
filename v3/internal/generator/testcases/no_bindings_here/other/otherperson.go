package other

// OtherPerson is like a person, but different.
type OtherPerson[T any] struct {
	// They have a name as well.
	Name string

	// But they may have many differences.
	Differences []T
}
