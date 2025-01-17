package collect

import (
	"slices"
	"strings"
)

// PackageIndex lists all services, models and unexported models
// to be generated from a single package.
//
// When obtained through a call to [PackageInfo.Index],
// each service and model appears at most once;
// services are sorted by name;
// exported models precede all unexported ones
// and both ranges are sorted by name.
type PackageIndex struct {
	Package *PackageInfo

	Services []*ServiceInfo

	Models            []*ModelInfo
	HasExportedModels bool // If true, there is at least one exported model.
}

// Index computes a [PackageIndex] for the selected language from the list
// of generated services and models and regenerates cached stats.
//
// Services and models appear at most once in the returned slices;
// services are sorted by name;
// exported models precede all unexported ones
// and both ranges are sorted by name.
//
// Index calls info.Collect, and therefore provides the same guarantees.
// It is safe for concurrent use.
func (info *PackageInfo) Index(TS bool) (index *PackageIndex) {
	// Init index.
	index = &PackageIndex{
		Package: info.Collect(),
	}

	// Init stats
	stats := &Stats{
		NumPackages: 1,
	}

	// Gather services.
	for _, value := range info.services.Range {
		service := value.(*ServiceInfo)
		if !service.IsEmpty() {
			if service.Object().Exported() {
				// Publish non-internal service on the local index.
				index.Services = append(index.Services, service)
			}
			// Update service stats.
			stats.NumServices++
			stats.NumMethods += len(service.Methods)
		}
	}

	// Sort services by name.
	slices.SortFunc(index.Services, func(b1 *ServiceInfo, b2 *ServiceInfo) int {
		if b1 == b2 {
			return 0
		}
		return strings.Compare(b1.Name, b2.Name)
	})

	// Gather models.
	for _, value := range info.models.Range {
		model := value.(*ModelInfo)
		index.Models = append(index.Models, model)
		// Mark presence of exported models
		if model.Object().Exported() {
			index.HasExportedModels = true
		}
		// Update model stats.
		if len(model.Values) > 0 {
			stats.NumEnums++
		} else {
			stats.NumModels++
		}
	}

	// Sort models by exported property (exported first), then by name.
	slices.SortFunc(index.Models, func(m1 *ModelInfo, m2 *ModelInfo) int {
		if m1 == m2 {
			return 0
		}

		m1e, m2e := m1.Object().Exported(), m2.Object().Exported()
		if m1e != m2e {
			if m1e {
				return -1
			} else {
				return 1
			}
		}

		return strings.Compare(m1.Name, m2.Name)
	})

	// Cache stats
	info.stats.Store(stats)

	return
}

// IsEmpty returns true if the given index
// contains no data for the selected language.
func (index *PackageIndex) IsEmpty() bool {
	return len(index.Package.Injections) == 0 && len(index.Services) == 0 && len(index.Models) == 0
}
