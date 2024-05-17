package collect

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"sync"
)

// ConstInfo records information about a constant declaration.
//
// Read accesses to any public field are only safe
// if a call to [ConstInfo.Collect] has completed before the access,
// for example by calling it in the accessing goroutine
// or before spawning the accessing goroutine.
type ConstInfo struct {
	Name  string
	Value any

	Pos  token.Pos
	Spec *GroupInfo
	Decl *GroupInfo

	obj  *types.Const
	node ast.Node

	collector *Collector
	once      sync.Once
}

func newConstInfo(collector *Collector, obj *types.Const) *ConstInfo {
	return &ConstInfo{
		obj:       obj,
		collector: collector,
	}
}

// Const returns the unique ConstInfo instance
// associated to the given object within a collector.
//
// Const is safe for concurrent use.
func (collector *Collector) Const(obj *types.Const) *ConstInfo {
	return collector.fromCache(obj).(*ConstInfo)
}

func (info *ConstInfo) Object() types.Object {
	return info.obj
}

func (info *ConstInfo) Type() types.Type {
	return info.obj.Type()
}

func (info *ConstInfo) Node() ast.Node {
	return info.Collect().node
}

// Collect gathers information about the constant described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *ConstInfo) Collect() *ConstInfo {
	if info == nil {
		return nil
	}

	info.once.Do(func() {
		collector := info.collector

		info.Name = info.obj.Name()
		info.Value = constant.Val(info.obj.Val())

		info.Pos = info.obj.Pos()

		path := collector.findDeclaration(info.obj)
		if path == nil {
			collector.logger.Warningf(
				"package %s: const %s: could not find declaration for constant object",
				info.obj.Pkg().Path(),
				info.Name,
			)

			// Provide dummy groups.
			dummyGroup := newGroupInfo(nil).Collect()
			info.Spec = dummyGroup
			info.Decl = dummyGroup
			return
		}

		// path shape: *ast.ValueSpec, *ast.GenDecl, *ast.File
		info.Spec = collector.fromCache(path[0]).(*GroupInfo).Collect()
		info.Decl = collector.fromCache(path[1]).(*GroupInfo).Collect()
		info.node = path[0]
	})

	return info
}
