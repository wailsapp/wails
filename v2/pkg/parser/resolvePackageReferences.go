package parser

import (
	"fmt"

	"github.com/leaanthony/slicer"
)

// resolvePackageNames will deterine the names for the packages, allowing
// us to create a flat structure for the imports in the frontend module
func (p *Parser) resolvePackageNames() {

	// A cache for the names
	var packageNameCache slicer.StringSlicer

	// Process each package
	for _, pkg := range p.packages {
		pkgName := pkg.gopackage.Name

		// Check for collision
		if packageNameCache.Contains(pkgName) {
			// https://www.youtube.com/watch?v=otNNGROI0Cs !!!!!

			// We start at 2 because having both "pkg" and "pkg1" is 🙄
			count := 2
			for ok := true; ok; ok = packageNameCache.Contains(pkgName) {
				pkgName = fmt.Sprintf("%s%d", pkg.gopackage.Name, count)
			}
		}

		// Save the name!
		packageNameCache.Add(pkgName)
		pkg.Name = pkgName
	}
}
