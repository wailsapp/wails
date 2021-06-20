# Build Directory

The build directory is used to house all the build files and assets for your application. 

The structure is:

  * bin - Output directory
  * dialog - Icons for dialogs
  * tray - Icons for the system tray
  * mac - MacOS specific files
  * linux - Linux specific files
  * windows - Windows specific files

## Dialog Icons

Place any PNG file in this directory to be able to use them in message dialogs.
The files should have names in the following format: `name[-(light|dark)][2x].png`

Examples:

* `mypic.png` - Standard definition icon with ID `mypic` 
* `mypic-light.png` - Standard definition icon with ID `mypic`, used when system theme is light  
* `mypic-dark.png` - Standard definition icon with ID `mypic`, used when system theme is dark
* `mypic2x.png` - High definition icon with ID `mypic`
* `mypic-light2x.png` - High definition icon with ID `mypic`, used when system theme is light
* `mypic-dark2x.png` - High definition icon with ID `mypic`, used when system theme is dark

### Order of preference

Icons are selected with the following order of preference:

For High Definition displays:
* name-(theme)2x.png
* name2x.png
* name-(theme).png
* name.png
  
For Standard Definition displays:
* name-(theme).png
* name.png

## Tray

Place any PNG file in this directory to be able to use them as tray icons.
The name of the filename will be the ID to reference the image.

Example:

* `mypic.png` - May be referenced using `runtime.Tray.SetIcon("mypic")` 

## Mac

The `darwin` directory holds files specific to Mac builds, such as `Info.plist`. 
These may be customised and used as part of the build. To return these files to the default state, simply delete them and
build with the `-package` flag.

## Windows 

The `windows` directory contains the manifest and rc files used when building with the `-package` flag. 
These may be customised for your application. To return these files to the default state, simply delete them and
build with the `-package` flag.