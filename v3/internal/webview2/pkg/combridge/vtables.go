//go:build windows

package combridge

import (
	"fmt"
	"reflect"
	"sync"

	"golang.org/x/sys/windows"
)

var (
	vTablesL sync.Mutex
	vTables  = make(map[string]*vTable)
)

// RegisterVTable registers the vtable trampoline methods for the specified ComInterface
// TBase is the base interface of T, and must be another ComInterface which roots in IUnknown or IUnknown itself.
// The first paramter of the fn is always the uintptr of the ComObject and the GoObject can be resolved with Resolve().
// After having resolved the GoObject the call must be redirected to the GoObject.
// Typically a trampoline FN looks like this.
//
//	func _ICoreWebView2NavigationCompletedEventHandlerInvoke(this uintptr, sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
//	   return combridge.Resolve[_ICoreWebView2NavigationCompletedEventHandler](this).NavigationCompleted(sender, args)
//	}
//
// The order of registration must be in the correct order as specified in the IDL of the interface.
func RegisterVTable[TParent, T IUnknown](guid string, fns ...interface{}) {
	registerVTableInternal[TParent, T](guid, false, fns...)
}

type vTable struct {
	Parent *vTable

	Name      string
	ComGUID   string
	ComVTable uintptr
	ComProcs  []uintptr
}

func registerVTableInternal[TParent, T IUnknown](guid string, isInternal bool, fns ...interface{}) {
	vTablesL.Lock()
	defer vTablesL.Unlock()

	t, tName := typeInterfaceToString[T]()
	tParent, tParentName := typeInterfaceToString[TParent]()
	if !t.Implements(tParent) {
		panic(fmt.Errorf("RegisterVTable '%s': '%s' must implement '%s'", tName, tName, tParentName))
	}

	if !isInternal {
		if t == reflect.TypeOf((*IUnknown)(nil)).Elem() {
			panic(fmt.Errorf("RegisterVTable '%s' IUnknown can't be registered", tName))
		}

		if t == tParent {
			panic(fmt.Errorf("RegisterVTable '%s': T and TParent can't be the same type", tName))
		}
	}

	var parent *vTable
	var parentProcs []uintptr
	var parentProcsCount int
	if t != tParent {
		parent = vTables[tParentName]
		if parent == nil {
			panic(fmt.Errorf("RegisterVTable '%s': Parent VTable '%s' not registered", tName, tParentName))
		}

		parentProcs = parent.ComProcs
		parentProcsCount = len(parentProcs)
	}

	comGuid, err := windows.GUIDFromString(guid)
	if err != nil {
		panic(fmt.Errorf("RegisterVTable '%s': invalid guid: %s", tName, err))
	}

	vTable := &vTable{
		Parent:  parent,
		Name:    tName,
		ComGUID: comGuid.String(),
	}
	vTable.ComVTable, vTable.ComProcs = allocUintptrObject(parentProcsCount + len(fns))

	for i, proc := range parentProcs {
		vTable.ComProcs[i] = proc
	}

	for i, fn := range fns {
		vTable.ComProcs[parentProcsCount+i] = windows.NewCallback(fn)
	}

	vTables[tName] = vTable
}

func typeInterfaceToString[T any]() (reflect.Type, string) {
	t := reflect.TypeOf((*T)(nil))
	if t.Kind() != reflect.Pointer {
		panic("must be a (*yourInterfaceType)(nil)")
	}
	t = t.Elem()
	return t, t.PkgPath() + "/" + t.Name()
}

func typeInterfaceToStringOnly[T any]() string {
	_, nane := typeInterfaceToString[T]()
	return nane
}

func guidOf[T any]() string {
	vtable := vTableOf[T]()
	if vtable == nil {
		return ""
	}
	return vtable.ComGUID
}

func vTableOf[T any]() *vTable {
	name := typeInterfaceToStringOnly[T]()
	vTablesL.Lock()
	defer vTablesL.Unlock()

	return vTables[name]
}

type ifceImpl interface {
	impl() any
	ifce() (*vTable, error)
}

type ifceDef[T any] struct {
	objImpl any
}

func (i ifceDef[T]) impl() any {
	return i.objImpl
}

func (i ifceDef[T]) ifce() (*vTable, error) {
	vtable := vTableOf[T]()
	if vtable == nil {
		return nil, fmt.Errorf("Unable to find vTable for %s", typeInterfaceToStringOnly[T]())
	}
	return vtable, nil
}
