package collect

import (
	"fmt"
	"go/ast"
	"go/types"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v3/internal/hash"
	"golang.org/x/tools/go/types/typeutil"
)

type (
	// ServiceInfo records all information that is required
	// to render JS/TS code for a service type.
	//
	// Read accesses to any public field are only safe
	// if a call to [ServiceInfo.Collect] has completed before the access,
	// for example by calling it in the accessing goroutine
	// or before spawning the accessing goroutine.
	ServiceInfo struct {
		*TypeInfo

		// Internal records whether the service
		// should be exported by the index file.
		Internal bool

		Imports *ImportMap
		Methods []*ServiceMethodInfo

		// HasInternalMethods records whether the service
		// defines lifecycle or http server methods.
		HasInternalMethods bool

		// Injections stores a list of JS code lines
		// that should be injected into the generated file.
		Injections []string

		collector *Collector
		once      sync.Once
	}

	// ServiceMethodInfo records all information that is required
	// to render JS/TS code for a service method.
	ServiceMethodInfo struct {
		*MethodInfo
		FQN      string
		ID       string
		Internal bool
		Params   []*ParamInfo
		Results  []types.Type
	}

	// ParamInfo records all information that is required
	// to render JS/TS code for a service method parameter.
	ParamInfo struct {
		Name     string
		Type     types.Type
		Blank    bool
		Variadic bool
	}
)

func newServiceInfo(collector *Collector, obj *types.TypeName) *ServiceInfo {
	return &ServiceInfo{
		TypeInfo:  collector.Type(obj),
		collector: collector,
	}
}

// Service returns the unique ServiceInfo instance
// associated to the given object within a collector
// and registers it for code generation.
//
// Service is safe for concurrent use.
func (collector *Collector) Service(obj *types.TypeName) *ServiceInfo {
	pkg := collector.Package(obj.Pkg())
	if pkg == nil {
		return nil
	}

	return pkg.recordService(obj)
}

// IsEmpty returns true if no methods or code injections
// are present for this service, for the selected language.
func (info *ServiceInfo) IsEmpty() bool {
	// Ensure information has been collected.
	info.Collect()
	return len(info.Methods) == 0 && len(info.Injections) == 0
}

// Collect gathers information about the service described by its receiver.
// It can be called concurrently by multiple goroutines;
// the computation will be performed just once.
//
// Collect returns the receiver for chaining.
// It is safe to call Collect with nil receiver.
//
// After Collect returns, the calling goroutine and all goroutines
// it might spawn afterwards are free to access
// the receiver's fields indefinitely.
func (info *ServiceInfo) Collect() *ServiceInfo {
	if info == nil {
		return nil
	}

	info.once.Do(func() {
		collector := info.collector
		obj := info.Object().(*types.TypeName)

		// Collect type information.
		info.TypeInfo.Collect()

		// Initialise import map.
		info.Imports = NewImportMap(collector.Package(obj.Pkg()))

		// Compute intuitive method set (i.e. both pointer and non-pointer receiver).
		// Do not use a method set cache because
		//   - it would hurt concurrency (requires mutual exclusion),
		//   - it is only useful when the same type is queried many times;
		//     this may only happen here if some embedded types appear frequently,
		//     which should be far from average.
		mset := typeutil.IntuitiveMethodSet(obj.Type(), nil)

		// Collect method information.
		info.Methods = make([]*ServiceMethodInfo, 0, len(mset))
		for _, sel := range mset {
			switch {
			case internalServiceMethods[sel.Obj().Name()]:
				info.HasInternalMethods = true
				continue
			case !sel.Obj().Exported():
				// Ignore unexported and internal methods.
				continue
			}

			methodInfo := info.collectMethod(sel.Obj().(*types.Func))
			if methodInfo != nil {
				info.Methods = append(info.Methods, methodInfo)
			}
		}

		// Record whether the service should be exported.
		info.Internal = !obj.Exported()

		// Parse directives.
		for _, doc := range []*ast.CommentGroup{info.Doc, info.Decl.Doc} {
			if doc == nil {
				continue
			}
			for _, comment := range doc.List {
				switch {
				case IsDirective(comment.Text, "internal"):
					info.Internal = true

				case IsDirective(comment.Text, "inject"):
					// Check condition.
					line, cond, err := ParseCondition(ParseDirective(comment.Text, "inject"))
					if err != nil {
						collector.logger.Errorf(
							"%s: in `wails:inject` directive: %v",
							collector.Package(obj.Pkg()).Fset.Position(comment.Pos()),
							err,
						)
						continue
					}

					if !cond.Satisfied(collector.options) {
						continue
					}

					// Record injected line.
					info.Injections = append(info.Injections, line)
				}
			}
		}
	})

	return info
}

// internalServiceMethod is a set of methods
// that are handled specially by the binding engine
// and must not be exposed to the frontend.
var internalServiceMethods = map[string]bool{
	"ServiceName":     true,
	"ServiceStartup":  true,
	"ServiceShutdown": true,
	"ServeHTTP":       true,
}

// typeError caches the type-checker type for the Go error interface.
var typeError = types.Universe.Lookup("error").Type()

// typeAny caches the empty interface type.
var typeAny = types.Universe.Lookup("any").Type().Underlying()

// collectMethod collects and returns information about a service method.
// It is intended to be called only by ServiceInfo.Collect.
func (info *ServiceInfo) collectMethod(method *types.Func) *ServiceMethodInfo {
	collector := info.collector
	obj := info.Object().(*types.TypeName)

	signature, _ := method.Type().(*types.Signature)
	if signature == nil {
		// Skip invalid interface method.
		// TODO: is this actually necessary?
		return nil
	}

	// Compute fully qualified name.
	path := obj.Pkg().Path()
	if obj.Pkg().Name() == "main" {
		// reflect.Method.PkgPath is always "main" for the main package.
		// This should not cause collisions because
		// other main packages are not importable.
		path = "main"
	}

	fqn := path + "." + obj.Name() + "." + method.Name()
	id := hash.Fnv(fqn)

	methodInfo := &ServiceMethodInfo{
		MethodInfo: collector.Method(method).Collect(),
		FQN:        fqn,
		ID:         strconv.FormatUint(uint64(id), 10),
		Params:     make([]*ParamInfo, 0, signature.Params().Len()),
		Results:    make([]types.Type, 0, signature.Results().Len()),
	}

	// Parse directives.
	if methodInfo.Doc != nil {
		var methodIdFound bool

		for _, comment := range methodInfo.Doc.List {
			switch {
			case IsDirective(comment.Text, "ignore"):
				return nil

			case IsDirective(comment.Text, "internal"):
				methodInfo.Internal = true

			case !methodIdFound && IsDirective(comment.Text, "id"):
				idString := ParseDirective(comment.Text, "id")
				idValue, err := strconv.ParseUint(idString, 10, 32)

				if err != nil {
					collector.logger.Errorf(
						"%s: invalid value '%s' in `wails:id` directive: expected a valid uint32 value",
						collector.Package(method.Pkg()).Fset.Position(comment.Pos()),
						idString,
					)
					continue
				}

				// Announce and record alias.
				collector.logger.Infof(
					"package %s: method %s.%s: default ID %s replaced by %d",
					path,
					obj.Name(),
					method.Name(),
					methodInfo.ID,
					idValue,
				)
				methodInfo.ID = strconv.FormatUint(idValue, 10)
				methodIdFound = true
			}
		}
	}

	// Collect parameters.
	for i := range signature.Params().Len() {
		param := signature.Params().At(i)

		if i == 0 {
			// Skip first parameter if it has context type.
			named, ok := types.Unalias(param.Type()).(*types.Named)
			if ok && named.Obj().Pkg().Path() == collector.systemPaths.ContextPackage && named.Obj().Name() == "Context" {
				continue
			}
		}

		if types.IsInterface(param.Type()) && !types.Identical(param.Type(), typeAny) {
			paramName := param.Name()
			if paramName == "" || paramName == "_" {
				paramName = fmt.Sprintf("#%d", i+1)
			}

			collector.logger.Warningf(
				"%s: parameter %s has non-empty interface type %s: this is not supported by encoding/json and will likely result in runtime errors",
				collector.Package(method.Pkg()).Fset.Position(param.Pos()),
				paramName,
				param.Type(),
			)
		}

		// Record type dependencies.
		info.Imports.AddType(param.Type())

		// Record parameter.
		methodInfo.Params = append(methodInfo.Params, &ParamInfo{
			Name:  param.Name(),
			Type:  param.Type(),
			Blank: param.Name() == "" || param.Name() == "_",
		})
	}

	if signature.Variadic() {
		methodInfo.Params[len(methodInfo.Params)-1].Type = methodInfo.Params[len(methodInfo.Params)-1].Type.(*types.Slice).Elem()
		methodInfo.Params[len(methodInfo.Params)-1].Variadic = true
	}

	// Collect results.
	for i := range signature.Results().Len() {
		result := signature.Results().At(i)

		if types.Identical(result.Type(), typeError) {
			// Skip error results, they are thrown as exceptions
			continue
		}

		// Record type dependencies.
		info.Imports.AddType(result.Type())

		// Record result.
		methodInfo.Results = append(methodInfo.Results, result.Type())
	}

	return methodInfo
}
