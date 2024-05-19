package flags

type ServiceCreateOptions struct {
	// Name of the service
	Name string `description:"Name of the service" required:"true"`
	// Description of the service
	Description string `description:"Description of the service" required:"true"`
	// Version
	Version string `description:"Version of the service"`
	// URL
	URL string `description:"URL of the service"`
	// Path
	OutputDir string `description:"Path to output the service" default:"./services"`
}
