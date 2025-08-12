package main

// An unexported model.
type unexportedModel struct {
	Field string
}

// An unexported service.
type unexportedService struct{}

func (unexportedService) Method(unexportedModel) {}
