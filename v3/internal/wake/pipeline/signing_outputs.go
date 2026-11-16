package pipeline

func signingOutputs(input, destination string, identity ArtifactIdentity) (string, string, ArtifactIdentity) {
	if identity.Target.OS == "linux" && identity.Format == "archlinux" {
		identity.Kind, identity.Format, identity.Notarized = ArtifactMetadata, "sig", false
		return input + ".sig", destination, identity
	}
	return input + ".signed", signedPublicationPath(destination, identity.Target.OS), identity
}

func publicationSpec(node Node, destination string) (PublishSpec, ArtifactIdentity) {
	spec := PublishSpec{Source: publishSource(node), Destination: destination}
	identity := node.Artifact
	if signing, ok := node.Spec.(SignSpec); ok && signing.TargetOS == "linux" && signing.Format == "archlinux" {
		spec.Source = signing.Input
		spec.Associated = []PublishedFile{{Source: node.Output, Destination: destination + ".sig", Identity: node.Artifact}}
		identity.Kind, identity.Format = ArtifactPackage, "archlinux"
	}
	return spec, identity
}

func publicationArtifacts(node Node) []ArtifactReference {
	references := []ArtifactReference{{Key: node.Key, Path: node.Output, Identity: node.Artifact}}
	if spec, ok := node.Spec.(PublishSpec); ok {
		for _, file := range spec.Associated {
			references = append(references, ArtifactReference{Key: node.Key, Path: file.Destination, Identity: file.Identity})
		}
	}
	return references
}
