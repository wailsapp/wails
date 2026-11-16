package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/pterm/pterm"
	"github.com/wailsapp/wails/v3/internal/codesignature"
	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/windowsign"
)

func signWindowsWithPassword(ctx context.Context, options *flags.Sign, password string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if options.Certificate == "" && options.Thumbprint == "" {
		return fmt.Errorf("either --certificate or --thumbprint is required for Windows signing")
	}
	optionsCopy := *options
	if optionsCopy.Timestamp == "" {
		optionsCopy.Timestamp = "http://timestamp.digicert.com"
	}
	if options.Certificate != "" && strings.ToLower(filepath.Ext(options.Input)) != ".msi" {
		err := signWindowsBuiltin(ctx, &optionsCopy, password)
		if !errors.Is(err, codesignature.ErrUnsupportedPFXEncryption) || runtime.GOOS != "windows" {
			return err
		}
		pterm.Info.Println("Using Windows SignTool for this legacy PFX encryption profile")
	}
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Windows certificate-store/provider signing and MSI signing require a Windows host with signtool.exe; portable PE/MSIX signing requires an exported AES/3DES PFX")
	}
	return signWindowsNative(ctx, &optionsCopy, password)
}

func signWindowsBuiltin(ctx context.Context, options *flags.Sign, password string) error {
	if options.Certificate == "" {
		return fmt.Errorf("an exported --certificate is required for portable Windows signing")
	}
	data, err := readSigningFile(options.Certificate, 16<<20)
	if err != nil {
		return err
	}
	defer clear(data)
	identity, err := codesignature.DecodePFX(data, password)
	if err != nil {
		return fmt.Errorf("Windows signing identity: %w", err)
	}
	output := options.Output
	if output == "" {
		output = options.Input
	}
	signing := codesignature.Options{TimestampURL: options.Timestamp}
	switch strings.ToLower(filepath.Ext(options.Input)) {
	case ".exe", ".dll":
		err = windowsign.SignPE(ctx, options.Input, output, identity, signing)
	case ".msix", ".appx":
		err = windowsign.SignMSIX(ctx, options.Input, output, identity, signing)
	default:
		return fmt.Errorf("portable Windows signing supports PE executables/DLLs and individual MSIX/APPX packages")
	}
	if err == nil {
		pterm.Success.Printfln("Signed: %s", output)
	}
	return err
}

func signWindowsNative(ctx context.Context, options *flags.Sign, password string) error {
	signtool, err := findSigntool()
	if err != nil {
		return err
	}
	args := []string{"sign"}
	if options.Certificate != "" {
		args = append(args, "/f", options.Certificate)
		if password != "" {
			args = append(args, "/p", password)
		}
	} else if options.Thumbprint != "" {
		args = append(args, "/sha1", options.Thumbprint)
	} else {
		return fmt.Errorf("either --certificate or --thumbprint is required")
	}
	timestamp := options.Timestamp
	if timestamp == "" {
		timestamp = "http://timestamp.digicert.com"
	}
	args = append(args, "/tr", timestamp, "/td", "SHA256", "/fd", "SHA256")
	output := options.Output
	if output == "" {
		output = options.Input
	}
	err = produceNamedPathTransactional(output, ".windows-sign-*", filepath.Base(options.Input), func(stage string) error {
		if err := copyManifestPath(options.Input, stage); err != nil {
			return err
		}
		command := exec.CommandContext(ctx, signtool, append(args, stage)...)
		detail, err := command.CombinedOutput()
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			message := strings.TrimSpace(string(detail))
			if password != "" {
				message = strings.ReplaceAll(message, password, "<redacted>")
			}
			return fmt.Errorf("signtool failed: %s: %w", message, err)
		}
		return ctx.Err()
	})
	if err == nil {
		pterm.Success.Printfln("Signed: %s", output)
	}
	return err
}

var manifestWindowsSign = func(ctx context.Context, options *flags.Sign) error {
	return signWindowsWithPassword(ctx, options, options.Password)
}

func windowsSigningPassword(credential string, environment []string) (string, error) {
	if credential == "" {
		return "", nil
	}
	password := environmentValues(mergeEnvironment(os.Environ(), environment))[credential]
	if password == "" {
		return "", fmt.Errorf("windows signing password environment variable %s is empty", credential)
	}
	return password, nil
}

func findSigntool() (string, error) {
	if path, err := exec.LookPath("signtool.exe"); err == nil {
		return path, nil
	}
	root := os.Getenv("ProgramFiles(x86)")
	if root == "" {
		root = `C:\Program Files (x86)`
	}
	arch := map[string]string{"amd64": "x64", "arm64": "arm64", "386": "x86"}[runtime.GOARCH]
	if arch == "" {
		return "", fmt.Errorf("signtool.exe not found")
	}
	return findSDKSigningTool(filepath.Join(root, "Windows Kits", "10", "bin"), arch)
}

func findSDKSigningTool(root, arch string) (string, error) {
	candidates, _ := filepath.Glob(filepath.Join(root, "10.*", arch, "signtool.exe"))
	slices.Reverse(candidates)
	candidates = append(candidates, filepath.Join(root, arch, "signtool.exe"))
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.Mode().IsRegular() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("signtool.exe not found; install the Windows SDK and put its signing tools on PATH")
}
