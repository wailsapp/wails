package backendjs

type Parser struct {
	Packages map[string]*Package
}

func NewParser() *Parser {
	return &Parser{
		Packages: make(map[string]*Package),
	}
}
