package collect

import (
	"go/ast"
	"go/types"
	"sync"
)

// MethodInfo records information about a method declaration.
//
// Read accesses to any public field are only safe
// if a call to [MethodInfo.Collect] has completed before the access,
// for example by calling it in the accessing goroutine
// or before spawning the accessing goroutine.
type MethodInfo struct {
	Name string

	// Abstract is true when the described method belongs to an interface.
	Abstract bool

	Doc  *ast.CommentGroup
	Decl *GroupInfo

	obj  *types.Func
	node ast.Node

	collector *Collector
	once      sync.Once
}

func newMethodInfo(collector *Collector, obj *types.Func) *MethodInfo {
	return &MethodInfo{
		obj:       obj,
		collector: collector,
	}
}

// Method returns the unique MethodInfo instance
// associated to the given object within a collector.
//
// Method is safe for concurrent use.
func (collector *Collector) Method(obj *types.Func) *MethodInfo {
	return collector.fromCache(obj).(*MethodInfo)
}

func (info *MethodInfo) Object() types.Object {
	return info.obj
}

func (info *MethodInfo) Type() types.Type {
	return info.obj.Type()
}

func (info *MethodInfo) Node() ast.Node {
	return info.Collect().node
}

// Collect gathers information about the method described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *MethodInfo) Collect() *MethodInfo {
	if info == nil {
		return nil
	}

	info.once.Do(func() {
		collector := info.collector

		info.Name = info.obj.Name()

		path := collector.findDeclaration(info.obj)
		if path == nil {
			recv := ""
			if info.obj.Type() != nil {
				recv = info.obj.Type().(*types.Signature).Recv().Type().String() + "."
			}

			collector.logger.Warningf(
				"package %s: method %s%s: could not find declaration for method object",
				info.obj.Pkg().Path(),
				recv,
				info.obj.Name(),
			)

			// Provide dummy group.
			info.Decl = newGroupInfo(nil).Collect()
			return
		}

		// path shape: *ast.FuncDecl/*ast.Field, ...
		info.node = path[0]

		// Retrieve doc comments.
		switch n := info.node.(type) {
		case *ast.FuncDecl:
			// Concrete method.
			info.Doc = n.Doc
			info.Decl = newGroupInfo(nil).Collect() // Provide dummy group.

		case *ast.Field:
			// Abstract method.
			info.Abstract = true
			info.Decl = newGroupInfo(path[0]).Collect()
		}
	})

	return info
}
