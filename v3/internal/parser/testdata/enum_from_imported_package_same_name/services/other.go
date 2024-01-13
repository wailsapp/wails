package services

type Title string

func (t Title) String() string {
	return string(t)
}

const (
	// Mister is a title
	Mister Title = "Mr"
	Miss   Title = "Miss"
	Ms     Title = "Ms"
	Mrs    Title = "Mrs"
	Dr     Title = "Dr"
)
