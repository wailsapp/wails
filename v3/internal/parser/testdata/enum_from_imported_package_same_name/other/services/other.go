package services

type Title2 string

func (t Title2) String() string {
	return string(t)
}

const (
	// Mister is a title
	Mister Title2 = "Mr"
	Miss   Title2 = "Miss"
	Ms     Title2 = "Ms"
	Mrs    Title2 = "Mrs"
	Dr     Title2 = "Dr"
)
