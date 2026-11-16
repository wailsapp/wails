package pipeline

import (
	"fmt"

	"github.com/wailsapp/wails/v3/internal/wake/buildinfo"
)

type targetCapability = buildinfo.TargetCapability
type formatCapability = buildinfo.FormatCapability

const (
	runnableNone   = buildinfo.RunnableNone
	runnableBinary = buildinfo.RunnableBinary
	runnableApp    = buildinfo.RunnableApp
)

func lookupTarget(platform, arch string) (targetCapability, bool) {
	return buildinfo.LookupTarget(platform, arch)
}

func lookupFormat(name string) (formatCapability, bool) { return buildinfo.LookupFormat(name) }

func supportedTargetNames() []string { return buildinfo.SupportedTargetNames() }

func hostBit(host string) buildinfo.HostMask { return buildinfo.HostMaskFor(host) }

func toolchainBit(name string) buildinfo.ToolchainMask { return buildinfo.ToolchainMaskFor(name) }

func unsupportedTargetError(platform, arch string) error {
	return fmt.Errorf("unsupported target %q; supported targets: %v", platform+"/"+arch, supportedTargetNames())
}
