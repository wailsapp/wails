package features

// Features holds generic and platform specific feature flags
type Features struct {
	Linux *Linux
}

// New creates a new Features object
func New() *Features {
	return &Features{
		Linux: &Linux{},
	}
}
