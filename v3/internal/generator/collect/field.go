package collect

import (
	"go/ast"
	"go/token"
	"go/types"
	"sync"
)

// FieldInfo records information about a struct field declaration.
//
// Read accesses to any public field are only safe
// if a call to [FieldInfo.Collect] has completed before the access,
// for example by calling it in the accessing goroutine
// or before spawning the accessing goroutine.
type FieldInfo struct {
	Name     string
	Blank    bool
	Embedded bool

	Pos  token.Pos
	Decl *GroupInfo

	obj  *types.Var
	node ast.Node

	collector *Collector
	once      sync.Once
}

// newFieldInfo initialises a descriptor for the given field object.
func newFieldInfo(collector *Collector, obj *types.Var) *FieldInfo {
	return &FieldInfo{
		obj:       obj,
		collector: collector,
	}
}

// Field returns the unique FieldInfo instance
// associated to the given object within a collector.
//
// Field is safe for concurrent use.
func (collector *Collector) Field(obj *types.Var) *FieldInfo {
	if !obj.IsField() {
		return nil
	}

	return collector.fromCache(obj).(*FieldInfo)
}

func (info *FieldInfo) Object() types.Object {
	return info.obj
}

func (info *FieldInfo) Type() types.Type {
	return info.obj.Type()
}

func (info *FieldInfo) Node() ast.Node {
	return info.Collect().node
}

// Collect gathers information about the struct field
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
func (info *FieldInfo) Collect() *FieldInfo {
	if info == nil {
		return nil
	}

	info.once.Do(func() {
		collector := info.collector

		info.Name = info.obj.Name()
		info.Blank = (info.Name == "" || info.Name == "_")
		info.Embedded = info.obj.Embedded()

		info.Pos = info.obj.Pos()

		path := collector.findDeclaration(info.obj)
		if path == nil {
			// Do not report failure: it is expected for anonymous struct fields.
			// Provide dummy group.
			info.Decl = newGroupInfo(nil).Collect()
			return
		}

		// path shape: *ast.Field, *ast.FieldList, ...
		info.Decl = collector.fromCache(path[0]).(*GroupInfo).Collect()
		info.node = path[0]
	})

	return info
}
