package staticanalysis

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type EmbedDetails struct {
	BaseDir   string
	EmbedPath string
	All       bool
}

func (e *EmbedDetails) GetFullPath() string {
	return filepath.Join(e.BaseDir, e.EmbedPath)
}

func GetEmbedDetails(sourcePath string) ([]*EmbedDetails, error) {
	// read in project files and determine which directories are used for embedding
	// return a list of directories

	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return nil, err
	}
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedCompiledGoFiles,
		Dir:  absPath,
	}, "./...")
	if err != nil {
		return nil, err
	}
	var result []*EmbedDetails
	for _, pkg := range pkgs {
		for index, file := range pkg.Syntax {
			baseDir := filepath.Dir(pkg.CompiledGoFiles[index])
			embedPaths := GetEmbedDetailsForFile(file, baseDir)
			if len(embedPaths) > 0 {
				result = append(result, embedPaths...)
			}
		}
	}
	return result, nil
}

func GetEmbedDetailsForFile(file *ast.File, baseDir string) []*EmbedDetails {
	var result []*EmbedDetails
	for _, comment := range file.Comments {
		for _, c := range comment.List {
			if strings.HasPrefix(c.Text, "//go:embed") {
				sl := strings.Split(c.Text, " ")
				if len(sl) == 1 {
					continue
				}
				// support for multiple paths in one comment
				for _, arg := range sl[1:] {
					embedPath := strings.TrimSpace(arg)
					// ignores all pattern matching characters except escape sequence
					if strings.Contains(embedPath, "*") || strings.Contains(embedPath, "?") || strings.Contains(embedPath, "[") {
						continue
					}
					if strings.HasPrefix(embedPath, "all:") {
						result = append(result, &EmbedDetails{
							EmbedPath: strings.TrimPrefix(embedPath, "all:"),
							All:       true,
							BaseDir:   baseDir,
						})
					} else {
						result = append(result, &EmbedDetails{
							EmbedPath: embedPath,
							All:       false,
							BaseDir:   baseDir,
						})
					}

				}
			}
		}
	}
	return result
}
