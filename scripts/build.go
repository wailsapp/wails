package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

/*
# Build runtime
echo "**** Building Runtime ****"
cd runtime/js
npm install
npm run build
cd ../..

echo "**** Packing Assets ****"
cd cmd
mewn
cd ..
cd lib/renderer
mewn
cd ../..

echo "**** Installing Wails locally ****"
cd cmd/wails
go install
cd ../..

echo "**** Tidying the mods! ****"
go mod tidy

echo "**** WE ARE DONE! ****"

*/

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		log.Fatal(err)
	}
	cmd.Run()
	fmt.Println(string(output))
}

// A build step that requires additional params, or platform specific steps for example
func main() {

	dir, _ := os.Getwd()

	// Build Runtime
	fmt.Println("**** Building Runtime ****")
	runtimeDir, _ := filepath.Abs(filepath.Join(dir, "..", "runtime", "js"))
	os.Chdir(runtimeDir)
	runCommand("npm", "install")
	runCommand("npm", "run", "build")

	// Pack assets
	fmt.Println("**** Packing Assets ****")
	rendererDir, _ := filepath.Abs(filepath.Join(dir, "..", "lib", "renderer"))
	os.Chdir(rendererDir)
	runCommand("mewn")
	cmdDir, _ := filepath.Abs(filepath.Join(dir, "..", "cmd"))
	os.Chdir(cmdDir)
	runCommand("mewn")

	// Install Wails
	fmt.Println("**** Installing Wails locally ****")
	execDir, _ := filepath.Abs(filepath.Join(dir, "..", "cmd", "wails"))
	os.Chdir(execDir)
	runCommand("go", "install")

	baseDir, _ := filepath.Abs(filepath.Join(dir, ".."))
	os.Chdir(baseDir)
	runCommand("go", "mod", "tidy")
}
