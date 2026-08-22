package manifest

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func Minimal(project Project) []byte {
	file := hclwrite.NewEmptyFile()
	body := file.Body()
	body.SetAttributeValue("version", cty.NumberIntVal(3))
	body.AppendNewline()
	projectBlock := hclwrite.NewBlock("project", nil)
	projectBody := projectBlock.Body()
	projectBody.SetAttributeValue("name", cty.StringVal(project.Name))
	projectBody.SetAttributeValue("product_name", cty.StringVal(project.ProductName))
	projectBody.SetAttributeValue("identifier", cty.StringVal(project.Identifier))
	projectBody.SetAttributeValue("version", cty.StringVal(project.Version))
	projectBody.SetAttributeValue("binary_name", cty.StringVal(deriveBinaryName(project.Name)))
	body.AppendBlock(projectBlock)
	body.AppendNewline()
	frontendBlock := hclwrite.NewBlock("frontend", nil)
	frontendBody := frontendBlock.Body()
	frontendBody.SetAttributeValue("directory", cty.StringVal("frontend"))
	frontendBody.SetAttributeValue("install", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("install")}))
	frontendBody.SetAttributeValue("build", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("run"), cty.StringVal("build")}))
	frontendBody.SetAttributeValue("dev", cty.ListVal([]cty.Value{cty.StringVal("npm"), cty.StringVal("run"), cty.StringVal("dev")}))
	frontendBody.SetAttributeValue("output", cty.StringVal("frontend/dist"))
	body.AppendBlock(frontendBlock)
	body.AppendNewline()
	buildBlock := hclwrite.NewBlock("build", nil)
	buildBlock.Body().SetAttributeValue("output", cty.StringVal("bin"))
	body.AppendBlock(buildBlock)
	body.AppendNewline()
	return file.Bytes()
}

func WriteMinimal(root string, project Project) error {
	if err := validateProject(project); err != nil {
		return err
	}
	return atomicWrite(filepath.Join(root, Filename), Minimal(project), 0o644)
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
