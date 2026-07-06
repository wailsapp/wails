package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/migrate"
	"github.com/wailsapp/wails/v3/internal/templates"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/version"
)

// Migrate converts a Wails v2 project into a Wails v3 project:
//
//	wails3 migrate -d ./myv2project -o ./myv3project
//
// It parses wails.json and the declarative options.App literal passed to
// wails.Run, generates an equivalent programmatic v3 main file, scaffolds the
// v3 build system (Taskfile + build assets), migrates the frontend (rewriting
// the generated wailsjs modules onto @wailsio/runtime) and rewrites v2
// runtime imports to the v3 compatibility bridge. Everything that cannot be
// migrated automatically is recorded in MIGRATION.md.
func Migrate(options *flags.Migrate) error {
	DisableFooter = true

	if options.Quiet {
		term.DisableOutput()
	}
	term.Header("Migrate Wails v2 project")
	term.Warningf("This command is experimental. Review the generated project and test it thoroughly.\n")
	term.Warningf("Please report problems at https://github.com/wailsapp/wails/issues - PRs improving the tool are very welcome.\n")

	if options.OutputDir == "" {
		return errors.New("please use the -o flag to specify an output directory for the migrated project")
	}

	proj, err := migrate.ParseV2Project(options.V2Dir)
	if err != nil {
		return err
	}

	outDir, err := filepath.Abs(options.OutputDir)
	if err != nil {
		return err
	}
	if isSubPath(proj.Dir, outDir) {
		return fmt.Errorf("the output directory must not be inside the v2 project directory")
	}
	if entries, err := os.ReadDir(outDir); err == nil && len(entries) > 0 && !options.Force {
		return fmt.Errorf("output directory %s is not empty (use -f to write into it anyway)", outDir)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	cfg := proj.Config
	term.Infof("Migrating %s (module %s)\n", cfg.Name, proj.ModulePath)

	// Map the declarative v2 options onto the v3 API and rewrite the main file.
	v3opts := migrate.MapOptions(proj)
	mainSrc, err := migrate.GenerateMain(proj, v3opts)
	if err != nil {
		return err
	}

	// Scaffold the v3 project pieces (Taskfile, .gitignore) from the project
	// template, then generate the build assets from the v2 project metadata.
	if err := scaffoldV3(proj, outDir, options.Quiet); err != nil {
		return err
	}

	// Write the migrated Go sources.
	mainRel, err := filepath.Rel(proj.Dir, proj.Main.Path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Join(outDir, mainRel)), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outDir, mainRel), mainSrc, 0o644); err != nil {
		return err
	}
	if err := migrate.CopyProjectFiles(proj, outDir); err != nil {
		return err
	}

	// The compatibility bridge is generated into the project (not shipped as
	// part of the v3 module) so that only migrated projects carry it, and its
	// owners can delete it as they finish porting to the v3 API.
	if proj.UsesV2Runtime || v3opts.NeedsLifecycleService() {
		if err := migrate.WriteCompatBridge(proj, outDir); err != nil {
			return err
		}
	}

	// go.mod: swap wails/v2 for wails/v3, keep everything else.
	// LatestStable is the released tag even in dev builds, so the generated
	// require is always resolvable.
	goMod, err := migrate.TransformGoMod(proj, version.LatestStable())
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outDir, "go.mod"), goMod, 0o644); err != nil {
		return err
	}

	// Frontend: copy, regenerate wailsjs as a compatibility layer, add the
	// @wailsio/runtime dependency.
	if err := migrate.MigrateFrontend(proj, outDir); err != nil {
		return err
	}

	// Write the migration report.
	reportPath := filepath.Join(outDir, "MIGRATION.md")
	if err := os.WriteFile(reportPath, []byte(proj.Report.Markdown()), 0o644); err != nil {
		return err
	}

	if !options.SkipGoModTidy {
		term.Info("Running go mod tidy...")
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = outDir
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			term.Warningf("go mod tidy failed (%v) - run it manually in %s\n", err, outDir)
		}
	}

	term.Infof("Migration complete: %s\n", outDir)
	if proj.Report.HasManualSteps() {
		term.Warningf("Some options need manual attention - see %s\n", reportPath)
	} else {
		term.Infof("See %s for the migration summary.\n", reportPath)
	}
	term.Infof("Next: cd %s && wails3 dev\n", options.OutputDir)
	return nil
}

// isSubPath reports whether child is inside parent.
func isSubPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, "..") && !filepath.IsAbs(rel))
}

// scaffoldV3 creates the v3 project skeleton pieces around the migrated
// sources: the root Taskfile, a .gitignore when the project has none, and
// the build directory (config.yml, platform Taskfiles, icons, packaging
// templates) populated from the v2 project metadata.
func scaffoldV3(proj *V2ProjectAlias, outDir string, quiet bool) error {
	cfg := proj.Config

	// Render the project template into a scratch dir and lift the Taskfile
	// (and .gitignore if needed) from it. This reuses the exact same
	// scaffolding machinery as `wails3 init`.
	tmpDir, err := os.MkdirTemp("", "wails3-migrate-scaffold")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	initFlags := &flags.Init{
		ProjectName:   cfg.OutputFilename,
		TemplateName:  "vanilla",
		ProjectDir:    tmpDir,
		ModulePath:    proj.ModulePath,
		Quiet:         true,
		SkipGoModTidy: true,
	}
	// templates.Install prints a project summary (partly via bare fmt.Print),
	// scaffolds into <ProjectDir>/<ProjectName> (mutating initFlags.ProjectDir
	// to match) and changes the working directory into it. Silence it and
	// restore the working directory afterwards - the scratch project is
	// deleted.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	term.DisableOutput()
	stdout := os.Stdout
	if devnull, dnErr := os.OpenFile(os.DevNull, os.O_WRONLY, 0); dnErr == nil {
		os.Stdout = devnull
		defer devnull.Close()
	}
	err = templates.Install(initFlags)
	os.Stdout = stdout
	if !quiet {
		term.EnableOutput()
	}
	if chdirErr := os.Chdir(wd); chdirErr != nil {
		return chdirErr
	}
	if err != nil {
		return err
	}
	scaffoldDir := initFlags.ProjectDir

	taskfile, err := os.ReadFile(filepath.Join(scaffoldDir, "Taskfile.yml"))
	if err != nil {
		return err
	}
	// v2 recorded the package manager in the frontend:install command; the v3
	// Taskfile uses the PACKAGE_MANAGER variable.
	if pm := cfg.PackageManager(); pm != "npm" {
		taskfile = []byte(strings.Replace(string(taskfile), `default "npm"`, `default "`+pm+`"`, 1))
	}
	if err := os.WriteFile(filepath.Join(outDir, "Taskfile.yml"), taskfile, 0o644); err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(proj.Dir, ".gitignore")); os.IsNotExist(err) {
		gitignore, err := os.ReadFile(filepath.Join(scaffoldDir, "gitignore"))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(outDir, ".gitignore"), gitignore, 0o644); err != nil {
			return err
		}
	} else {
		proj.Report.Note("Kept the project's .gitignore; the v3 build outputs binaries to `bin/`, so consider adding `bin/` to it.")
	}

	// Build assets from v2 metadata. The product identifier keeps v2's
	// `com.wails.<name>` convention so the migrated app keeps its identity.
	company := cfg.Info.CompanyName
	if company == "" {
		company = cfg.Author.Name
	}
	copyright := ""
	if cfg.Info.Copyright != nil {
		copyright = *cfg.Info.Copyright
	}
	comments := ""
	if cfg.Info.Comments != nil {
		comments = *cfg.Info.Comments
	}
	buildAssetsOptions := &BuildAssetsOptions{
		Dir:                filepath.Join(outDir, "build"),
		Name:               cfg.Name,
		BinaryName:         cfg.OutputFilename,
		ProductName:        cfg.Info.ProductName,
		ProductVersion:     cfg.Info.ProductVersion,
		ProductIdentifier:  safeBundleID(cfg.Name),
		ProductCompany:     company,
		ProductCopyright:   copyright,
		ProductComments:    comments,
		ProductDescription: cfg.Info.ProductName,
		Silent:             true,
		Typescript:         isTypescriptFrontend(proj.FrontendDir),
		UseInterfaces:      true,
	}
	if err := GenerateBuildAssets(buildAssetsOptions); err != nil {
		return err
	}

	// Persist the values into build/config.yml so future
	// `wails3 task common:update:build-assets` runs keep them.
	configFlags := &flags.Init{
		ProjectDir:         outDir,
		ProductName:        buildAssetsOptions.ProductName,
		ProductVersion:     buildAssetsOptions.ProductVersion,
		ProductIdentifier:  buildAssetsOptions.ProductIdentifier,
		ProductCompany:     buildAssetsOptions.ProductCompany,
		ProductCopyright:   buildAssetsOptions.ProductCopyright,
		ProductComments:    buildAssetsOptions.ProductComments,
		ProductDescription: buildAssetsOptions.ProductDescription,
	}
	if err := writeProjectConfigYML(configFlags); err != nil {
		return err
	}

	// File associations and protocols move from wails.json to config.yml.
	if len(cfg.Info.FileAssociations) > 0 || len(cfg.Info.Protocols) > 0 {
		if err := appendAssociations(proj, outDir); err != nil {
			return err
		}
		updateOptions := &UpdateBuildAssetsOptions{
			Dir:                filepath.Join(outDir, "build"),
			Config:             filepath.Join(outDir, "build", "config.yml"),
			Name:               cfg.Name,
			BinaryName:         cfg.OutputFilename,
			ProductName:        buildAssetsOptions.ProductName,
			ProductVersion:     buildAssetsOptions.ProductVersion,
			ProductIdentifier:  buildAssetsOptions.ProductIdentifier,
			ProductCompany:     buildAssetsOptions.ProductCompany,
			ProductCopyright:   buildAssetsOptions.ProductCopyright,
			ProductComments:    buildAssetsOptions.ProductComments,
			ProductDescription: buildAssetsOptions.ProductDescription,
			Silent:             true,
		}
		if err := UpdateBuildAssets(updateOptions); err != nil {
			return err
		}
	}

	// Keep the application icon.
	srcIcon := filepath.Join(proj.Dir, cfg.BuildDir, "appicon.png")
	if _, err := os.Stat(srcIcon); err == nil {
		data, err := os.ReadFile(srcIcon)
		if err == nil {
			_ = os.WriteFile(filepath.Join(outDir, "build", "appicon.png"), data, 0o644)
			proj.Report.Note("Copied build/appicon.png from the v2 project. Run `wails3 task common:generate:icons` to regenerate the platform icon files (icons.icns / icon.ico) from it.")
		}
	}
	proj.Report.Note("Custom changes to the v2 build assets (Info.plist, wails.exe.manifest, NSIS scripts) were not carried over; the v3 equivalents live in `build/` and are configured through `build/config.yml`.")

	return nil
}

// V2ProjectAlias keeps the commands package decoupled from the migrate
// package's internals in function signatures.
type V2ProjectAlias = migrate.V2Project

// safeBundleID mirrors the v2 bundle-id derivation (com.wails.<name> with
// non-alphanumerics replaced) so migrated apps keep their identity.
func safeBundleID(name string) string {
	var sb strings.Builder
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	return "com.wails." + sb.String()
}

func isTypescriptFrontend(frontendDir string) bool {
	if _, err := os.Stat(filepath.Join(frontendDir, "tsconfig.json")); err == nil {
		return true
	}
	return false
}

var fileAssociationsRe = regexp.MustCompile(`(?m)^fileAssociations:\s*$`)

// appendAssociations writes the v2 file associations and protocols into the
// generated build/config.yml.
func appendAssociations(proj *V2ProjectAlias, outDir string) error {
	configPath := filepath.Join(outDir, "build", "config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	content := string(data)

	if fas := proj.Config.Info.FileAssociations; len(fas) > 0 {
		var sb strings.Builder
		sb.WriteString("fileAssociations:\n")
		for _, fa := range fas {
			sb.WriteString(fmt.Sprintf("  - ext: %s\n", fa.Ext))
			sb.WriteString(fmt.Sprintf("    name: %s\n", fa.Name))
			sb.WriteString(fmt.Sprintf("    description: %s\n", fa.Description))
			sb.WriteString(fmt.Sprintf("    iconName: %s\n", fa.IconName))
			sb.WriteString(fmt.Sprintf("    role: %s\n", fa.Role))
		}
		if fileAssociationsRe.MatchString(content) {
			content = fileAssociationsRe.ReplaceAllString(content, strings.TrimSuffix(sb.String(), "\n"))
		} else {
			content += "\n" + sb.String()
		}
		proj.Report.Note("File associations were moved to `build/config.yml`. Copy their icons into `build/` (macOS: `<iconName>.icns`, Windows: `<iconName>.ico`).")
	}

	if protos := proj.Config.Info.Protocols; len(protos) > 0 {
		var sb strings.Builder
		sb.WriteString("\nprotocols:\n")
		for _, p := range protos {
			sb.WriteString(fmt.Sprintf("  - scheme: %s\n", p.Scheme))
			if p.Description != "" {
				sb.WriteString(fmt.Sprintf("    description: %s\n", p.Description))
			}
		}
		content += sb.String()
		proj.Report.Note("Custom protocol schemes were moved to `build/config.yml`.")
	}

	return os.WriteFile(configPath, []byte(content), 0o644)
}
