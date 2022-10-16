package staticanalysis

import (
	"go/ast"
	"golang.org/x/tools/go/packages"
	"path/filepath"
	"strings"
)

type EmbedDetails struct {
	SourcePath string
	All        bool
}

func GetEmbedDetails(sourcePath string) ([]*EmbedDetails, error) {
	// read in project files and determine which directories are used for embedding
	// return a list of directories

	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return nil, err
	}
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:  absPath,
	}, "./...")
	if err != nil {
		return nil, err
	}
	var result []*EmbedDetails
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			result = append(result, GetEmbedDetailsForFile(file)...)
		}
	}
	return result, nil
}

func GetEmbedDetailsForFile(file *ast.File) []*EmbedDetails {
	var result []*EmbedDetails
	for _, comment := range file.Comments {
		for _, c := range comment.List {
			if strings.HasPrefix(c.Text, "//go:embed") {
				sl := strings.Split(c.Text, " ")
				path := ""
				all := false
				if len(sl) == 1 {
					continue
				}
				embedPath := strings.TrimSpace(sl[1])
				switch true {
				case strings.HasPrefix(embedPath, "all:"):
					path = strings.TrimPrefix(embedPath, "all:")
					all = true
				default:
					path = embedPath
				}
				result = append(result, &EmbedDetails{
					SourcePath: path,
					All:        all,
				})
			}
		}
	}
	return result
}
