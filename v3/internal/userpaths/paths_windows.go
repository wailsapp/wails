package userpaths

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func folder(id *windows.KNOWNFOLDERID, fallback string) string {
	if path, err := windows.KnownFolderPath(id, windows.KF_FLAG_DONT_VERIFY); err == nil && path != "" {
		return path
	}
	return fallback
}

func platformDefaults(home string) Locations {
	home = folder(windows.FOLDERID_Profile, home)
	p := personalDefaults(home)
	local := folder(windows.FOLDERID_LocalAppData, os.Getenv("LOCALAPPDATA"))
	roaming := folder(windows.FOLDERID_RoamingAppData, os.Getenv("APPDATA"))
	shared := folder(windows.FOLDERID_ProgramData, os.Getenv("ProgramData"))
	p.DataHome, p.ConfigHome, p.StateHome = local, local, local
	p.CacheHome = filepath.Join(local, "cache")
	p.RuntimeDir = local
	p.DataDirs = uniqueAbsolute([]string{roaming, shared})
	p.ConfigDirs = uniqueAbsolute([]string{shared, roaming})
	for _, entry := range []struct {
		id   *windows.KNOWNFOLDERID
		path *string
	}{
		{windows.FOLDERID_Desktop, &p.Desktop},
		{windows.FOLDERID_Downloads, &p.Download},
		{windows.FOLDERID_Documents, &p.Documents},
		{windows.FOLDERID_Music, &p.Music},
		{windows.FOLDERID_Pictures, &p.Pictures},
		{windows.FOLDERID_Videos, &p.Videos},
	} {
		*entry.path = folder(entry.id, *entry.path)
	}
	p.Templates = folder(windows.FOLDERID_Templates, filepath.Join(roaming, "Microsoft/Windows/Templates"))
	p.PublicShare = folder(windows.FOLDERID_Public, os.Getenv("PUBLIC"))
	p.FontDirs = uniqueAbsolute([]string{folder(windows.FOLDERID_Fonts, filepath.Join(os.Getenv("SystemRoot"), "Fonts")), filepath.Join(local, "Microsoft/Windows/Fonts")})
	p.ApplicationDirs = uniqueAbsolute([]string{
		folder(windows.FOLDERID_Programs, filepath.Join(roaming, "Microsoft/Windows/Start Menu/Programs")),
		folder(windows.FOLDERID_CommonPrograms, filepath.Join(shared, "Microsoft/Windows/Start Menu/Programs")),
		folder(windows.FOLDERID_ProgramFiles, os.Getenv("ProgramFiles")),
		folder(windows.FOLDERID_ProgramFilesCommon, filepath.Join(os.Getenv("ProgramFiles"), "Common Files")),
		folder(windows.FOLDERID_UserProgramFiles, filepath.Join(local, "Programs")),
		folder(windows.FOLDERID_UserProgramFilesCommon, filepath.Join(local, "Programs/Common")),
	})
	return p
}
func loadPersonalDirectories(*Locations) {}
func (*Locations) searchDirectories()    {}
