package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"gopkg.in/yaml.v3"
)

type legacyBuildConfiguration struct {
	Version string `yaml:"version"`
	Info    struct {
		CompanyName       string `yaml:"companyName"`
		ProductName       string `yaml:"productName"`
		ProductIdentifier string `yaml:"productIdentifier"`
		Description       string `yaml:"description"`
		Copyright         string `yaml:"copyright"`
		Comments          string `yaml:"comments"`
		Version           string `yaml:"version"`
		CFBundleIconName  string `yaml:"cfBundleIconName"`
	} `yaml:"info"`
	IOS struct {
		BundleID        string   `yaml:"bundleID"`
		DisplayName     string   `yaml:"displayName"`
		Version         string   `yaml:"version"`
		Company         string   `yaml:"company"`
		Comments        string   `yaml:"comments"`
		BackgroundModes []string `yaml:"backgroundModes"`
		MinimumVersion  string   `yaml:"minimumVersion"`
		MinVersion      string   `yaml:"minVersion"`
		MinIOSVersion   string   `yaml:"minIOSVersion"`
	} `yaml:"ios"`
	Android struct {
		ApplicationID string `yaml:"applicationId"`
		DisplayName   string `yaml:"displayName"`
		VersionCode   int    `yaml:"versionCode"`
		VersionName   string `yaml:"versionName"`
		MinimumSDK    int    `yaml:"minSdkVersion"`
		TargetSDK     int    `yaml:"targetSdkVersion"`
		Company       string `yaml:"company"`
		Comments      string `yaml:"comments"`
	} `yaml:"android"`
	DevMode struct {
		RootPath string `yaml:"root_path"`
		LogLevel string `yaml:"log_level"`
		Debounce int    `yaml:"debounce"`
		Ignore   struct {
			Directories      []string `yaml:"dir"`
			Files            []string `yaml:"file"`
			WatchedExtension []string `yaml:"watched_extension"`
			GitIgnore        *bool    `yaml:"git_ignore"`
		} `yaml:"ignore"`
		Executes []struct {
			Command string `yaml:"cmd"`
			Type    string `yaml:"type"`
		} `yaml:"executes"`
	} `yaml:"dev_mode"`
	FileAssociations []FileAssociation `yaml:"fileAssociations"`
	Protocols        []ProtocolConfig  `yaml:"protocols"`
	Other            any               `yaml:"other"`
}

func applyLegacyBuildConfiguration(root string, report *MigrationReport, doc *manifest.Document) error {
	path := filepath.Join(root, "build", "config.yml")
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	recordMigrationSource(report, root, path)
	var legacy legacyBuildConfiguration
	if err := yaml.Unmarshal(data, &legacy); err != nil {
		return fmt.Errorf("parse build/config.yml: %w", err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&legacy); err != nil {
		var fields *yaml.TypeError
		if !errors.As(err, &fields) {
			return fmt.Errorf("parse build/config.yml: %w", err)
		}
		for _, detail := range fields.Errors {
			blockMigration(report, "unknown-config-field", "build/config.yml", "", detail)
		}
	}
	if legacy.Version != "" && legacy.Version != "3" {
		blockMigration(report, "config-version", "build/config.yml", "", "unsupported legacy configuration version "+legacy.Version)
	}

	doc.Project.ProductName = first(legacy.Info.ProductName, doc.Project.ProductName)
	doc.Project.Identifier = first(legacy.Info.ProductIdentifier, doc.Project.Identifier)
	doc.Project.Version = first(legacy.Info.Version, doc.Project.Version)
	doc.Project.CompanyName = legacy.Info.CompanyName
	doc.Project.Description = legacy.Info.Description
	doc.Project.Copyright = legacy.Info.Copyright
	doc.Project.Comments = legacy.Info.Comments
	doc.Targets.Darwin.CFBundleIconName = legacy.Info.CFBundleIconName

	doc.Targets.IOS.Identifier = legacy.IOS.BundleID
	doc.Targets.IOS.ProductName = legacy.IOS.DisplayName
	doc.Targets.IOS.CompanyName = legacy.IOS.Company
	doc.Targets.IOS.Comments = legacy.IOS.Comments
	doc.Targets.IOS.BackgroundModes = append([]string(nil), legacy.IOS.BackgroundModes...)
	doc.Targets.IOS.MinimumVersion = first(legacy.IOS.MinimumVersion, first(legacy.IOS.MinVersion, legacy.IOS.MinIOSVersion))
	if legacy.IOS.Version != "" && legacy.IOS.Version != doc.Project.Version {
		blockMigration(report, "platform-version", "build/config.yml", "", fmt.Sprintf("iOS version %q differs from project version %q and has no independent manifest field", legacy.IOS.Version, doc.Project.Version))
	}

	doc.Targets.Android.Identifier = legacy.Android.ApplicationID
	doc.Targets.Android.ProductName = legacy.Android.DisplayName
	doc.Targets.Android.VersionCode = legacy.Android.VersionCode
	doc.Targets.Android.VersionName = legacy.Android.VersionName
	doc.Targets.Android.MinimumSDK = legacy.Android.MinimumSDK
	doc.Targets.Android.TargetSDK = legacy.Android.TargetSDK
	doc.Targets.Android.CompanyName = legacy.Android.Company
	doc.Targets.Android.Comments = legacy.Android.Comments

	if legacy.DevMode.RootPath != "" && filepath.Clean(legacy.DevMode.RootPath) != "." {
		blockMigration(report, "dev-root", "build/config.yml", "", "dev_mode.root_path is not representable unless it is the project root")
	}
	if legacy.DevMode.LogLevel != "" {
		doc.Dev.LogLevel = legacy.DevMode.LogLevel
	}
	if legacy.DevMode.Debounce != 0 {
		doc.Dev.DebounceMS = legacy.DevMode.Debounce
	}
	if legacy.DevMode.Ignore.GitIgnore != nil {
		doc.Dev.UseGitIgnore = *legacy.DevMode.Ignore.GitIgnore
	}
	if len(legacy.DevMode.Ignore.WatchedExtension) > 0 {
		doc.Dev.Watch = append([]string(nil), legacy.DevMode.Ignore.WatchedExtension...)
	}
	if len(legacy.DevMode.Ignore.Directories)+len(legacy.DevMode.Ignore.Files) > 0 {
		doc.Dev.Exclude = append(append([]string(nil), legacy.DevMode.Ignore.Directories...), legacy.DevMode.Ignore.Files...)
	}
	for _, execute := range legacy.DevMode.Executes {
		if !legacyDevCommand(execute.Command) {
			blockMigration(report, "unsupported-dev-command", "build/config.yml", "", "dev command requires a user-owned hook or manual migration: "+execute.Command)
		}
	}
	for _, association := range legacy.FileAssociations {
		doc.Associations = append(doc.Associations, manifest.Association{Extensions: []string{association.Ext}, Name: association.Name, Description: association.Description, Icon: association.IconName, Role: association.Role, MIMEType: association.MimeType})
	}
	for _, protocol := range legacy.Protocols {
		doc.Protocols = append(doc.Protocols, manifest.Protocol{Scheme: protocol.Scheme, Description: protocol.Description})
	}
	return nil
}

func applyLegacyFrontendConfiguration(root, taskPackageManager string, report *MigrationReport, doc *manifest.Document) error {
	frontend := filepath.Join(root, "frontend")
	info, err := os.Stat(frontend)
	if errors.Is(err, fs.ErrNotExist) {
		report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "info", Code: "frontend-missing", File: "frontend", Message: "no conventional frontend directory was found; manifest defaults were retained"})
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		blockMigration(report, "frontend-invalid", "frontend", "", "frontend path is not a directory")
		return nil
	}
	doc.Frontend.Directory = "frontend"
	doc.Frontend.OutputDirectory = "frontend/dist"
	if outputInfo, statErr := os.Stat(filepath.Join(frontend, "dist")); statErr == nil && !outputInfo.IsDir() {
		blockMigration(report, "frontend-output", "frontend/dist", "", "frontend output path is not a directory")
	} else if errors.Is(statErr, fs.ErrNotExist) {
		blockMigration(report, "frontend-output", "frontend", "", "frontend/dist was not found; confirm the frontend build output directory")
	} else if statErr != nil {
		return statErr
	}

	type packageMetadata struct {
		PackageManager string            `json:"packageManager"`
		Scripts        map[string]string `json:"scripts"`
	}
	var metadata packageMetadata
	packagePath := filepath.Join(frontend, "package.json")
	if data, readErr := os.ReadFile(packagePath); readErr == nil {
		recordMigrationSource(report, root, packagePath)
		if err := json.Unmarshal(data, &metadata); err != nil {
			blockMigration(report, "frontend-package-json", "frontend/package.json", "", "cannot parse package metadata: "+err.Error())
		} else {
			for _, script := range []string{"build", "dev"} {
				if strings.TrimSpace(metadata.Scripts[script]) == "" {
					blockMigration(report, "frontend-script", "frontend/package.json", "", "required "+script+" script is missing")
				}
			}
		}
	} else if errors.Is(readErr, fs.ErrNotExist) {
		blockMigration(report, "frontend-package-json", "frontend/package.json", "", "package.json is required to prove frontend build and dev commands")
	} else {
		return readErr
	}

	managers := map[string][]string{}
	for manager, lockfiles := range map[string][]string{
		"npm": {"package-lock.json", "npm-shrinkwrap.json"}, "pnpm": {"pnpm-lock.yaml"},
		"yarn": {"yarn.lock"}, "bun": {"bun.lock", "bun.lockb"},
	} {
		for _, lockfile := range lockfiles {
			path := filepath.Join(frontend, lockfile)
			if info, statErr := os.Stat(path); statErr == nil && !info.IsDir() {
				managers[manager] = append(managers[manager], lockfile)
				recordMigrationSource(report, root, path)
			} else if statErr != nil && !errors.Is(statErr, fs.ErrNotExist) {
				return statErr
			}
		}
	}
	declaredManager := metadata.PackageManager
	if before, _, found := strings.Cut(declaredManager, "@"); found {
		declaredManager = before
	}
	managerSet := map[string]bool{}
	for manager := range managers {
		managerSet[manager] = true
	}
	if declaredManager != "" {
		managerSet[declaredManager] = true
	}
	if taskPackageManager != "" {
		managerSet[taskPackageManager] = true
	}
	managerNames := make([]string, 0, len(managerSet))
	for manager := range managerSet {
		managerNames = append(managerNames, manager)
	}
	sort.Strings(managerNames)
	if len(managerNames) > 1 {
		blockMigration(report, "package-manager-conflict", "frontend", "", "package manager signals disagree: "+strings.Join(managerNames, ", "))
	}
	manager := taskPackageManager
	if manager == "" {
		manager = declaredManager
	}
	if manager == "" && len(managerNames) > 0 {
		manager = managerNames[0]
	}
	if manager == "" {
		manager = "npm"
	}
	if !containsString([]string{"npm", "pnpm", "yarn", "bun"}, manager) {
		blockMigration(report, "package-manager", "frontend", "", "unsupported package manager "+manager)
		return nil
	}
	doc.Frontend.PackageManager = manager
	doc.Frontend.Install = []string{manager, "install"}
	doc.Frontend.Build = []string{manager, "run", "build"}
	doc.Frontend.Dev = []string{manager, "run", "dev"}

	tsconfig := filepath.Join(frontend, "tsconfig.json")
	if info, statErr := os.Stat(tsconfig); statErr == nil && !info.IsDir() {
		doc.Frontend.Bindings.TypeScript = true
		recordMigrationSource(report, root, tsconfig)
	} else if statErr != nil && !errors.Is(statErr, fs.ErrNotExist) {
		return statErr
	}
	return nil
}

func applyConventionalLegacyAssets(root string, report *MigrationReport, doc *manifest.Document) error {
	set := func(target *string, candidates ...string) error {
		for _, candidate := range candidates {
			path := filepath.Join(root, filepath.FromSlash(candidate))
			info, err := os.Stat(path)
			if err == nil && !info.IsDir() {
				*target = candidate
				recordMigrationSource(report, root, path)
				return nil
			}
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
		return nil
	}
	for _, binding := range []struct {
		target     *string
		candidates []string
	}{
		{&doc.Project.Icon, []string{"build/appicon.png"}},
		{&doc.Targets.Windows.Icon, []string{"build/windows/icon.ico", "build/icon.ico"}},
		{&doc.Targets.Windows.Manifest, []string{"build/windows/wails.exe.manifest", "build/wails.exe.manifest"}},
		{&doc.Targets.Darwin.Icon, []string{"build/darwin/icons.icns", "build/icons.icns"}},
		{&doc.Targets.Darwin.AssetsCar, []string{"build/darwin/Assets.car"}},
		{&doc.Targets.Darwin.InfoPlist, []string{"build/darwin/Info.plist", "build/Info.plist"}},
		{&doc.Targets.IOS.Icon, []string{"build/ios/icon.png"}},
		{&doc.Targets.IOS.AssetsCar, []string{"build/ios/Assets.car"}},
		{&doc.Targets.IOS.InfoPlist, []string{"build/ios/Info.plist"}},
		{&doc.Targets.Linux.DesktopEntry, []string{"build/linux/desktop"}},
		{&doc.Signing.IOS.Entitlements, []string{"build/ios/entitlements.plist"}},
		{&doc.Package.Windows.NSIS.Template, []string{"build/windows/nsis/project.nsi", "build/nsis/project.nsi"}},
		{&doc.Package.Windows.MSIX.Manifest, []string{"build/windows/msix/app_manifest.xml"}},
	} {
		if err := set(binding.target, binding.candidates...); err != nil {
			return err
		}
	}
	if doc.Signing.IOS.Entitlements != "" {
		doc.Signing.IOS.Enabled = true
	}
	doc.Package.Linux.AppImage.Icon = doc.Project.Icon
	doc.Package.Linux.AppImage.DesktopEntry = doc.Targets.Linux.DesktopEntry
	return nil
}

func recordMigrationSource(report *MigrationReport, root, path string) {
	digest, err := digestFile(path)
	if err != nil {
		return
	}
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return
	}
	report.Sources[filepath.ToSlash(relative)] = digest
}

func blockMigration(report *MigrationReport, code, file, task, message string) {
	report.Complete = false
	report.Diagnostics = append(report.Diagnostics, MigrationDiagnostic{Severity: "warning", Code: code, File: file, Task: task, Message: message})
}
