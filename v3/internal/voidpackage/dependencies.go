package voidpackage

import (
	"fmt"
	"regexp"
	"slices"
)

var (
	packageNamePattern       = regexp.MustCompile(`^[A-Za-z0-9_+][A-Za-z0-9+._-]*$`)
	packageVersionPattern    = regexp.MustCompile(`^[A-Za-z0-9_+][A-Za-z0-9+._-]*-[A-Za-z0-9][A-Za-z0-9.+~]*_[0-9]+$`)
	packageComparisonPattern = regexp.MustCompile(`^[A-Za-z0-9_+][A-Za-z0-9+._-]*[<>]=?[A-Za-z0-9][A-Za-z0-9.+_~*-]*$`)
)

// Package lookup accepts bare names, but libxbps requires versioned patterns in
// dependency metadata. Exact pkgvers and existing comparisons stay unchanged.
func normalizePatterns(patterns []string) ([]string, error) {
	result := slices.Clone(patterns)
	for i, value := range result {
		switch {
		case packageVersionPattern.MatchString(value), packageComparisonPattern.MatchString(value):
		case packageNamePattern.MatchString(value):
			result[i] = value + ">=0"
		default:
			return nil, fmt.Errorf("invalid XBPS package pattern %q", value)
		}
	}
	return result, nil
}

func normalizeDependencies(metadata Metadata) (Metadata, error) {
	for _, field := range []struct {
		name     string
		patterns *[]string
	}{{"run_depends", &metadata.Depends}, {"conflicts", &metadata.Conflicts}, {"replaces", &metadata.Replaces}} {
		patterns, err := normalizePatterns(*field.patterns)
		if err != nil {
			return Metadata{}, fmt.Errorf("%s: %w", field.name, err)
		}
		*field.patterns = patterns
	}
	return metadata, nil
}
