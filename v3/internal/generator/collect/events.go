package collect

import (
	_ "embed"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"slices"
	"strings"
	"sync"

	"golang.org/x/tools/go/ast/astutil"
)

type (
	// EventMap holds information about a set of custom events
	// and their associated data types.
	//
	// Read accesses to any public field are only safe
	// if a call to [EventMap.Collect] has completed before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	EventMap struct {
		Imports *ImportMap
		Defs    []*EventInfo

		registerEvent types.Object
		collector     *Collector
		once          sync.Once
	}

	// EventInfo holds information about a single event definition.
	EventInfo struct {
		// Name is the name of the event.
		Name string

		// Data is the data type the event has been registered with.
		// It may be nil in case of conflicting definitions.
		Data types.Type

		// Pos records the position
		// of the first discovered definition for this event.
		Pos token.Position
	}
)

func newEventMap(collector *Collector, registerEvent types.Object) *EventMap {
	return &EventMap{
		registerEvent: registerEvent,
		collector:     collector,
	}
}

// EventMap returns the unique event map associated with the given collector,
// or nil if event collection is disabled.
func (collector *Collector) EventMap() *EventMap {
	return collector.events
}

// Stats returns statistics for this event map.
// It is an error to call stats before a call to [EventMap.Collect] has completed.
func (em *EventMap) Stats() *Stats {
	return &Stats{
		NumEvents: len(em.Defs),
	}
}

// Collect gathers information for the event map described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (em *EventMap) Collect() *EventMap {
	if em == nil {
		return nil
	}

	em.once.Do(func() {
		// XXX: initialise the import map with a fake package.
		// At present this works fine; let's hope it doesn't come back and haunt us later.
		em.Imports = NewImportMap(&PackageInfo{
			Path:      em.collector.systemPaths.InternalPackage,
			collector: em.collector,
		})

		if em.registerEvent == nil {
			return
		}

		var (
			wg   sync.WaitGroup
			defs sync.Map
		)

		for pkg := range em.collector.Iterate {
			if !pkg.IsOrImportsApp {
				// Packages that are not, and do not import, the Wails application package
				// cannot define any events.
				continue
			}

			wg.Add(1)

			// Process packages in parallel.
			em.collector.scheduler.Schedule(func() {
				em.collectEventsInPackage(pkg, &defs)
				wg.Done()
			})
		}

		wg.Wait()

		// Collect valid events.
		em.Defs = slices.Collect(func(yield func(*EventInfo) bool) {
			for _, v := range defs.Range {
				event := v.(*EventInfo)
				if event.Data != nil && !IsParametric(event.Data) {
					if !yield(event) {
						break
					}
				}
			}
		})

		// Sort by name, ascending.
		slices.SortFunc(em.Defs, func(a, b *EventInfo) int {
			return strings.Compare(a.Name, b.Name)
		})

		// Record required types.
		// This must be done at the end because:
		//   - [ImportMap.AddType] does not support concurrent calls, and
		//   - we only know the set of valid events after inspecting all definitions.
		for _, def := range em.Defs {
			em.Imports.AddType(def.Data)
		}
	})

	return em
}

func (em *EventMap) collectEventsInPackage(pkg *PackageInfo, defs *sync.Map) {
	for ident, inst := range pkg.TypesInfo.Instances {
		if pkg.TypesInfo.Uses[ident] != em.registerEvent {
			continue
		}

		file := findEnclosingFile(pkg.Collect().Files, ident.Pos())
		if file == nil {
			em.collector.logger.Warningf(
				"package %s: found event declaration with no associated source file",
				pkg.Path,
			)
			continue
		}

		path, _ := astutil.PathEnclosingInterval(file, ident.Pos(), ident.End())
		if path[0] != ident {
			em.collector.logger.Warningf(
				"%v: event declaration not found in source file",
				pkg.Fset.Position(ident.Pos()),
			)
			continue
		}

		// Walk up the path: *ast.Ident -> *ast.SelectorExpr? -> (*ast.IndexExpr | *ast.IndexListExpr)? -> *ast.CallExpr?
		path = path[1:]

		if _, ok := path[0].(*ast.SelectorExpr); ok {
			path = path[1:]
		}

		if _, ok := path[0].(*ast.IndexExpr); ok {
			path = path[1:]
		} else if _, ok := path[0].(*ast.IndexListExpr); ok {
			path = path[1:]
		}

		call, ok := path[0].(*ast.CallExpr)
		if !ok {
			em.collector.logger.Warningf(
				"%v: `application.RegisterEvent` is instantiated here but not called",
				pkg.Fset.Position(path[0].Pos()),
			)
			em.collector.logger.Warningf("events registered through indirect calls are not discoverable by the binding generator: it is recommended to invoke `application.RegisterEvent` directly")
			continue
		}

		if len(call.Args) == 0 {
			// Invalid calls result in compile-time failures and can be ignored safely.
			continue
		}

		eventName, ok := pkg.TypesInfo.Types[call.Args[0]]
		if !ok || !types.AssignableTo(eventName.Type, types.Universe.Lookup("string").Type()) {
			// Mistyped calls result in compile-time failures and can be ignored safely.
			continue
		}

		if eventName.Value == nil {
			em.collector.logger.Warningf(
				"%v: `application.RegisterEvent` called here with non-constant event name",
				pkg.Fset.Position(call.Pos()),
			)
			em.collector.logger.Warningf("dynamically registered event names are not discoverable by the binding generator: it is recommended to invoke `application.RegisterEvent` with constant arguments only")
			continue
		}

		event := &EventInfo{
			Data: inst.TypeArgs.At(0),
			Pos:  pkg.Fset.Position(call.Pos()),
		}
		if eventName.Value.Kind() == constant.String {
			event.Name = constant.StringVal(eventName.Value)
		} else {
			event.Name = eventName.Value.ExactString()
		}

		if IsKnownEvent(event.Name) {
			em.collector.logger.Warningf(
				"%v: event '%s' is a known system event and cannot be overridden; this call to `application.RegisterEvent` will panic",
				event.Pos,
				event.Name,
			)
			continue
		}

		if v, ok := defs.LoadOrStore(event.Name, event); ok {
			prev := v.(*EventInfo)
			if prev.Data != nil && !types.Identical(event.Data, prev.Data) {
				next := &EventInfo{
					Name: prev.Name,
					Pos:  prev.Pos,
				}

				if defs.CompareAndSwap(prev.Name, prev, next) {
					em.collector.logger.Warningf("event '%s' has multiple conflicting definitions and will be ignored", event.Name)
					em.collector.logger.Warningf(
						"%v: event '%s' has one of multiple definitions here with data type %s",
						prev.Pos,
						prev.Name,
						prev.Data,
					)
				}

				prev = next
			}
			if prev.Data == nil {
				em.collector.logger.Warningf(
					"%v: event '%s' has one of multiple definitions here with data type %s",
					event.Pos,
					event.Name,
					event.Data,
				)
			}
			continue
		}

		// Emit unsupported type warnings only for first definition
		if IsParametric(event.Data) {
			em.collector.logger.Warningf(
				"%v: data type %s for event '%s' contains unresolved type parameters and will be ignored`",
				event.Pos,
				event.Data,
				event.Name,
			)
			em.collector.logger.Warningf("generic wrappers for calls to `application.RegisterEvent` are not analysable by the binding generator: it is recommended to call `application.RegisterEvent` with concrete types only")
		} else if types.IsInterface(event.Data) && !types.Identical(event.Data.Underlying(), typeAny) {
			em.collector.logger.Warningf(
				"%v: data type %s for event '%s' is a non-empty interface: emitting events from the frontend with data other than `null` is not supported by encoding/json and will likely result in runtime errors",
				event.Pos,
				event.Data,
				event.Name,
			)
		}
	}
}
