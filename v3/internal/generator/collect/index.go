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
// services and models are sorted
// by internal property (all exported first), then by name.
type PackageIndex struct {
	Package *PackageInfo

	Services            []*ServiceInfo
	HasExportedServices bool // If true, there is at least one exported service.

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
			index.Services = append(index.Services, service)
			// Mark presence of exported services
			if !service.Internal {
				index.HasExportedServices = true
			}
			// Update service stats.
			stats.NumServices++
			stats.NumMethods += len(service.Methods)
		}
	}

	// Sort services by internal property (exported first), then by name.
	slices.SortFunc(index.Services, func(s1 *ServiceInfo, s2 *ServiceInfo) int {
		if s1 == s2 {
			return 0
		}

		if s1.Internal != s2.Internal {
			if s1.Internal {
				return 1
			} else {
				return -1
			}
		}

		return strings.Compare(s1.Name, s2.Name)
	})

	// Gather models.
	for _, value := range info.models.Range {
		model := value.(*ModelInfo)
		index.Models = append(index.Models, model)
		// Mark presence of exported models
		if !model.Internal {
			index.HasExportedModels = true
		}
		// Update model stats.
		if len(model.Values) > 0 {
			stats.NumEnums++
		} else {
			stats.NumModels++
		}
	}

	// Sort models by internal property (exported first), then by name.
	slices.SortFunc(index.Models, func(m1 *ModelInfo, m2 *ModelInfo) int {
		if m1 == m2 {
			return 0
		}

		if m1.Internal != m2.Internal {
			if m1.Internal {
				return 1
			} else {
				return -1
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
	return !index.HasExportedServices && !index.HasExportedModels && len(index.Package.Injections) == 0
}
