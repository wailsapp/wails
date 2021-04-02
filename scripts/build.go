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
		log.Println(string(output))
		log.Fatal(err)
	}
	fmt.Println(string(output))
}

// A build step that requires additional params, or platform specific steps for example
func main() {

	dir, _ := os.Getwd()

	// Build Runtime
	fmt.Println("**** Building Runtime ****")
	runtimeDir, _ := filepath.Abs(filepath.Join(dir, "..", "runtime", "js"))
	err := os.Chdir(runtimeDir)
	if err != nil {
		log.Fatal(err)
	}
	runCommand("npm", "install")
	runCommand("npm", "run", "build")

	// Install Wails
	fmt.Println("**** Installing Wails locally ****")
	execDir, _ := filepath.Abs(filepath.Join(dir, "..", "cmd", "wails"))
	err = os.Chdir(execDir)
	if err != nil {
		log.Fatal(err)
	}
	runCommand("go", "install")

	baseDir, _ := filepath.Abs(filepath.Join(dir, ".."))
	err = os.Chdir(baseDir)
	if err != nil {
		log.Fatal(err)
	}
	runCommand("go", "mod", "tidy")
}
