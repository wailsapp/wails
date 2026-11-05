package packagesign

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// No verifier is used by the signing implementation. Run this acceptance test
// in an environment with a writable RPM database and both official readers.
func TestIndependentPackageReaders(t *testing.T) {
	if os.Getenv("WAILS_TEST_PACKAGE_SIGNING") != "1" {
		t.Skip("set WAILS_TEST_PACKAGE_SIGNING=1 for independent package verification")
	}
	for _, format := range []string{"deb", "rpm", "archlinux"} {
		t.Run(format, func(t *testing.T) {
			tool := "rpm"
			if format == "deb" {
				tool = "dpkg-sig"
			} else if format == "archlinux" {
				tool = "gpg"
			}
			for _, name := range []string{"gpg", tool} {
				if _, err := exec.LookPath(name); err != nil {
					t.Fatalf("required verifier %s: %v", name, err)
				}
			}
			root := t.TempDir()
			command := func(name string, args ...string) string {
				t.Helper()
				cmd := exec.CommandContext(t.Context(), name, args...)
				cmd.Env = append(os.Environ(), "GNUPGHOME="+root)
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("%s: %v: %s", name, err, output)
				}
				return string(output)
			}
			out := filepath.Join(root, "signed."+format)
			input := "testdata/unsigned." + format
			if format == "archlinux" {
				input = "testdata/unsigned.deb"
			}
			if err := Sign(t.Context(), input, out, fixtureOptions(format)); err != nil {
				t.Fatal(err)
			}
			if format == "archlinux" {
				command("gpg", "--batch", "--import", "testdata/public.asc")
				command("gpg", "--batch", "--verify", out, input)
			} else if format == "deb" {
				command("gpg", "--batch", "--import", "testdata/public.asc")
				if result := command(tool, "--verify", out); !strings.Contains(result, "GOODSIG") {
					t.Fatalf("missing independent signature verification: %s", result)
				}
			} else {
				db := filepath.Join(root, "rpmdb")
				command(tool, "--dbpath", db, "--initdb")
				command(tool, "--dbpath", db, "--import", "testdata/public.asc")
				result := command(tool, "--dbpath", db, "--checksig", "--verbose", out)
				if !strings.Contains(result, "Signature") || !strings.Contains(result, "OK") {
					t.Fatalf("missing independent signature verification: %s", result)
				}
			}
		})
	}
}
