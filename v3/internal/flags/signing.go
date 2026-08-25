package flags

type SigningSetup struct {
	Platforms []string `name:"platform" description:"Comma-separated platforms to configure (darwin,windows,linux). If omitted, uses the current platform for HCL projects or detects legacy Taskfiles."`
}

type EntitlementsSetup struct {
	Output string `name:"output" description:"Output path for entitlements.plist (default: build/darwin/entitlements.plist)"`
}
