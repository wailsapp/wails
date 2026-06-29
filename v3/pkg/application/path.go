package application

import "github.com/adrg/xdg"

type PathType int

const (
	// PathHome is the user's home directory.
	PathHome PathType = iota

	// PathDataHome defines the base directory relative to which user-specific
	// data files should be stored. This directory is defined by the
	// $XDG_DATA_HOME environment variable. If the variable is not set,
	// a default equal to $HOME/.local/share should be used.
	PathDataHome

	// PathConfigHome defines the base directory relative to which user-specific
	// configuration files should be written. This directory is defined by
	// the $XDG_CONFIG_HOME environment variable. If the variable is
	// not set, a default equal to $HOME/.config should be used.
	PathConfigHome

	// PathStateHome defines the base directory relative to which user-specific
	// state files should be stored. This directory is defined by the
	// $XDG_STATE_HOME environment variable. If the variable is not set,
	// a default equal to ~/.local/state should be used.
	PathStateHome

	// PathCacheHome defines the base directory relative to which user-specific
	// non-essential (cached) data should be written. This directory is
	// defined by the $XDG_CACHE_HOME environment variable. If the variable
	// is not set, a default equal to $HOME/.cache should be used.
	PathCacheHome

	// PathRuntimeDir defines the base directory relative to which user-specific
	// non-essential runtime files and other file objects (such as sockets,
	// named pipes, etc.) should be stored. This directory is defined by the
	// $XDG_RUNTIME_DIR environment variable. If the variable is not set,
	// applications should fall back to a replacement directory with similar
	// capabilities. Applications should use this directory for communication
	// and synchronization purposes and should not place larger files in it,
	// since it might reside in runtime memory and cannot necessarily be
	// swapped out to disk.
	PathRuntimeDir

	// PathDesktop defines the location of the user's desktop directory.
	PathDesktop

	// PathDownload defines a suitable location for user downloaded files.
	PathDownload

	// PathDocuments defines a suitable location for user document files.
	PathDocuments

	// PathMusic defines a suitable location for user audio files.
	PathMusic

	// PathPictures defines a suitable location for user image files.
	PathPictures

	// PathVideos defines a suitable location for user video files.
	PathVideos

	// PathTemplates defines a suitable location for user template files.
	PathTemplates

	// PathPublicShare defines a suitable location for user shared files.
	PathPublicShare
)

var paths = map[PathType]string{
	PathHome:        xdg.Home,
	PathDataHome:    xdg.DataHome,
	PathConfigHome:  xdg.ConfigHome,
	PathStateHome:   xdg.StateHome,
	PathCacheHome:   xdg.CacheHome,
	PathRuntimeDir:  xdg.RuntimeDir,
	PathDesktop:     xdg.UserDirs.Desktop,
	PathDownload:    xdg.UserDirs.Download,
	PathDocuments:   xdg.UserDirs.Documents,
	PathMusic:       xdg.UserDirs.Music,
	PathPictures:    xdg.UserDirs.Pictures,
	PathVideos:      xdg.UserDirs.Videos,
	PathTemplates:   xdg.UserDirs.Templates,
	PathPublicShare: xdg.UserDirs.PublicShare,
}

type PathTypes int

const (
	// PathsDataDirs defines the preference-ordered set of base directories to
	// search for data files in addition to the DataHome base directory.
	// This set of directories is defined by the $XDG_DATA_DIRS environment
	// variable. If the variable is not set, the default directories
	// to be used are /usr/local/share and /usr/share, in that order. The
	// DataHome directory is considered more important than any of the
	// directories defined by DataDirs. Therefore, user data files should be
	// written relative to the DataHome directory, if possible.
	PathsDataDirs PathTypes = iota

	// PathsConfigDirs defines the preference-ordered set of base directories
	// search for configuration files in addition to the ConfigHome base
	// directory. This set of directories is defined by the $XDG_CONFIG_DIRS
	// environment variable. If the variable is not set, a default equal
	// to /etc/xdg should be used. The ConfigHome directory is considered
	// more important than any of the directories defined by ConfigDirs.
	// Therefore, user config files should be written relative to the
	// ConfigHome directory, if possible.
	PathsConfigDirs

	// PathsFontDirs defines the common locations where font files are stored.
	PathsFontDirs

	// PathsApplicationDirs defines the common locations of applications.
	PathsApplicationDirs
)

var pathdirs = map[PathTypes][]string{
	PathsDataDirs:        xdg.DataDirs,
	PathsConfigDirs:      xdg.ConfigDirs,
	PathsFontDirs:        xdg.FontDirs,
	PathsApplicationDirs: xdg.ApplicationDirs,
}

// Path returns the path for the given selector
func Path(selector PathType) string {
	return paths[selector]
}

// Paths returns the paths for the given selector
func Paths(selector PathTypes) []string {
	return pathdirs[selector]
}
