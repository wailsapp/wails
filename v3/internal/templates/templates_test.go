package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

// --- parseTemplate ---

func TestParseTemplate_YAML_Valid(t *testing.T) {
	fsys := fstest.MapFS{
		"template.yaml": &fstest.MapFile{
			Data: []byte("name: Test\nshortname: test\nauthor: Me\ndescription: A test\nhelpurl: https://example.com\nversion: v1.0.0\nwailsVersion: 3\n"),
		},
	}
	tmpl, err := parseTemplate(fsys, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Name != "Test" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Test")
	}
	if tmpl.WailsVersion != 3 {
		t.Errorf("WailsVersion = %d, want 3", tmpl.WailsVersion)
	}
}

func TestParseTemplate_YAML_MissingWailsVersion(t *testing.T) {
	fsys := fstest.MapFS{
		"template.yaml": &fstest.MapFile{
			Data: []byte("name: Test\nshortname: test\nauthor: Me\ndescription: A test\n"),
		},
	}
	_, err := parseTemplate(fsys, "")
	if err == nil {
		t.Fatal("expected error for missing wailsVersion, got nil")
	}
}

func TestParseTemplate_YAML_WrongWailsVersion(t *testing.T) {
	fsys := fstest.MapFS{
		"template.yaml": &fstest.MapFile{
			Data: []byte("name: Test\nwailsVersion: 2\n"),
		},
	}
	_, err := parseTemplate(fsys, "")
	if err == nil {
		t.Fatal("expected error for wailsVersion 2, got nil")
	}
}

func TestParseTemplate_JSON_Schema3_BackwardsCompat(t *testing.T) {
	fsys := fstest.MapFS{
		"template.json": &fstest.MapFile{
			Data: []byte(`{"name":"Old","shortname":"old","author":"Me","description":"Old v3","version":"v0.0.1","schema":3}`),
		},
	}
	tmpl, err := parseTemplate(fsys, "")
	if err != nil {
		t.Fatalf("unexpected error for legacy schema:3 template: %v", err)
	}
	if tmpl.Name != "Old" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Old")
	}
}

func TestParseTemplate_JSON_NoSchema_V2Error(t *testing.T) {
	fsys := fstest.MapFS{
		"template.json": &fstest.MapFile{
			Data: []byte(`{"name":"OldV2","shortname":"oldv2","author":"Me","description":"A v2 template"}`),
		},
	}
	_, err := parseTemplate(fsys, "")
	if err == nil {
		t.Fatal("expected error for v2 template (no schema), got nil")
	}
}

func TestParseTemplate_JSON_WrongSchema_Error(t *testing.T) {
	fsys := fstest.MapFS{
		"template.json": &fstest.MapFile{
			Data: []byte(`{"name":"Bad","schema":99}`),
		},
	}
	_, err := parseTemplate(fsys, "")
	if err == nil {
		t.Fatal("expected error for unsupported schema, got nil")
	}
}

func TestParseTemplate_NoFiles_Error(t *testing.T) {
	fsys := fstest.MapFS{}
	_, err := parseTemplate(fsys, "")
	if err == nil {
		t.Fatal("expected error when neither template.yaml nor template.json exist, got nil")
	}
}

// YAML takes precedence over JSON when both exist.
func TestParseTemplate_YAML_TakesPrecedenceOverJSON(t *testing.T) {
	fsys := fstest.MapFS{
		"template.yaml": &fstest.MapFile{
			Data: []byte("name: FromYAML\nwailsVersion: 3\n"),
		},
		"template.json": &fstest.MapFile{
			// This JSON has no schema — would error if parsed.
			Data: []byte(`{"name":"FromJSON"}`),
		},
	}
	tmpl, err := parseTemplate(fsys, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Name != "FromYAML" {
		t.Errorf("Name = %q, want %q (YAML should take precedence)", tmpl.Name, "FromYAML")
	}
}

// parseTemplate should work with a subdirectory name (built-in template path).
func TestParseTemplate_WithSubdirPrefix(t *testing.T) {
	fsys := fstest.MapFS{
		"mytemplate/template.yaml": &fstest.MapFile{
			Data: []byte("name: Sub\nwailsVersion: 3\n"),
		},
	}
	tmpl, err := parseTemplate(fsys, "mytemplate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Name != "Sub" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "Sub")
	}
}

// --- stripUnsafe ---

func TestStripUnsafe_RemovesESC(t *testing.T) {
	input := "normal\x1b[31mred\x1b[0m text"
	got := stripUnsafe(input)
	want := "normal[31mred[0m text"
	if got != want {
		t.Errorf("stripUnsafe(%q) = %q, want %q", input, got, want)
	}
}

func TestStripUnsafe_RemovesControlChars(t *testing.T) {
	input := "abc\x00\x01\x02def"
	got := stripUnsafe(input)
	want := "abcdef"
	if got != want {
		t.Errorf("stripUnsafe(%q) = %q, want %q", input, got, want)
	}
}

func TestStripUnsafe_StripsNewlineAndTab(t *testing.T) {
	input := "line1\nline2\ttabbed"
	got := stripUnsafe(input)
	want := "line1line2tabbed"
	if got != want {
		t.Errorf("stripUnsafe(%q) = %q, want %q", input, got, want)
	}
}

func TestStripUnsafe_RemovesDEL(t *testing.T) {
	input := "abc\x7fdef"
	got := stripUnsafe(input)
	want := "abcdef"
	if got != want {
		t.Errorf("stripUnsafe(%q) = %q, want %q", input, got, want)
	}
}

func TestStripUnsafe_CleanString_Unchanged(t *testing.T) {
	input := "A perfectly normal template name"
	got := stripUnsafe(input)
	if got != input {
		t.Errorf("stripUnsafe modified clean string: got %q", got)
	}
}

// --- GenerateTemplate ---

func TestGenerateTemplate_CreatesExpectedFiles(t *testing.T) {
	dir := t.TempDir()
	opts := &BaseTemplate{
		Name:    "MyTemplate",
		Author:  "Test Author",
		Version: "v1.0.0",
		Dir:     dir,
	}

	if err := GenerateTemplate(opts); err != nil {
		t.Fatalf("GenerateTemplate error: %v", err)
	}

	outDir := filepath.Join(dir, "MyTemplate")

	mustExist := []string{
		"template.yaml",
		"NEXTSTEPS.md",
		"main.go.tmpl",
		"go.mod.tmpl",
		"greetservice.go",
		"Taskfile.tmpl.yml",
		"gitignore.tmpl",
		filepath.Join("frontend", "index.html"),
	}
	for _, f := range mustExist {
		path := filepath.Join(outDir, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", f, err)
		}
	}
}

func TestGenerateTemplate_TemplateYAMLContainsWailsVersion3(t *testing.T) {
	dir := t.TempDir()
	opts := &BaseTemplate{
		Name:    "VersionTest",
		Version: "v0.1.0",
		Dir:     dir,
	}
	if err := GenerateTemplate(opts); err != nil {
		t.Fatalf("GenerateTemplate error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "VersionTest", "template.yaml"))
	if err != nil {
		t.Fatalf("reading template.yaml: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "wailsVersion: 3") {
		t.Errorf("template.yaml does not contain 'wailsVersion: 3':\n%s", content)
	}
}

func TestGenerateTemplate_NoTemplateJSON(t *testing.T) {
	dir := t.TempDir()
	opts := &BaseTemplate{Name: "NoJSON", Dir: dir}
	if err := GenerateTemplate(opts); err != nil {
		t.Fatalf("GenerateTemplate error: %v", err)
	}

	jsonPath := filepath.Join(dir, "NoJSON", "template.json")
	if _, err := os.Stat(jsonPath); err == nil {
		t.Error("template.json should not be created by GenerateTemplate")
	}
}

func TestGenerateTemplate_ErrorsIfOutputExists(t *testing.T) {
	dir := t.TempDir()
	// Pre-create the output directory
	if err := os.Mkdir(filepath.Join(dir, "Exists"), 0755); err != nil {
		t.Fatal(err)
	}
	opts := &BaseTemplate{Name: "Exists", Dir: dir}
	if err := GenerateTemplate(opts); err == nil {
		t.Error("expected error when output directory already exists, got nil")
	}
}

func TestGenerateTemplate_GeneratedTemplateCanBeInstalled(t *testing.T) {
	genDir := t.TempDir()
	opts := &BaseTemplate{Name: "RoundTrip", Dir: genDir}
	if err := GenerateTemplate(opts); err != nil {
		t.Fatalf("GenerateTemplate error: %v", err)
	}

	// The generated template should pass parseTemplate without errors.
	tmplPath := filepath.Join(genDir, "RoundTrip")
	tmpl, err := parseTemplate(os.DirFS(tmplPath), "")
	if err != nil {
		t.Fatalf("generated template fails parseTemplate: %v", err)
	}
	if tmpl.WailsVersion != 3 {
		t.Errorf("WailsVersion = %d, want 3", tmpl.WailsVersion)
	}
}

