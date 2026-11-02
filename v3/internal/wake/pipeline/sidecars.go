package pipeline

import "strings"

func publishSource(node Node) string {
	if spec, ok := node.Spec.(PackageSpec); ok && spec.Format == "appimage" {
		return spec.Output
	}
	return node.Output
}

// A package containing update files is cached as a directory. Each public file
// remains independently represented and verified in the artifact receipt.
func publishAppImageSidecars(plan *Plan) error {
	var primary []Node
	for _, key := range plan.Artifacts {
		if node := plan.Nodes[key]; node.Kind == PublishArtifact && node.Artifact.Format == "appimage" {
			primary = append(primary, node)
		}
	}
	for _, node := range primary {
		published := node.Spec.(PublishSpec)
		var config PackageSpec
		for _, key := range node.Dependencies {
			if spec, ok := plan.Nodes[key].Spec.(PackageSpec); ok && spec.Format == "appimage" {
				config = spec
				break
			}
		}
		var suffixes []string
		if config.Config.UpdateInfo != "" {
			suffixes = append(suffixes, ".zsync")
		}
		if config.Config.Checksum {
			suffixes = append(suffixes, ".sha256")
		}
		for _, suffix := range suffixes {
			key := NodeKey(string(node.Key) + suffix)
			identity := node.Artifact
			identity.Format = strings.TrimPrefix(suffix, ".")
			sidecar := Node{Key: key, Kind: PublishArtifact, Label: "Publish AppImage " + identity.Format, Scope: node.Scope, Dependencies: append([]NodeKey(nil), node.Dependencies...), Spec: PublishSpec{Source: published.Source + suffix, Destination: published.Destination + suffix}, Output: published.Destination + suffix, Cache: node.Cache, Claims: node.Claims, EstimateMS: node.EstimateMS, Artifact: identity}
			addNode(plan, sidecar, nil)
			plan.Artifacts = append(plan.Artifacts, key)
			plan.Roots = append(plan.Roots, key)
		}
	}
	return nil
}
