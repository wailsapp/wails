package build

// DesktopBuilder builds applications for the desktop
type DesktopBuilder struct {
	*BaseBuilder
}

func newDesktopBuilder(options *Options) *DesktopBuilder {
	return &DesktopBuilder{
		BaseBuilder: NewBaseBuilder(options),
	}
}
