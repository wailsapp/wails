//go:build windows

package combridge

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	comIfcePointersL sync.RWMutex
	comIfcePointers  = map[uintptr]*comObject{} // Map from ComInterfacePointer to the Go ComObject
)

// Resolve the GoInterface of the specified ComInterfacePointer
func Resolve[T IUnknown](ifceP uintptr) T {
	comIfcePointersL.RLock()
	comObj := comIfcePointers[ifceP]
	comIfcePointersL.RUnlock()

	var n T
	if comObj != nil {
		t := comObj.resolve(ifceP)
		if t != nil {
			n = t.(T)
		}
	}

	return n
}

// New returns a new ComObject which implements the specified Com Interface, com calls will be redirected
// to the specified go interface.
func New[T IUnknown](obj T) *ComObject[T] {
	cObj := new(
		ifceDef[T]{obj},
	)
	return newComObject[T](cObj)
}

// New2 returns a new ComObject which implements the two specified Com Interfaces, com calls will be redirected
// to those interfaces accordingly.
// This is needed if a ComObject should implement two interfaces that are not descendants of each other,
// then you get multiple inheritance.
func New2[T IUnknown, T2 IUnknown](obj T, obj2 T2) *ComObject[T] {
	cObj := new(
		ifceDef[T]{obj},
		ifceDef[T2]{obj2},
	)
	return newComObject[T](cObj)
}

// new returns a new ComObject which implements multiple specified Com Interfaces, com calls will be redirected
// to the specified go interfaces accordingly.
// This is needed if a ComObject should implement multiple interfaces that are not descendants of each other,
// then you get multiple inheritance.
func new(impls ...ifceImpl) *comObject {
	impls = append([]ifceImpl{ifceDef[IUnknown]{}}, impls...)

	cObj := &comObject{
		refCount:  1,
		ifces:     map[string]int{},
		ifcesImpl: make([]comInterfaceDesc, len(impls)),
	}

	for i, ifceDef := range impls {
		vtable, err := ifceDef.ifce()
		if err != nil {
			panic(err)
		}

		needsImplement := false
		for table := vtable; table != nil; table = table.Parent {
			guid := table.ComGUID
			if i, found := cObj.ifces[guid]; found {
				// This Interface is already implemented
				if guid == iUnknownGUID {
					// IUnknown is a special interface and never has an user specific implementation
				} else if cObj.ifcesImpl[i].impl != ifceDef.impl() {
					panic(fmt.Sprintf("Interface '%s' is already implemented by another object", table.Name))
				}

				break
			}

			needsImplement = true
			cObj.ifces[guid] = i
		}

		if !needsImplement {
			continue
		}

		ifceP, ifcePSlice := allocUintptrObject(1)
		ifcePSlice[0] = vtable.ComVTable
		cObj.ifcesImpl[i] = comInterfaceDesc{ifceP, ifceDef.impl()}
	}

	comIfcePointersL.Lock()
	for _, ifceImpl := range cObj.ifcesImpl {
		comIfcePointers[ifceImpl.ref] = cObj
	}
	comIfcePointersL.Unlock()

	return cObj
}

func newComObject[T IUnknown](comObj *comObject) *ComObject[T] {
	c := &ComObject[T]{obj: comObj}
	// Make sure to async release since release needs locks and might block the finalizer goroutine for a longer period
	runtime.SetFinalizer(c, func(obj *ComObject[T]) { obj.close(true) })
	return c
}

// ComObject describes an exported go instance to be used as a ComObject which implements
// the specified Interface.
type ComObject[T IUnknown] struct {
	obj    *comObject
	closed int32
}

// Ref returns the native uintptr that points to the ComObject that is an interface pointer to T.
// This can be used in native calls. If the object has been closed this function will panic.
func (o *ComObject[T]) Ref() uintptr {
	if atomic.LoadInt32(&o.closed) != 0 {
		panic("ComObject has been released")
	}
	return o.obj.queryInterface(guidOf[T](), false)
}

// Close releases the native com object from the go side. It will only be destroyed if the ref counter
// reaches zero.
// After closing `Ref()` will panic.
func (o *ComObject[T]) Close() error {
	o.close(false)
	return nil
}

// close releases the native com object from the go side. It will only be destroyed if the ref counter
// reaches zero.
// After closing `Ref()` will panic.
func (o *ComObject[T]) close(asyncRelease bool) {
	if atomic.CompareAndSwapInt32(&o.closed, 0, 1) {
		runtime.SetFinalizer(o, nil)
		if asyncRelease {
			go o.obj.release()
		} else {
			o.obj.release()
		}
	}
}

type comInterfaceDesc struct {
	ref  uintptr // The native Com InterfacePointer
	impl any     // The golang target object
}

type comObject struct {
	l sync.Mutex

	refCount  int32
	ifces     map[string]int     // Map of ComInterfaceGUID to Interface Slots
	ifcesImpl []comInterfaceDesc // Slots with InterfaceDescriptors
}

func (c *comObject) queryInterface(ifceGUID string, withAddRef bool) uintptr {
	c.l.Lock()
	defer c.l.Unlock()
	if c.refCount <= 0 {
		panic("call on released com object")
	}

	i, found := c.ifces[ifceGUID]
	if !found {
		return 0
	}

	if withAddRef {
		c.refCount++
	}
	return c.ifcesImpl[i].ref
}

func (c *comObject) resolve(ifceP uintptr) any {
	c.l.Lock()
	defer c.l.Unlock()
	if c.refCount <= 0 {
		panic("call on destroyed com object")
	}

	for _, ifce := range c.ifcesImpl {
		if ifce.ref != ifceP {
			continue
		}

		return ifce.impl
	}
	return nil
}

func (c *comObject) addRef() int32 {
	c.l.Lock()
	defer c.l.Unlock()
	if c.refCount <= 0 {
		panic("call on destroyed com object")
	}

	c.refCount++
	return c.refCount
}

func (c *comObject) release() int32 {
	c.l.Lock()
	defer c.l.Unlock()
	if c.refCount <= 0 {
		panic("call on destroyed com object")
	}

	if c.refCount--; c.refCount == 0 {
		comIfcePointersL.Lock()
		for _, ref := range c.ifcesImpl {
			delete(comIfcePointers, ref.ref)
		}
		comIfcePointersL.Unlock()

		for _, impl := range c.ifcesImpl {
			ref := impl.ref
			if ref == 0 {
				continue
			}

			globalFree(ref)
		}
	}

	return c.refCount
}
