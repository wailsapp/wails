# Assets Directory

The assets directory is used to house all the assets of your application. 

The structure is:

  * dialog - Icons for dialogs
  * tray - Icons for the system tray
  * custom - A place for assets you wish to bundle in the application
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

## Custom

Any file in this directory will be embedded into the app using the Wails asset bundler.
Assets can be retrieved using the following methods:

* `wails.Assets().Read(filename string) ([]byte, error)`  
* `wails.Assets().String(filename string) (string, error)`

The filename should include the path to the file relative to the `custom` directory.