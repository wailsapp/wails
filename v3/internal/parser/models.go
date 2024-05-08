package parser

import "go/types"

// generateModels collects information about the given list of models,
// which must belong to the package described by info,
// and generates a JS/TS model file.
//
// If internal is true, the generated file is named "internal"
// and the types declared therein are not exported by the package index file.
//
// A call to index.Info.Collect must complete before entering generateModels.
func (generator *Generator) generateModels(info *PackageInfo, models []*types.TypeName, internal bool) {
	panic("model generation not implemented yet")
}
