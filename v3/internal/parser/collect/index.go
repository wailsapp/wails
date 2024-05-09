package collect

import (
	"slices"
	"strings"
)

// PackageIndex lists all bindings, models and unexported models
// generated from a package.
//
// When obtained through a call to [PackageInfo.Index],
// each binding and model name appears at most once.
type PackageIndex struct {
	Package *PackageInfo

	Services []*ServiceInfo
	Models   []*ModelInfo
	Internal []*ModelInfo
}

// Index computes a [PackageIndex] for the selected language from the list
// of generated services and models and regenerates cached stats.
//
// Services and models appear at most once in the returned slices,
// which are sorted by name.
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
	info.services.Range(func(key, value any) bool {
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
		return true
	})

	// Sort services by name.
	slices.SortFunc(index.Services, func(b1 *ServiceInfo, b2 *ServiceInfo) int {
		if b1 == b2 {
			return 0
		}
		return strings.Compare(b1.Name, b2.Name)
	})

	// Gather models.
	info.models.Range(func(key, value any) bool {
		model := value.(*ModelInfo)
		index.Models = append(index.Models, model)
		// Update model stats.
		if len(model.Values) > 0 {
			stats.NumEnums++
		} else {
			stats.NumModels++
		}
		return true
	})

	// Sort models by internal property (non-internal first), then by name.
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

	// Find first internal model.
	split, _ := slices.BinarySearchFunc(index.Models, struct{}{}, func(m *ModelInfo, _ struct{}) int {
		if m.Object().Exported() {
			return -1
		} else {
			return 1
		}
	})

	// Separate internal and non-internal models.
	index.Internal = index.Models[split:]
	index.Models = index.Models[:split]

	// Cache stats
	info.stats.Store(stats)

	return
}

// IsEmpty returns true if the given index
// contains no data for the selected language.
func (index *PackageIndex) IsEmpty() bool {
	return len(index.Package.Injections) == 0 && len(index.Services) == 0 && len(index.Models) == 0
}
