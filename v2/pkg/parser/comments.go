package parser

import (
	"go/ast"
	"strings"
)

func parseComments(comments *ast.CommentGroup) []string {
	var result []string

	if comments == nil {
		return result
	}

	for _, comment := range comments.List {
		commentText := strings.TrimPrefix(comment.Text, "//")
		result = append(result, commentText)
	}

	return result
}
