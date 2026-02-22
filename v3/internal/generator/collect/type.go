package collect

import (
	"go/ast"
	"go/types"
	"slices"
	"sync"
)

// TypeInfo records information about a type declaration.
//
// Read accesses to any public field are only safe
// if a call to [TypeInfo.Collect] has completed before the access,
// for example by calling it in the accessing goroutine
// or before spawning the accessing goroutine.
type TypeInfo struct {
	Name string

	// Alias is true for type aliases.
	Alias bool

	Doc  *ast.CommentGroup
	Decl *GroupInfo

	obj  *types.TypeName
	node ast.Node

	collector *Collector
	once      sync.Once
}

// newTypeInfo initialises a descriptor for the given named type object.
func newTypeInfo(collector *Collector, obj *types.TypeName) *TypeInfo {
	return &TypeInfo{
		obj:       obj,
		collector: collector,
	}
}

// Type returns the unique TypeInfo instance
// associated to the given object within a collector.
//
// Type is safe for concurrent use.
func (collector *Collector) Type(obj *types.TypeName) *TypeInfo {
	return collector.fromCache(obj).(*TypeInfo)
}

func (info *TypeInfo) Object() types.Object {
	return info.obj
}

func (info *TypeInfo) Type() types.Type {
	return info.obj.Type()
}

func (info *TypeInfo) Node() ast.Node {
	return info.Collect().node
}

// Collect gathers information about the type described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *TypeInfo) Collect() *TypeInfo {
	if info == nil {
		return nil
	}

	info.once.Do(func() {
		collector := info.collector

		info.Name = info.obj.Name()
		info.Alias = info.obj.IsAlias()

		path := collector.findDeclaration(info.obj)
		if path == nil {
			collector.logger.Warningf(
				"package %s: type %s: could not find declaration for type object",
				info.obj.Pkg().Path(),
				info.Name,
			)

			// Provide dummy group.
			info.Decl = newGroupInfo(nil).Collect()
			return
		}

		// path shape: *ast.TypeSpec, *ast.GenDecl, *ast.File
		tspec := path[0].(*ast.TypeSpec)

		// Retrieve doc comments.
		info.Doc = tspec.Doc
		if info.Doc == nil {
			info.Doc = tspec.Comment
		} else if tspec.Comment != nil {
			info.Doc = &ast.CommentGroup{
				List: slices.Concat(tspec.Doc.List, tspec.Comment.List),
			}
		}

		info.Decl = collector.fromCache(path[1]).(*GroupInfo).Collect()

		info.node = path[0]
	})

	return info
}
