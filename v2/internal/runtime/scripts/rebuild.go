package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/internal/fs"
	"github.com/wailsapp/wails/v2/internal/shell"
)

func main() {

	sourceDir := fs.RelativePath("../../../internal/runtime/js")

	platforms := []string{"darwin"}

	for _, platform := range platforms {
		println("Building JS Runtime for " + platform)
		envvars := []string{"WAILSPLATFORM=" + platform}

		// Split up the InstallCommand and execute it
		stdout, stderr, err := shell.RunCommand(sourceDir, "npm", "install")
		if err != nil {
			for _, l := range strings.Split(stdout, "\n") {
				fmt.Printf("    %s\n", l)
			}
			for _, l := range strings.Split(stderr, "\n") {
				fmt.Printf("    %s\n", l)
			}
		}

		runtimeDir := fs.RelativePath("../js")
		cmd := shell.CreateCommand(runtimeDir, "npm", "run", "build:desktop")
		cmd.Env = append(os.Environ(), envvars...)
		var stdo, stde bytes.Buffer
		cmd.Stdout = &stdo
		cmd.Stderr = &stde
		err = cmd.Run()
		if err != nil {
			for _, l := range strings.Split(stdo.String(), "\n") {
				fmt.Printf("    %s\n", l)
			}
			for _, l := range strings.Split(stde.String(), "\n") {
				fmt.Printf("    %s\n", l)
			}
			log.Fatal(err)
		}

		wailsJS := fs.RelativePath("../../../internal/runtime/assets/desktop_" + platform + ".js")
		runtimeData, err := ioutil.ReadFile(wailsJS)
		if err != nil {
			log.Fatal(err)
		}

		// Convert to C structure
		runtimeC := `
// runtime.c (c) 2019-Present Lea Anthony.
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file was auto-generated. DO NOT MODIFY.
const unsigned char runtime[]={`
		for _, b := range runtimeData {
			runtimeC += fmt.Sprintf("0x%x, ", b)
		}
		runtimeC += "0x00};"

		// Save file
		outputFile := fs.RelativePath(fmt.Sprintf("../../ffenestri/runtime_%s.c", platform))

		if err := ioutil.WriteFile(outputFile, []byte(runtimeC), 0600); err != nil {
			log.Fatal(err)
		}
	}
}
