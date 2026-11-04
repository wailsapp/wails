package pipeline

import "github.com/wailsapp/wails/v3/internal/wake/manifest"

func packageSigningPolicy(config manifest.Config, request Request, format, assets string) (CachePolicy, *manifest.SigningPlatform) {
	if request.TargetOS == "ios" {
		// Even simulator assembly uses an ad-hoc signature.
		return CacheNever, nil
	}
	if format != "dmg" || !request.sign {
		return CacheArtifact, nil
	}
	for _, phase := range []manifest.HookPhase{manifest.AfterPackage, manifest.BeforeSign} {
		if _, exists := config.Hooks[phase]; exists {
			// Preserve hooks that edit the unsigned image before signing.
			return CacheArtifact, nil
		}
	}
	signing := signingPlatformForPlan(config, request.TargetOS, assets)
	return CacheNever, &signing
}
