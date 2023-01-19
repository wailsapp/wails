package flags

type Update struct {
	Common
	Version    string `description:"The version to update to"`
	PreRelease bool   `name:"pre" description:"Update to latest pre-release"`
}
