package flags

type ShowReleaseNotes struct {
	Common
	Version string `description:"The version to show the release notes for"`
}
