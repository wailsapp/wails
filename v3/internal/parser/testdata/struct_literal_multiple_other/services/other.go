package services

// OtherService is a struct
// that does things
type OtherService struct {
	t int
}

// Yay does this and that
func (o *OtherService) Yay() []int {
	return []int{0, 1, 2}
}
