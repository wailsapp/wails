package collect

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"
	"sync"
)

// GroupInfo records information about a group
// of type, field or constant declarations.
// This may be either a list of distinct specifications
// wrapped in parentheses, or a single specification
// declaring multiple fields or constants.
//
// Read accesses to any public field are only safe
// if a call to [GroupInfo.Collect] has completed before the access,
// for example by calling it in the accessing goroutine
// or before spawning the accessing goroutine.
type GroupInfo struct {
	Pos token.Pos
	Doc *ast.CommentGroup

	node ast.Node

	once sync.Once
}

func newGroupInfo(node ast.Node) *GroupInfo {
	return &GroupInfo{
		node: node,
	}
}

func (*GroupInfo) Object() types.Object {
	return nil
}

func (*GroupInfo) Type() types.Type {
	return nil
}

func (info *GroupInfo) Node() ast.Node {
	return info.node
}

// Collect gathers information about the declaration group
// described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *GroupInfo) Collect() *GroupInfo {
	if info == nil {
		return nil
	}

	info.once.Do(func() {
		switch n := info.node.(type) {
		case *ast.GenDecl:
			info.Pos = n.Pos()
			info.Doc = n.Doc

		case *ast.ValueSpec:
			info.Pos = n.Pos()
			info.Doc = n.Doc
			if info.Doc == nil {
				info.Doc = n.Comment
			} else if n.Comment != nil {
				info.Doc = &ast.CommentGroup{
					List: slices.Concat(n.Doc.List, n.Comment.List),
				}
			}

		case *ast.Field:
			info.Pos = n.Pos()
			info.Doc = n.Doc
			if info.Doc == nil {
				info.Doc = n.Comment
			} else if n.Comment != nil {
				info.Doc = &ast.CommentGroup{
					List: slices.Concat(n.Doc.List, n.Comment.List),
				}
			}
		}
	})

	return info
}
