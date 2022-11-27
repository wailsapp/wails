package linux

// Options specific to Linux builds
type Options struct {
	Icon                []byte
	WindowIsTranslucent bool

	// User messages that can be customised
	Messages *Messages
}

type Messages struct {
	WebKit2GTKMinRequired string
}

func DefaultMessages() *Messages {
	return &Messages{
		WebKit2GTKMinRequired: "This application requires at least WebKit2GTK %s to be installed.",
	}
}
