package analyser

import (
	"go/types"
	"slices"

	"github.com/pterm/pterm"
)

// target represents an additional discovered target for static analysis.
// This may be (a subfield of) either a local or a global variable.
type target struct {
	// pkgi holds the index (relative to [Analyser.pkgs]) of the package
	// where a local or unexported variable lives.
	// For exported global variables, pkgi is -1 and the analysis
	// is performed over all packages.
	pkgi int

	// variable holds the type-checker object describing the target variable.
	variable *types.Var

	// path holds the path to a subfield of the target variable.
	path PathSetElement
}

// schedule adds the given target to the analyser queue
// only if it has never been scheduled before.
func (analyser *Analyser) schedule(tgt target) {
	if analyser.scheduled.Add(tgt) {
		analyser.queue = append(analyser.queue, tgt)

		// Target analysis with a strong path yields
		// the same or more results than analysis of the same target
		// with the weak version of the same path.
		// If we are queueing the target with a strong path,
		// mark the weak version as done too.
		if tgt.path.StrongRef() {
			tgt.path = analyser.paths.Get(tgt.path.Path().Weaken())
			if !analyser.scheduled.Add(tgt) {
				analyser.queue = slices.DeleteFunc(analyser.queue, func(t target) bool {
					return t == tgt
				})
			}
		}
	}
}

// processTarget finds all references to a local or global target variable
// and feeds them to [Analyser.processReference] for further analysis.
func (analyser *Analyser) processTarget(tgt target) {
	pkgstart, pkgend := 0, len(analyser.pkgs)

	if tgt.pkgi >= 0 {
		// tgt is either a local variable or an unexported global,
		// restrict analysis to the defining package.
		pkgstart, pkgend = tgt.pkgi, min(len(analyser.pkgs), tgt.pkgi+1)
	}

	defInPkgs := false

	// Iterate over selected package (either exactly one, or all of them).
	for pkgi := pkgstart; pkgi < pkgend; pkgi++ {
		if tgt.variable.Pkg() == analyser.pkgs[pkgi].Types {
			defInPkgs = true
		}

		// Analyse references to the target variable in the i-th package.
		for _, ref := range analyser.refs[pkgi][tgt.variable] {
			analyser.processReference(pkgi, tgt.variable, ref, tgt.path.Path())
		}
	}

	if tgt.pkgi < 0 && !defInPkgs {
		// The defining package for an exported global target
		// is not under analysis: emit a warning.
		pterm.Warning.Printfln(
			"global %s might provide bound services, but its declaring package is not under analysis",
			types.ObjectString(tgt.variable, nil),
		)
	}
}
