package main

import (
	"fmt"
	"github.com/leaanthony/slicer"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

func convertToHexLiteral(bytes []byte) string {
	result := ""
	for _, b := range bytes {
		result += fmt.Sprintf("0x%x, ", b)
	}
	return result
}

func main() {
	dialogIconFilenames, err := filepath.Glob("*.png")
	if err != nil {
		log.Fatal(err)
	}

	// Build icons for Mac
	err = buildMacIcons(dialogIconFilenames)
	if err != nil {
		log.Fatal(err)
	}
}

func buildMacIcons(dialogIconFilenames []string) error {

	// Setup target
	targetFile := "../../../../../../../internal/ffenestri/defaultdialogicons_darwin.c"

	var dataBytes []byte
	var err error

	// Use a strings builder
	var cdata strings.Builder

	// Write header
	header := `// defaultdialogicons_darwin.c
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL.
// This file was auto-generated. DO NOT MODIFY.

`
	cdata.WriteString(header)

	var variableList slicer.StringSlicer

	// Loop over icons
	for count, filename := range dialogIconFilenames {

		// Load the tray icon
		dataBytes, err = ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		iconname := strings.TrimSuffix(filepath.Base(filename), ".png")
		dialogIconName := fmt.Sprintf("defaultDialogIcon%dName", count)
		variableList.Add(dialogIconName)
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", dialogIconName, convertToHexLiteral([]byte(iconname))))

		dialogIconLength := fmt.Sprintf("defaultDialogIcon%dLength", count)
		variableList.Add(dialogIconLength)
		lengthAsString := strconv.Itoa(len(dataBytes))
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", dialogIconLength, convertToHexLiteral([]byte(lengthAsString))))

		dialogIconData := fmt.Sprintf("defaultDialogIcon%dData", count)
		variableList.Add(dialogIconData)
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { ", dialogIconData))

		// Convert each byte to hex
		for _, b := range dataBytes {
			cdata.WriteString(fmt.Sprintf("0x%x, ", b))
		}

		cdata.WriteString("0x00 };\n")
	}

	// Write out main dialogIcons data
	cdata.WriteString("const unsigned char *defaultDialogIcons[] = { ")
	cdata.WriteString(variableList.Join(", "))
	if len(dialogIconFilenames) > 0 {
		cdata.WriteString(", ")
	}
	cdata.WriteString("0x00 };\n")

	err = ioutil.WriteFile(targetFile, []byte(cdata.String()), 0600)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
