package main

// An exported but internal model.
//
//wails:internal
type InternalModel struct {
	Field string
}

// An exported but internal service.
//
//wails:internal
type InternalService struct{}

func (InternalService) Method(InternalModel) {}
