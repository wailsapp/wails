package optional

// True is a Bool set to true.
var True = NewBool(true)

// False is a Bool set to false.
var False = NewBool(false)

// Bool is an optional bool value.
type Bool = Var[bool]

// NewBool creates a new Bool with the given value.
func NewBool(val bool) Bool {
	return NewVar(val)
}

// Var is a generic optional value that tracks whether it has been set.
type Var[T any] struct {
	val T
	set bool
}

// Get returns the value, or the zero value if unset.
func (v *Var[T]) Get() T {
	return v.val
}

// Set assigns a value and marks the variable as set.
func (v *Var[T]) Set(val T) {
	v.val = val
	v.set = true
}

// IsSet reports whether a value has been assigned.
func (v *Var[T]) IsSet() bool {
	return v.set
}

// Unset resets the variable to the zero value and marks it as unset.
func (v *Var[T]) Unset() {
	v.set = false
	var zero T
	v.val = zero
}

// NewVar creates a new Var with the given value, marked as set.
func NewVar[T any](val T) Var[T] {
	return Var[T]{val: val, set: true}
}
