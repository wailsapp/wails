
# Packing linux

  * create app, app.desktop, app.png (512x512)
  * chmod +x app!
  * ./linuxdeploy-x86_64.AppImage --appdir AppDir -i react.png -d react.desktop  -e react --output appimage


# Wails Doctor

Tested on:

  * Debian 8
  * Ubuntu 20.04
  * Ubuntu 19.10
  * Solus 4.1
  * Centos 8
  * Gentoo
  * OpenSUSE/leap
  * Fedora 31

### Development

Add a new package manager processor here: `v2/internal/system/packagemanager/`. IsAvailable should work even if the package is installed.
Add your new package manager to the list of package managers in `v2/internal/system/packagemanager/packagemanager.go`:

```
var db = map[string]PackageManager{
	"eopkg":  NewEopkg(),
	"apt":    NewApt(),
	"yum":    NewYum(),
	"pacman": NewPacman(),
	"emerge": NewEmerge(),
	"zypper": NewZypper(),
}
```

## Gentoo

  * Setup docker image using: emerge-webrsync -x -v
