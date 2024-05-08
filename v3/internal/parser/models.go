package parser

import (
	"go/types"

	"github.com/wailsapp/wails/v3/internal/parser/collect"
)

// generateModels collects information about the given list of models,
// which must belong to the package described by info,
// and generates a JS/TS model file.
//
// If internal is true, the generated file is named "internal"
// and the types declared therein are not exported by the package index file.
//
// A call to index.Info.Collect must complete before entering generateModels.
func (generator *Generator) generateModels(info *collect.PackageInfo, models []*types.TypeName, internal bool) {
	defer generator.wg.Done()
	panic("model generation not implemented yet")
}
