package commands

import "testing"

func Test_GitURLToModuleName(t *testing.T) {
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
			if got := GitURLToModuleName(tt.gitURL); got != tt.want {
				t.Errorf("GitURLToModuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}
