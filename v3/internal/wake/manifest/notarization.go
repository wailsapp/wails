package manifest

func ValidateNotarization(signing SigningPlatform) error {
	if signing.NotarizationKey == "" && signing.NotarizationKeyID == "" && signing.NotarizationIssuerID == "" {
		return nil
	}
	for _, field := range []struct{ name, value string }{
		{"api_key", signing.NotarizationKey}, {"key_id", signing.NotarizationKeyID}, {"issuer_id", signing.NotarizationIssuerID},
	} {
		if field.value == "" {
			return fieldValidationError("darwin.notarization."+field.name, "is required with App Store Connect API credentials")
		}
	}
	if signing.NotarizationCredential != "" {
		return fieldValidationError("darwin.notarization.credential", "choose either a keychain profile or API credentials")
	}
	return nil
}
