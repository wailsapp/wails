package flags

// ToolMSIX represents the options for the MSIX packaging command
type ToolMSIX struct {
	Common

	// Project configuration
	ConfigPath string `name:"config" description:"Path to the project configuration file" default:"wails.json"`

	// MSIX package information
	Publisher string `name:"publisher" description:"Publisher name for the MSIX package (e.g., CN=CompanyName)" default:""`

	// Certificate for signing
	CertificatePath     string `name:"cert" description:"Path to the certificate file for signing the MSIX package" default:""`
	CertificatePassword string `name:"cert-password" description:"Password for the certificate file" default:""`

	// Build options
	Arch           string `name:"arch" description:"Architecture of the package (x64, x86, arm64)" default:"x64"`
	ExecutableName string `name:"name" description:"Name of the executable in the package" default:""`
	ExecutablePath string `name:"executable" description:"Path to the executable file to package" default:""`
	OutputPath     string `name:"out" description:"Path where the MSIX package will be saved" default:""`

	// Tool selection
	UseMsixPackagingTool bool `name:"use-msix-tool" description:"Use the Microsoft MSIX Packaging Tool for packaging" default:"false"`
	UseMakeAppx          bool `name:"use-makeappx" description:"Use MakeAppx.exe for packaging" default:"true"`
}
