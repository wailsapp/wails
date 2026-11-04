package installerbuild

import (
	"bytes"
	"context"
	"debug/pe"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestInstallerCustomDefinesAreLiteralAndDeterministic(t *testing.T) {
	options := testOptions(t)
	options.Defines = map[string]string{"Z_OPTION": "yes", "APP_MESSAGE": "two words $INSTDIR\nquote\""}
	wrapper, err := wrapperScript(options.Script, options.Output, options)
	if err != nil {
		t.Fatal(err)
	}
	first := `!define APP_MESSAGE "two words $$INSTDIR$\nquote$\""`
	if !strings.Contains(wrapper, first) || strings.Index(wrapper, "APP_MESSAGE") > strings.Index(wrapper, "Z_OPTION") {
		t.Fatalf("unsafe or unstable defines:\n%s", wrapper)
	}
	err = Create(context.Background(), options, func(_ context.Context, _ string, args []string) error {
		staged := filepath.Dir(args[len(args)-1])
		text, err := os.ReadFile(args[len(args)-1])
		if err != nil {
			return err
		}
		if !strings.Contains(string(text), first) {
			t.Fatal("defines did not reach compiler")
		}
		return os.WriteFile(filepath.Join(staged, "installer.exe"), []byte("installer"), 0600)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstallerRejectsUnsafeOrOwnedDefinesBeforeCompiler(t *testing.T) {
	for _, name := range []string{"", "bad name", "bad\n!include", "0start", "PRODUCT_EXECUTABLE", "ARG_WAILS_ARM64_BINARY", "WAILS_INSTALL_SCOPE", "INFO_PRODUCTVERSION", "REQUEST_EXECUTION_LEVEL", "WAILS_SKIP_WEBVIEW2"} {
		options := testOptions(t)
		options.Defines = map[string]string{name: "value"}
		err := Create(context.Background(), options, func(context.Context, string, []string) error { t.Error("invalid define reached compiler"); return nil })
		if err == nil {
			t.Errorf("accepted define %q", name)
		}
	}
	options := testOptions(t)
	options.Defines = map[string]string{"VALID": "nul\x00value"}
	if err := options.Validate(); err == nil {
		t.Error("accepted NUL in define")
	}
}

func TestInstallerCustomDefinesWithRealCompiler(t *testing.T) {
	compiler := os.Getenv("WAILS_TEST_MAKENSIS")
	if compiler == "" {
		t.Skip("set WAILS_TEST_MAKENSIS for native compiler interoperability")
	}
	options := testOptions(t)
	options.Defines = map[string]string{"APP_MESSAGE": "two words $INSTDIR\nquote\"", "DEFINITION_SENTINEL": "custom-define-ok"}
	source, err := os.ReadFile(options.Script)
	if err != nil {
		t.Fatal(err)
	}
	source = append(source, []byte(`
!if "${DEFINITION_SENTINEL}" != "custom-define-ok"
!error "Wake did not supply the custom definition"
!endif
BrandingText "${APP_MESSAGE}"
`)...)
	if err := os.WriteFile(options.Script, source, 0600); err != nil {
		t.Fatal(err)
	}
	if err := Create(context.Background(), options, func(ctx context.Context, directory string, args []string) error {
		cmd := exec.CommandContext(ctx, compiler, args...)
		cmd.Dir = directory
		cmd.Env = CompilerEnvironment(runtime.GOOS, os.Environ())
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%w: %s", err, output)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	image, err := pe.Open(options.Output)
	if err != nil {
		t.Fatal(err)
	}
	if err := image.Close(); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(options.Script)
	if err != nil || !bytes.Equal(source, after) {
		t.Fatalf("compiler changed project script: %v", err)
	}
}
