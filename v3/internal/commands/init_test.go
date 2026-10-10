package commands

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/internal/defaults"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestPrepareRemoteTemplateConfig(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		wantPath   string
		wantCopied bool
	}{
		{name: "yml", filename: "config.yml", wantPath: "build/config.yml"},
		{name: "yaml", filename: "config.yaml", wantPath: "build/config.yaml", wantCopied: true},
		{name: "missing", wantPath: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir := t.TempDir()
			buildDir := filepath.Join(projectDir, "build")
			if tt.filename != "" {
				if err := os.MkdirAll(buildDir, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(buildDir, tt.filename), []byte("template config\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			gotPath, err := prepareRemoteTemplateConfig(projectDir)
			if err != nil {
				t.Fatalf("prepareRemoteTemplateConfig() error = %v", err)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("prepareRemoteTemplateConfig() path = %q, want %q", gotPath, tt.wantPath)
			}

			canonicalPath := filepath.Join(buildDir, "config.yml")
			_, canonicalErr := os.Stat(canonicalPath)
			if tt.wantCopied {
				if canonicalErr != nil {
					t.Fatalf("expected config.yaml to be copied to config.yml: %v", canonicalErr)
				}
				got, err := os.ReadFile(canonicalPath)
				if err != nil {
					t.Fatal(err)
				}
				if string(got) != "template config\n" {
					t.Errorf("copied config = %q, want original template config", got)
				}
			} else if tt.filename == "" && !os.IsNotExist(canonicalErr) {
				t.Errorf("unexpected canonical config file for missing source")
			}
		})
	}
}

func TestPrepareRemoteTemplateConfigRejectsSymlink(t *testing.T) {
	projectDir := t.TempDir()
	buildDir := filepath.Join(projectDir, "build")
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(projectDir, "outside-config.yml")
	if err := os.WriteFile(target, []byte("not for copying\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(buildDir, "config.yaml")); err != nil {
		t.Skipf("symbolic links are not supported: %v", err)
	}

	if _, err := prepareRemoteTemplateConfig(projectDir); err == nil {
		t.Fatal("expected symbolic-link configuration to be rejected")
	}
}

func TestGitURLToModulePath(t *testing.T) {
	tests := []struct {
		name   string
		gitURL string
		want   string
	}{
		{
			name:   "Simple GitHub URL",
			gitURL: "github.com/username/project",
			want:   "github.com/username/project",
		},
		{
			name:   "GitHub URL with .git suffix",
			gitURL: "github.com/username/project.git",
			want:   "github.com/username/project",
		},
		{
			name:   "HTTPS GitHub URL",
			gitURL: "https://github.com/username/project",
			want:   "github.com/username/project",
		},
		{
			name:   "HTTPS GitHub URL with .git suffix",
			gitURL: "https://github.com/username/project.git",
			want:   "github.com/username/project",
		},
		{
			name:   "HTTP GitHub URL",
			gitURL: "http://github.com/username/project",
			want:   "github.com/username/project",
		},
		{
			name:   "HTTP GitHub URL with .git suffix",
			gitURL: "http://github.com/username/project.git",
			want:   "github.com/username/project",
		},
		{
			name:   "SSH GitHub URL",
			gitURL: "git@github.com:username/project",
			want:   "github.com/username/project",
		},
		{
			name:   "SSH GitHub URL with .git suffix",
			gitURL: "git@github.com:username/project.git",
			want:   "github.com/username/project",
		},
		{
			name:   "Alternative SSH URL format",
			gitURL: "ssh://git@github.com/username/project.git",
			want:   "github.com/username/project",
		},
		{
			name:   "Git protocol URL",
			gitURL: "git://github.com/username/project.git",
			want:   "github.com/username/project",
		},
		{
			name:   "File system URL",
			gitURL: "file:///path/to/project.git",
			want:   "path/to/project",
		},
		{
			name:   "SSH GitLab URL",
			gitURL: "git@gitlab.com:username/project.git",
			want:   "gitlab.com/username/project",
		},
		{
			name:   "SSH Custom Domain",
			gitURL: "git@git.company.com:username/project.git",
			want:   "git.company.com/username/project",
		},
		{
			name:   "GitLab URL",
			gitURL: "gitlab.com/username/project",
			want:   "gitlab.com/username/project",
		},
		{
			name:   "BitBucket URL",
			gitURL: "bitbucket.org/username/project",
			want:   "bitbucket.org/username/project",
		},
		{
			name:   "Custom domain",
			gitURL: "git.company.com/username/project",
			want:   "git.company.com/username/project",
		},
		{
			name:   "Custom domain with HTTPS and .git",
			gitURL: "https://git.company.com/username/project.git",
			want:   "git.company.com/username/project",
		},
		{
			name:   "Empty string",
			gitURL: "",
			want:   "",
		},
		{
			name:   "Just .git suffix",
			gitURL: ".git",
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gitURLToModulePath(tt.gitURL); got != tt.want {
				t.Errorf("gitURLToModulePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyGlobalDefaults(t *testing.T) {
	currentYear := time.Now().Format("2006")

	tests := []struct {
		name           string
		options        flags.Init
		globalDefaults defaults.GlobalDefaults
		wantOptions    flags.Init
	}{
		{
			name: "applies template default when using vanilla",
			options: flags.Init{
				TemplateName: "vanilla",
				ProjectName:  "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Project: defaults.ProjectDefaults{
					DefaultTemplate: "react",
				},
			},
			wantOptions: flags.Init{
				TemplateName: "react",
				ProjectName:  "myapp",
			},
		},
		{
			name: "does not override non-default template",
			options: flags.Init{
				TemplateName: "svelte",
				ProjectName:  "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Project: defaults.ProjectDefaults{
					DefaultTemplate: "react",
				},
			},
			wantOptions: flags.Init{
				TemplateName: "svelte",
				ProjectName:  "myapp",
			},
		},
		{
			name: "applies company default when using My Company",
			options: flags.Init{
				ProductCompany: "My Company",
				ProjectName:    "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Author: defaults.AuthorDefaults{
					Company: "Acme Corp",
				},
			},
			wantOptions: flags.Init{
				ProductCompany: "Acme Corp",
				ProjectName:    "myapp",
			},
		},
		{
			name: "does not override non-default company",
			options: flags.Init{
				ProductCompany: "Custom Company",
				ProjectName:    "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Author: defaults.AuthorDefaults{
					Company: "Acme Corp",
				},
			},
			wantOptions: flags.Init{
				ProductCompany: "Custom Company",
				ProjectName:    "myapp",
			},
		},
		{
			name: "applies copyright from global defaults",
			options: flags.Init{
				ProductCopyright: "\u00a9 now, My Company",
				ProjectName:      "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Author: defaults.AuthorDefaults{
					Company: "Acme Corp",
				},
				Project: defaults.ProjectDefaults{
					CopyrightTemplate: "© {year}, {company}",
				},
			},
			wantOptions: flags.Init{
				ProductCopyright: "© " + currentYear + ", Acme Corp",
				ProjectName:      "myapp",
			},
		},
		{
			name: "applies product identifier from global defaults",
			options: flags.Init{
				ProductIdentifier: "",
				ProjectName:       "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Project: defaults.ProjectDefaults{
					ProductIdentifierPrefix: "com.acme",
				},
			},
			wantOptions: flags.Init{
				ProductIdentifier: "com.acme.myapp",
				ProjectName:       "myapp",
			},
		},
		{
			name: "does not override explicitly set product identifier",
			options: flags.Init{
				ProductIdentifier: "com.custom.app",
				ProjectName:       "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Project: defaults.ProjectDefaults{
					ProductIdentifierPrefix: "com.acme",
				},
			},
			wantOptions: flags.Init{
				ProductIdentifier: "com.custom.app",
				ProjectName:       "myapp",
			},
		},
		{
			name: "applies description from global defaults",
			options: flags.Init{
				ProductDescription: "My Product Description",
				ProjectName:        "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Project: defaults.ProjectDefaults{
					DescriptionTemplate: "A {name} application",
				},
			},
			wantOptions: flags.Init{
				ProductDescription: "A myapp application",
				ProjectName:        "myapp",
			},
		},
		{
			name: "does not override non-default description",
			options: flags.Init{
				ProductDescription: "Custom description",
				ProjectName:        "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Project: defaults.ProjectDefaults{
					DescriptionTemplate: "A {name} application",
				},
			},
			wantOptions: flags.Init{
				ProductDescription: "Custom description",
				ProjectName:        "myapp",
			},
		},
		{
			name: "applies version from global defaults",
			options: flags.Init{
				ProductVersion: "0.1.0",
				ProjectName:    "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Project: defaults.ProjectDefaults{
					DefaultVersion: "1.0.0",
				},
			},
			wantOptions: flags.Init{
				ProductVersion: "1.0.0",
				ProjectName:    "myapp",
			},
		},
		{
			name: "does not override non-default version",
			options: flags.Init{
				ProductVersion: "2.0.0",
				ProjectName:    "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{
				Project: defaults.ProjectDefaults{
					DefaultVersion: "1.0.0",
				},
			},
			wantOptions: flags.Init{
				ProductVersion: "2.0.0",
				ProjectName:    "myapp",
			},
		},
		{
			name: "applies all defaults together",
			options: flags.Init{
				TemplateName:       "vanilla",
				ProductCompany:     "My Company",
				ProductCopyright:   "\u00a9 now, My Company",
				ProductIdentifier:  "",
				ProductDescription: "My Product Description",
				ProductVersion:     "0.1.0",
				ProjectName:        "testproject",
			},
			globalDefaults: defaults.GlobalDefaults{
				Author: defaults.AuthorDefaults{
					Company: "Wails Inc",
				},
				Project: defaults.ProjectDefaults{
					DefaultTemplate:         "vue",
					ProductIdentifierPrefix: "com.wails",
					CopyrightTemplate:       "© {year}, {company}",
					DescriptionTemplate:     "A {name} application",
					DefaultVersion:          "1.0.0",
				},
			},
			wantOptions: flags.Init{
				TemplateName:       "vue",
				ProductCompany:     "Wails Inc",
				ProductCopyright:   "© " + currentYear + ", Wails Inc",
				ProductIdentifier:  "com.wails.testproject",
				ProductDescription: "A testproject application",
				ProductVersion:     "1.0.0",
				ProjectName:        "testproject",
			},
		},
		{
			name: "empty global defaults normalize copyright year but leave other options unchanged",
			options: flags.Init{
				TemplateName:       "vanilla",
				ProductCompany:     "My Company",
				ProductCopyright:   "\u00a9 now, My Company",
				ProductIdentifier:  "",
				ProductDescription: "My Product Description",
				ProductVersion:     "0.1.0",
				ProjectName:        "myapp",
			},
			globalDefaults: defaults.GlobalDefaults{},
			wantOptions: flags.Init{
				TemplateName:       "vanilla",
				ProductCompany:     "My Company",
				ProductCopyright:   "© " + currentYear + ", My Company",
				ProductIdentifier:  "",
				ProductDescription: "My Product Description",
				ProductVersion:     "0.1.0",
				ProjectName:        "myapp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := tt.options
			applyGlobalDefaults(&options, tt.globalDefaults)

			if options.TemplateName != tt.wantOptions.TemplateName {
				t.Errorf("TemplateName = %q, want %q", options.TemplateName, tt.wantOptions.TemplateName)
			}
			if options.ProductCompany != tt.wantOptions.ProductCompany {
				t.Errorf("ProductCompany = %q, want %q", options.ProductCompany, tt.wantOptions.ProductCompany)
			}
			if options.ProductCopyright != tt.wantOptions.ProductCopyright {
				t.Errorf("ProductCopyright = %q, want %q", options.ProductCopyright, tt.wantOptions.ProductCopyright)
			}
			if options.ProductIdentifier != tt.wantOptions.ProductIdentifier {
				t.Errorf("ProductIdentifier = %q, want %q", options.ProductIdentifier, tt.wantOptions.ProductIdentifier)
			}
			if options.ProductDescription != tt.wantOptions.ProductDescription {
				t.Errorf("ProductDescription = %q, want %q", options.ProductDescription, tt.wantOptions.ProductDescription)
			}
			if options.ProductVersion != tt.wantOptions.ProductVersion {
				t.Errorf("ProductVersion = %q, want %q", options.ProductVersion, tt.wantOptions.ProductVersion)
			}
		})
	}
}
