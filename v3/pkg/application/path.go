package application

import (
	"sync"

	"github.com/wailsapp/wails/v3/internal/userpaths"
)

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

// The user's folders are resolved on the first call to Path or Paths, not at
// package init: on Windows that asks the shell (shell32, and user32 and COM
// with it), which a program that never asks should not pay for. The snapshot
// is taken then, so environment changes made before that call are seen.
var resolvedPaths = sync.OnceValue(userpaths.Resolve)

var paths = sync.OnceValue(func() map[PathType]string {
	p := resolvedPaths()
	return map[PathType]string{
		PathHome:        p.Home,
		PathDataHome:    p.DataHome,
		PathConfigHome:  p.ConfigHome,
		PathStateHome:   p.StateHome,
		PathCacheHome:   p.CacheHome,
		PathRuntimeDir:  p.RuntimeDir,
		PathDesktop:     p.Desktop,
		PathDownload:    p.Download,
		PathDocuments:   p.Documents,
		PathMusic:       p.Music,
		PathPictures:    p.Pictures,
		PathVideos:      p.Videos,
		PathTemplates:   p.Templates,
		PathPublicShare: p.PublicShare,
	}
})

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

var pathdirs = sync.OnceValue(func() map[PathTypes][]string {
	p := resolvedPaths()
	return map[PathTypes][]string{
		PathsDataDirs:        p.DataDirs,
		PathsConfigDirs:      p.ConfigDirs,
		PathsFontDirs:        p.FontDirs,
		PathsApplicationDirs: p.ApplicationDirs,
	}
})

// Path returns the path for the given selector
func Path(selector PathType) string {
	return paths()[selector]
}

// Paths returns the paths for the given selector
func Paths(selector PathTypes) []string {
	return pathdirs()[selector]
}
