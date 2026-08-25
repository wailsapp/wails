package manifest

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// InitialState is the complete manifest intent selected while scaffolding a
// project. The writer owns how that state is represented in sparse HCL.
type InitialState struct {
	Project    Project
	TypeScript bool
	Interfaces bool
}

func Minimal(project Project) []byte {
	return EncodeInitial(InitialState{Project: project, TypeScript: true, Interfaces: true})
}

func EncodeInitial(state InitialState) []byte {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.SetAttributeValue("version", cty.NumberIntVal(3))
	body.AppendNewline()
	projectBlock := hclwrite.NewBlock("project", nil)
	projectBody := projectBlock.Body()
	projectBody.SetAttributeValue("name", cty.StringVal(state.Project.Name))
	projectBody.SetAttributeValue("product_name", cty.StringVal(state.Project.ProductName))
	projectBody.SetAttributeValue("identifier", cty.StringVal(state.Project.Identifier))
	projectBody.SetAttributeValue("version", cty.StringVal(state.Project.Version))
	projectBody.SetAttributeValue("binary_name", cty.StringVal(deriveBinaryName(state.Project.Name)))
	setOptionalStringAttribute(projectBody, "company", state.Project.CompanyName)
	setOptionalStringAttribute(projectBody, "description", state.Project.Description)
	setOptionalStringAttribute(projectBody, "copyright", state.Project.Copyright)
	setOptionalStringAttribute(projectBody, "comments", state.Project.Comments)
	body.AppendBlock(projectBlock)
	body.AppendNewline()
	frontendBlock := hclwrite.NewBlock("frontend", nil)
	frontendBody := frontendBlock.Body()
	frontendBody.SetAttributeValue("directory", cty.StringVal("frontend"))
	frontendBody.SetAttributeValue("install", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("install")}))
	frontendBody.SetAttributeValue("build", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("run"), cty.StringVal("build")}))
	frontendBody.SetAttributeValue("dev", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("run"), cty.StringVal("dev")}))
	frontendBody.SetAttributeValue("output", cty.StringVal("frontend/dist"))
	frontendBody.AppendNewline()
	bindingsBlock := hclwrite.NewBlock("bindings", nil)
	bindingsBlock.Body().SetAttributeValue("typescript", cty.BoolVal(state.TypeScript))
	bindingsBlock.Body().SetAttributeValue("interfaces", cty.BoolVal(state.TypeScript && state.Interfaces))
	frontendBody.AppendBlock(bindingsBlock)
	body.AppendBlock(frontendBlock)
	body.AppendNewline()
	buildBlock := hclwrite.NewBlock("build", nil)
	buildBlock.Body().SetAttributeValue("output", cty.StringVal("bin"))
	body.AppendBlock(buildBlock)
	body.AppendNewline()
	return file.Bytes()
}

func setOptionalStringAttribute(body *hclwrite.Body, name, value string) {
	if value != "" {
		body.SetAttributeValue(name, cty.StringVal(value))
	}
}

func WriteMinimal(root string, project Project) error {
	if err := validateProject(project); err != nil {
		return err
	}
	return atomicWrite(filepath.Join(root, Filename), Minimal(project), 0o644)
}

func WriteInitial(root string, state InitialState) error {
	if err := validateProject(state.Project); err != nil {
		return err
	}
	return atomicWrite(filepath.Join(root, Filename), EncodeInitial(state), 0o644)
}

// UpdateProjectMetadata replaces only the project values owned by project
// setup. All other Manifest intent, comments, ordering, and file permissions
// remain user-owned and are preserved.
func UpdateProjectMetadata(start string, project Project) error {
	return updateInitialState(start, project, nil)
}

// UpdateInitialState applies scaffold-owned intent to an existing manifest
// while preserving all template-owned configuration, comments, ordering, and
// file permissions.
func UpdateInitialState(start string, state InitialState) error {
	return updateInitialState(start, state.Project, &state)
}

func updateInitialState(start string, project Project, initial *InitialState) error {
	if err := validateProject(project); err != nil {
		return err
	}
	return updateManifest(start, func(file *hclwrite.File) error {
		var projectBody *hclwrite.Body
		for _, block := range file.Body().Blocks() {
			if block.Type() == "project" {
				projectBody = block.Body()
				break
			}
		}
		if projectBody == nil {
			return fmt.Errorf("%s: project block is required", Filename)
		}
		projectBody.SetAttributeValue("name", cty.StringVal(project.Name))
		projectBody.SetAttributeValue("product_name", cty.StringVal(project.ProductName))
		projectBody.SetAttributeValue("identifier", cty.StringVal(project.Identifier))
		projectBody.SetAttributeValue("version", cty.StringVal(project.Version))
		updateOptionalStringAttribute(projectBody, "company", project.CompanyName)
		updateOptionalStringAttribute(projectBody, "description", project.Description)
		updateOptionalStringAttribute(projectBody, "copyright", project.Copyright)
		updateOptionalStringAttribute(projectBody, "comments", project.Comments)
		if initial != nil {
			frontendBody := ensureBlock(file.Body(), "frontend").Body()
			bindingsBody := ensureBlock(frontendBody, "bindings").Body()
			bindingsBody.SetAttributeValue("typescript", cty.BoolVal(initial.TypeScript))
			bindingsBody.SetAttributeValue("interfaces", cty.BoolVal(initial.TypeScript && initial.Interfaces))
		}
		return nil
	})
}

// UpdateSigningPlatform replaces the project signing intent for one platform
// without disturbing any other manifest content.
func UpdateSigningPlatform(start, platform string, signing SigningPlatform) error {
	if !contains([]string{"windows", "darwin", "linux", "ios", "android"}, platform) {
		return fmt.Errorf("unsupported signing platform %q", platform)
	}
	return updateManifest(start, func(file *hclwrite.File) error {
		platformBody := ensureBlock(file.Body(), platform).Body()
		if signing.Enabled || signingHasValues(signing) {
			signingBody := ensureBlock(platformBody, "signing").Body()
			updateOptionalStringAttribute(signingBody, "credential", signing.Credential)
			updateOptionalStringAttribute(signingBody, "identity", signing.Identity)
			updateOptionalStringAttribute(signingBody, "certificate", signing.Certificate)
			updateOptionalStringAttribute(signingBody, "thumbprint", signing.Thumbprint)
			updateOptionalStringAttribute(signingBody, "timestamp_server", signing.TimestampServer)
			updateOptionalStringAttribute(signingBody, "entitlements", signing.Entitlements)
			updateOptionalStringAttribute(signingBody, "provisioning_profile", signing.ProvisioningProfile)
			updateOptionalStringAttribute(signingBody, "key_alias", signing.KeyAlias)
		} else {
			removeBlocks(platformBody, "signing")
		}
		if signing.Notarize {
			notarizationBody := ensureBlock(platformBody, "notarization").Body()
			updateOptionalStringAttribute(notarizationBody, "credential", signing.NotarizationCredential)
		} else {
			removeBlocks(platformBody, "notarization")
		}
		return nil
	})
}

func signingHasValues(signing SigningPlatform) bool {
	return signing.Identity != "" || signing.Certificate != "" || signing.Thumbprint != "" || signing.TimestampServer != "" || signing.Entitlements != "" || signing.ProvisioningProfile != "" || signing.KeyAlias != "" || signing.Credential != ""
}

func updateManifest(start string, mutate func(*hclwrite.File) error) error {
	loaded, err := Load(start, "")
	if err != nil {
		return err
	}
	file, diagnostics := hclwrite.ParseConfig(loaded.Raw, loaded.Path, hcl.InitialPos)
	if diagnostics.HasErrors() {
		return validationFromDiagnostics(diagnostics)
	}
	if err := mutate(file); err != nil {
		return err
	}
	data := file.Bytes()
	if _, err := decodeHCL(loaded.Config.Root, loaded.Path, data, ""); err != nil {
		return err
	}
	info, err := os.Stat(loaded.Path)
	if err != nil {
		return err
	}
	return atomicWrite(loaded.Path, data, info.Mode().Perm())
}

func ensureBlock(body *hclwrite.Body, blockType string) *hclwrite.Block {
	for _, block := range body.Blocks() {
		if block.Type() == blockType {
			return block
		}
	}
	if len(body.Attributes()) > 0 || len(body.Blocks()) > 0 {
		body.AppendNewline()
	}
	block := hclwrite.NewBlock(blockType, nil)
	body.AppendBlock(block)
	return block
}

func removeBlocks(body *hclwrite.Body, blockType string) {
	for _, block := range body.Blocks() {
		if block.Type() == blockType {
			body.RemoveBlock(block)
		}
	}
}

func updateOptionalStringAttribute(body *hclwrite.Body, name, value string) {
	if value == "" {
		body.RemoveAttribute(name)
		return
	}
	body.SetAttributeValue(name, cty.StringVal(value))
}

func EncodeConfig(config Config) ([]byte, error) {
	return encodeConfigHCL(config, "")
}

func EncodeDocument(doc Document) ([]byte, error) {
	if err := validateProject(doc.Project); err != nil {
		return nil, err
	}
	return EncodeConfig(configFromDocument(".", "", doc))
}

func WriteDocument(root string, doc Document) error {
	data, err := EncodeDocument(doc)
	if err != nil {
		return err
	}
	return atomicWrite(filepath.Join(root, Filename), data, 0o644)
}

// WriteMigrationDraft writes the inactive result of legacy analysis. Only
// `wails3 migrate --activate` may rename this file into the opt-in manifest.
func WriteMigrationDraft(root string, doc Document) error {
	return WriteMigrationDraftAt(root, MigratedFilename, doc, nil)
}

// WriteMigrationDraftAt exclusively creates an inactive, project-owned HCL
// proposal. Migration analysis may be rerun safely without replacing a draft
// the user has already reviewed or edited.
func WriteMigrationDraftAt(root, output string, doc Document, comments []string) error {
	return writeMigrationDraftAt(root, output, doc, comments, exclusiveWrite)
}

func writeMigrationDraftAt(root, output string, doc Document, comments []string, write func(string, []byte, os.FileMode) error) error {
	data, err := EncodeDocument(doc)
	if err != nil {
		return err
	}
	clean := filepath.ToSlash(filepath.Clean(output))
	if strings.EqualFold(clean, Filename) || strings.EqualFold(clean, EjectedFilename) || clean == "." || strings.HasPrefix(strings.ToLower(clean), ".wails/") {
		return fmt.Errorf("migration output %q must be an inactive project-owned HCL file", output)
	}
	path, err := ResolveProjectPath(root, "migration output", clean, false)
	if err != nil {
		return err
	}
	if !strings.EqualFold(filepath.Ext(path), ".hcl") {
		return fmt.Errorf("migration output %q must use the .hcl extension", output)
	}
	if len(comments) > 0 {
		var header strings.Builder
		for _, comment := range comments {
			for _, line := range strings.Split(comment, "\n") {
				header.WriteString("# ")
				header.WriteString(line)
				header.WriteByte('\n')
			}
		}
		header.WriteByte('\n')
		data = append([]byte(header.String()), data...)
	}
	if err := write(path, data, 0o644); err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("migration output %s already exists; refusing to overwrite it: %w", clean, err)
		}
		return err
	}
	return nil
}

func Eject(root, profile, cliVersion string, force bool) error {
	return ejectWithWriters(root, profile, cliVersion, force, EncodeEjectedHCL, exclusiveWrite, atomicWrite)
}

func ejectWithWriters(root, profile, cliVersion string, force bool, encode func(Config, string) ([]byte, error), exclusive, replace func(string, []byte, os.FileMode) error) error {
	if profile != "" {
		return fmt.Errorf("wails3 eject does not accept profiles; it writes the complete resolved manifest")
	}
	loaded, err := Load(root, "")
	if err != nil {
		return err
	}
	output := filepath.Join(loaded.Config.Root, EjectedFilename)
	data, err := encode(loaded.Config, cliVersion)
	if err != nil {
		return err
	}
	write := exclusive
	if force {
		write = replace
	}
	if err := write(output, data, 0o644); err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%s already exists; use --force to replace it", EjectedFilename)
		}
		return err
	}
	return nil
}

type temporaryFile interface {
	io.Writer
	Name() string
	Chmod(os.FileMode) error
	Sync() error
	Close() error
}

type writeOperations struct {
	mkdirAll   func(string, os.FileMode) error
	createTemp func(string, string) (temporaryFile, error)
	remove     func(string) error
	replace    func(string, string) error
	link       func(string, string) error
}

func osWriteOperations() writeOperations {
	return writeOperations{
		mkdirAll:   os.MkdirAll,
		createTemp: func(directory, pattern string) (temporaryFile, error) { return os.CreateTemp(directory, pattern) },
		remove:     os.Remove, replace: replaceFile, link: os.Link,
	}
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	return atomicWriteWithOperations(path, data, mode, osWriteOperations())
}

func atomicWriteWithOperations(path string, data []byte, mode os.FileMode, ops writeOperations) error {
	dir := filepath.Dir(path)
	if err := ops.mkdirAll(dir, 0o755); err != nil {
		return err
	}
	name, err := writeTemporary(dir, data, mode, ops)
	if err != nil {
		return err
	}
	defer ops.remove(name)
	return ops.replace(name, path)
}

func exclusiveWrite(path string, data []byte, mode os.FileMode) error {
	return exclusiveWriteWithOperations(path, data, mode, osWriteOperations())
}

func exclusiveWriteWithOperations(path string, data []byte, mode os.FileMode, ops writeOperations) error {
	dir := filepath.Dir(path)
	if err := ops.mkdirAll(dir, 0o755); err != nil {
		return err
	}
	name, err := writeTemporary(dir, data, mode, ops)
	if err != nil {
		return err
	}
	defer ops.remove(name)
	// Linking a complete same-directory temporary file publishes the initial
	// ejection atomically while retaining O_EXCL semantics on every platform.
	return ops.link(name, path)
}

func writeTemporary(dir string, data []byte, mode os.FileMode, ops writeOperations) (string, error) {
	tmp, err := ops.createTemp(dir, ".wails-hcl-*")
	if err != nil {
		return "", err
	}
	name := tmp.Name()
	complete := false
	defer func() {
		if !complete {
			_ = tmp.Close()
			_ = ops.remove(name)
		}
	}()
	if err := tmp.Chmod(mode); err != nil {
		return "", err
	}
	if _, err := tmp.Write(data); err != nil {
		return "", err
	}
	if err := tmp.Sync(); err != nil {
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	complete = true
	return name, nil
}
