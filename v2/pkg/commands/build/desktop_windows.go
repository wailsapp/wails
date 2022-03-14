//go:build windows
// +build windows

package build

// PostCompilation is called after the compilation step, if successful
func (d *DesktopBuilder) PostCompilation(options *Options) error {
	// Dump the DLLs
	//userTags := slicer.String(options.UserTags)
	//if userTags.Contains("cgo") {
	//	err := os.WriteFile(filepath.Join(options.BuildDirectory, "WebView2Loader.dll"), x64.WebView2Loader, 0755)
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

// We will compile all tray icons found at <projectdir>/assets/trayicons/*.png into the application
func (d *DesktopBuilder) processTrayIcons(assetDir string, options *Options) error {
	//
	//	var err error
	//
	//	// Get all the tray icon filenames
	//	trayIconDirectory := filepath.Join(options.ProjectData.BuildDir, "tray")
	//
	//	// If the directory doesn't exist, create it
	//	if !fs.DirExists(trayIconDirectory) {
	//		err = fs.MkDirs(trayIconDirectory)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	var trayIconFilenames []string
	//	trayIconFilenames, err = filepath.Glob(trayIconDirectory + "/*.png")
	//	if err != nil {
	//		log.Fatal(err)
	//		return err
	//	}
	//
	//	// Setup target
	//	targetFilename := "trayicons"
	//	targetFile := filepath.Join(assetDir, targetFilename+".h")
	//	d.addFileToDelete(targetFile)
	//
	//	var dataBytes []byte
	//
	//	// Use a strings builder
	//	var cdata strings.Builder
	//
	//	// Write header
	//	header := `// trayicons.h
	//// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL.
	//// This file was auto-generated. DO NOT MODIFY.
	//
	//`
	//	cdata.WriteString(header)
	//
	//	var variableList slicer.StringSlicer
	//
	//	// Loop over icons
	//	for count, filename := range trayIconFilenames {
	//
	//		// Load the tray icon
	//		dataBytes, err = ioutil.ReadFile(filename)
	//		if err != nil {
	//			return err
	//		}
	//
	//		iconname := strings.TrimSuffix(filepath.Base(filename), ".png")
	//		trayIconName := fmt.Sprintf("trayIcon%dName", count)
	//		variableList.Add(trayIconName)
	//		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", trayIconName, d.convertToHexLiteral([]byte(iconname))))
	//
	//		trayIconLength := fmt.Sprintf("trayIcon%dLength", count)
	//		variableList.Add(trayIconLength)
	//		lengthAsString := strconv.Itoa(len(dataBytes))
	//		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", trayIconLength, d.convertToHexLiteral([]byte(lengthAsString))))
	//
	//		trayIconData := fmt.Sprintf("trayIcon%dData", count)
	//		variableList.Add(trayIconData)
	//		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { ", trayIconData))
	//
	//		// Convert each byte to hex
	//		for _, b := range dataBytes {
	//			cdata.WriteString(fmt.Sprintf("0x%x, ", b))
	//		}
	//
	//		cdata.WriteString("0x00 };\n")
	//	}
	//
	//	// Write out main trayIcons data
	//	cdata.WriteString("const unsigned char *trayIcons[] = { ")
	//	cdata.WriteString(variableList.Join(", "))
	//	if len(trayIconFilenames) > 0 {
	//		cdata.WriteString(", ")
	//	}
	//	cdata.WriteString("0x00 };\n")
	//
	//	err = ioutil.WriteFile(targetFile, []byte(cdata.String()), 0600)
	//	if err != nil {
	//		return err
	//	}
	return nil
}

// compileIcon will compile the icon found at <projectdir>/icon.png into the application
func (d *DesktopBuilder) compileIcon(assetDir string, iconFile []byte) error {
	return nil
}

// We will compile all dialog icons found at <projectdir>/icons/dialog/*.png into the application
func (d *DesktopBuilder) processDialogIcons(assetDir string, options *Options) error {

	//	var err error
	//
	//	// Get all the dialog icon filenames
	//	dialogIconDirectory := filepath.Join(options.ProjectData.BuildDir, "dialog")
	//	var dialogIconFilenames []string
	//
	//	// If the directory does not exist, create it
	//	if !fs.DirExists(dialogIconDirectory) {
	//		err = fs.MkDirs(dialogIconDirectory)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	dialogIconFilenames, err = filepath.Glob(dialogIconDirectory + "/*.png")
	//	if err != nil {
	//		log.Fatal(err)
	//		return err
	//	}
	//
	//	// Setup target
	//	targetFilename := "userdialogicons"
	//	targetFile := filepath.Join(assetDir, targetFilename+".h")
	//	d.addFileToDelete(targetFile)
	//
	//	var dataBytes []byte
	//
	//	// Use a strings builder
	//	var cdata strings.Builder
	//
	//	// Write header
	//	header := `// userdialogicons.h
	//// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL.
	//// This file was auto-generated. DO NOT MODIFY.
	//
	//`
	//	cdata.WriteString(header)
	//
	//	var variableList slicer.StringSlicer
	//
	//	// Loop over icons
	//	for count, filename := range dialogIconFilenames {
	//
	//		// Load the tray icon
	//		dataBytes, err = ioutil.ReadFile(filename)
	//		if err != nil {
	//			return err
	//		}
	//
	//		iconname := strings.TrimSuffix(filepath.Base(filename), ".png")
	//		dialogIconName := fmt.Sprintf("userDialogIcon%dName", count)
	//		variableList.Add(dialogIconName)
	//		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", dialogIconName, d.convertToHexLiteral([]byte(iconname))))
	//
	//		dialogIconLength := fmt.Sprintf("userDialogIcon%dLength", count)
	//		variableList.Add(dialogIconLength)
	//		lengthAsString := strconv.Itoa(len(dataBytes))
	//		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", dialogIconLength, d.convertToHexLiteral([]byte(lengthAsString))))
	//
	//		dialogIconData := fmt.Sprintf("userDialogIcon%dData", count)
	//		variableList.Add(dialogIconData)
	//		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { ", dialogIconData))
	//
	//		// Convert each byte to hex
	//		for _, b := range dataBytes {
	//			cdata.WriteString(fmt.Sprintf("0x%x, ", b))
	//		}
	//
	//		cdata.WriteString("0x00 };\n")
	//	}
	//
	//	// Write out main dialogIcons data
	//	cdata.WriteString("const unsigned char *userDialogIcons[] = { ")
	//	cdata.WriteString(variableList.Join(", "))
	//	if len(dialogIconFilenames) > 0 {
	//		cdata.WriteString(", ")
	//	}
	//	cdata.WriteString("0x00 };\n")
	//
	//	err = ioutil.WriteFile(targetFile, []byte(cdata.String()), 0600)
	//	if err != nil {
	//		return err
	//	}
	return nil
}
