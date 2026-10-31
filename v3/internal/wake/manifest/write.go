package manifest

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// InitialState is the complete manifest intent selected while scaffolding a
// project. The writer owns how that state is represented in sparse YAML.
type InitialState struct {
	Project    Project
	TypeScript bool
	Interfaces bool
}

func Minimal(project Project) []byte {
	return EncodeInitial(InitialState{Project: project, TypeScript: true, Interfaces: true})
}

func EncodeInitial(state InitialState) []byte {
	return encodeInitialYAML(state)
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
	return updateManifest(start, func(root *yaml.Node) error {
		projectBody := mappingValue(root, "project")
		if projectBody == nil {
			return fmt.Errorf("%s: project section is required", Filename)
		}
		yamlUpdate(projectBody, "name", yamlString(project.Name))
		yamlUpdate(projectBody, "product_name", yamlString(project.ProductName))
		yamlUpdate(projectBody, "identifier", yamlString(project.Identifier))
		yamlUpdate(projectBody, "version", yamlString(project.Version))
		yamlUpdateOptionalString(projectBody, "company", project.CompanyName)
		yamlUpdateOptionalString(projectBody, "description", project.Description)
		yamlUpdateOptionalString(projectBody, "copyright", project.Copyright)
		yamlUpdateOptionalString(projectBody, "comments", project.Comments)
		if initial != nil {
			frontendBody := yamlEnsureMap(root, "frontend")
			bindingsBody := yamlEnsureMap(frontendBody, "bindings")
			yamlUpdate(bindingsBody, "typescript", yamlBool(initial.TypeScript))
			yamlUpdate(bindingsBody, "interfaces", yamlBool(initial.TypeScript && initial.Interfaces))
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
	return updateManifest(start, func(root *yaml.Node) error {
		platformBody := yamlEnsureMap(root, platform)
		if signing.Enabled || signingHasValues(signing) {
			signingBody := yamlEnsureMap(platformBody, "signing")
			yamlUpdateOptionalString(signingBody, "credential", signing.Credential)
			yamlUpdateOptionalString(signingBody, "identity", signing.Identity)
			yamlUpdateOptionalString(signingBody, "certificate", signing.Certificate)
			yamlUpdateOptionalString(signingBody, "thumbprint", signing.Thumbprint)
			yamlUpdateOptionalString(signingBody, "timestamp_server", signing.TimestampServer)
			yamlUpdateOptionalString(signingBody, "entitlements", signing.Entitlements)
			yamlUpdateOptionalString(signingBody, "provisioning_profile", signing.ProvisioningProfile)
			yamlUpdateOptionalString(signingBody, "key_alias", signing.KeyAlias)
		} else {
			yamlRemove(platformBody, "signing")
		}
		if signing.Notarize {
			notarizationBody := yamlEnsureMap(platformBody, "notarization")
			yamlUpdateOptionalString(notarizationBody, "credential", signing.NotarizationCredential)
		} else {
			yamlRemove(platformBody, "notarization")
		}
		return nil
	})
}

func signingHasValues(signing SigningPlatform) bool {
	return signing.Identity != "" || signing.Certificate != "" || signing.Thumbprint != "" || signing.TimestampServer != "" || signing.Entitlements != "" || signing.ProvisioningProfile != "" || signing.KeyAlias != "" || signing.Credential != ""
}

func updateManifest(start string, mutate func(*yaml.Node) error) error {
	loaded, err := Load(start, "")
	if err != nil {
		return err
	}
	var document yaml.Node
	if err := yaml.Unmarshal(loaded.Raw, &document); err != nil {
		return yamlParseError(loaded.Path, err)
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return fmt.Errorf("%s: a YAML mapping is required", Filename)
	}
	root := document.Content[0]
	if err := mutate(root); err != nil {
		return err
	}
	data, err := marshalYAML(&document)
	if err != nil {
		return err
	}
	if _, err := decodeYAML(loaded.Config.Root, loaded.Path, data, ""); err != nil {
		return err
	}
	info, err := os.Stat(loaded.Path)
	if err != nil {
		return err
	}
	return atomicWrite(loaded.Path, data, info.Mode().Perm())
}

func yamlEnsureMap(parent *yaml.Node, key string) *yaml.Node {
	if existing := mappingValue(parent, key); existing != nil && existing.Kind == yaml.MappingNode {
		return existing
	}
	result := yamlMap()
	yamlUpdate(parent, key, result)
	return result
}

func yamlUpdate(mapping *yaml.Node, key string, value *yaml.Node) {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			old := mapping.Content[index+1]
			value.HeadComment, value.LineComment, value.FootComment = old.HeadComment, old.LineComment, old.FootComment
			mapping.Content[index+1] = value
			return
		}
	}
	yamlSet(mapping, key, value)
}

func yamlRemove(mapping *yaml.Node, key string) {
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value == key {
			mapping.Content = append(mapping.Content[:index], mapping.Content[index+2:]...)
			return
		}
	}
}

func yamlUpdateOptionalString(mapping *yaml.Node, name, value string) {
	if value == "" {
		yamlRemove(mapping, name)
		return
	}
	yamlUpdate(mapping, name, yamlString(value))
}

func EncodeConfig(config Config) ([]byte, error) {
	return encodeConfigYAML(config, "")
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

// WriteMigrationDraftAt exclusively creates an inactive, project-owned YAML
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
		return fmt.Errorf("migration output %q must be an inactive project-owned YAML file", output)
	}
	path, err := ResolveProjectPath(root, "migration output", clean, false)
	if err != nil {
		return err
	}
	if extension := strings.ToLower(filepath.Ext(path)); extension != ".yaml" && extension != ".yml" {
		return fmt.Errorf("migration output %q must use the .yaml or .yml extension", output)
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
	return ejectWithWriters(root, profile, cliVersion, force, EncodeEjectedYAML, exclusiveWrite, atomicWrite)
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
	tmp, err := ops.createTemp(dir, ".wails-yaml-*")
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
