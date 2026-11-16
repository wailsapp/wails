package commands

import (
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
)

func (h *manifestHandler) prebuiltSource(spec pipeline.PrebuiltSpec) (string, error) {
	info, err := os.Lstat(filepath.Join(h.root, filepath.FromSlash(spec.Source)))
	if err != nil {
		return "", fmt.Errorf("prebuilt binary %s: %w", spec.Source, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("prebuilt binary %s must be a regular file", spec.Source)
	}
	source, err := existingPathInsideProject(h.root, spec.Source)
	if err != nil {
		return "", fmt.Errorf("prebuilt binary %s: %w", spec.Source, err)
	}
	if err := validatePrebuiltTarget(source, spec.Target); err != nil {
		return "", fmt.Errorf("prebuilt binary %s for %s/%s: %w", spec.Source, spec.Target.OS, spec.Target.Arch, err)
	}
	return source, nil
}

func validatePrebuiltTarget(path string, target pipeline.Target) error {
	switch target.OS {
	case "linux":
		file, err := elf.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		machine := map[string]elf.Machine{"amd64": elf.EM_X86_64, "arm64": elf.EM_AARCH64}[target.Arch]
		if machine == 0 || file.Machine != machine || file.Class != elf.ELFCLASS64 || file.Data != elf.ELFDATA2LSB || file.Entry == 0 || (file.Type != elf.ET_EXEC && file.Type != elf.ET_DYN) {
			return fmt.Errorf("expected a 64-bit %s ELF executable, found %s %s", target.Arch, file.Machine, file.Type)
		}
	case "windows":
		file, err := pe.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		machine := map[string]uint16{"amd64": pe.IMAGE_FILE_MACHINE_AMD64, "arm64": pe.IMAGE_FILE_MACHINE_ARM64}[target.Arch]
		if machine == 0 || file.Machine != machine || file.Characteristics&pe.IMAGE_FILE_EXECUTABLE_IMAGE == 0 || file.Characteristics&pe.IMAGE_FILE_DLL != 0 {
			return fmt.Errorf("expected a %s PE executable, found machine %#x", target.Arch, file.Machine)
		}
	case "darwin":
		return validatePrebuiltMachO(path, target.Arch)
	default:
		return fmt.Errorf("unsupported prebuilt target %s/%s", target.OS, target.Arch)
	}
	return nil
}

func validatePrebuiltMachO(path, arch string) error {
	fat, err := macho.OpenFat(path)
	if err == nil {
		defer fat.Close()
		found := map[macho.Cpu]bool{}
		for _, component := range fat.Arches {
			if component.Type != macho.TypeExec {
				return fmt.Errorf("Mach-O component %s is not an executable", component.Cpu)
			}
			found[component.Cpu] = true
		}
		if (arch == "universal" && found[macho.CpuAmd64] && found[macho.CpuArm64]) || (arch == "amd64" && found[macho.CpuAmd64]) || (arch == "arm64" && found[macho.CpuArm64]) {
			return nil
		}
		return fmt.Errorf("Mach-O binary does not contain the requested %s architecture", arch)
	}
	file, err := macho.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	cpu := map[string]macho.Cpu{"amd64": macho.CpuAmd64, "arm64": macho.CpuArm64}[arch]
	if cpu == 0 || file.Cpu != cpu || file.Type != macho.TypeExec {
		return fmt.Errorf("expected a %s Mach-O executable, found %s %s", arch, file.Cpu, file.Type)
	}
	return nil
}

func (h *manifestHandler) importBinary(spec pipeline.PrebuiltSpec) (pipeline.RunResult, error) {
	source, err := h.prebuiltSource(spec)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	if !strings.HasPrefix(filepath.ToSlash(filepath.Clean(spec.Output)), ".wails/") {
		return pipeline.RunResult{}, fmt.Errorf("prebuilt destination must be in the generated .wails workspace")
	}
	destination, err := hookScopeOutputPath(h.root, spec.Output)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	if err := producePathTransactional(destination, ".wails-import-*", func(staged string) error {
		if err := copyManifestPath(source, staged); err != nil {
			return err
		}
		return os.Chmod(staged, 0o755)
	}); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Output: spec.Output}, nil
}
