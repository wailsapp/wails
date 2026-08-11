package generator

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"iter"

	"github.com/wailsapp/wails/v3/internal/generator/config"
	"golang.org/x/tools/go/packages"
)

// ServiceBinding describes a concrete service type discovered by the static
// analyser and, when non-nil, the named interface that projects its frontend
// binding surface.
type ServiceBinding struct {
	Type       *types.TypeName
	Projection *types.TypeName
}

// FindServices scans the given packages for invocations
// of service registration functions from the Wails application package.
//
// Whenever one is found and the type of its unique argument
// is a valid service type, the corresponding service binding
// is fed into the returned iterator.
//
// Results are deduplicated, i.e. the iterator yields any given object at most once.
func FindServices(pkgs []*packages.Package, systemPaths *config.SystemPaths, logger config.Logger) (iter.Seq[*ServiceBinding], types.Object, error) {
	type instanceInfo struct {
		args *types.TypeList
		pos  token.Position
	}

	type target struct {
		obj   types.Object
		param int
	}

	type targetInfo struct {
		target
		cause token.Position
	}

	// instances maps objects (TypeName or Func) to their instance list.
	instances := make(map[types.Object][]instanceInfo)

	// owner maps type parameter objects to their parent object (TypeName or Func)
	owner := make(map[*types.TypeName]types.Object)

	// scheduled holds the set of type parameters
	// that have been already scheduled for analysis,
	// for deduplication.
	scheduled := make(map[target]bool)

	// registerEvent holds the `application.RegisterEvent` function if found.
	var registerEvent types.Object

	// next lists type parameter objects that have yet to be analysed.
	var next []targetInfo

	// Initialise instance/owner maps and detect application.NewService.
	for _, pkg := range pkgs {
		for ident, instance := range pkg.TypesInfo.Instances {
			obj := pkg.TypesInfo.Uses[ident]

			// Add to instance map.
			objInstances, seen := instances[obj]
			instances[obj] = append(objInstances, instanceInfo{
				instance.TypeArgs,
				pkg.Fset.Position(ident.Pos()),
			})

			if seen {
				continue
			}

			// Object seen for the first time:
			// add type params to owner map.
			var tp *types.TypeParamList

			if t, ok := obj.Type().(interface{ TypeParams() *types.TypeParamList }); ok {
				tp = t.TypeParams()
			} else {
				// Instantiated object has unexpected kind:
				// the spec might have changed.
				logger.Warningf(
					"unexpected instantiation for %s: please report this to Wails maintainers",
					types.ObjectString(obj, nil),
				)
				continue
			}

			// Add type params to owner map.
			for i := range tp.Len() {
				if param := tp.At(i).Obj(); param != nil {
					owner[param] = obj
				}
			}

			// If this is a named type, process methods.
			if recv, ok := obj.Type().(*types.Named); ok && recv.NumMethods() > 0 {
				// Register receiver type params.
				for i := range recv.NumMethods() {
					tp := recv.Method(i).Type().(*types.Signature).RecvTypeParams()
					for j := range tp.Len() {
						if param := tp.At(j).Obj(); param != nil {
							owner[param] = obj
						}
					}
				}
			}

			// Detect application.RegisterEvent
			if registerEvent == nil && obj.Name() == "RegisterEvent" && obj.Pkg().Path() == systemPaths.ApplicationPackage {
				fn, ok := obj.(*types.Func)
				if !ok {
					return nil, nil, ErrBadApplicationPackage
				}

				signature := fn.Type().(*types.Signature)
				if signature.Params().Len() != 1 || signature.Results().Len() != 0 || signature.TypeParams().Len() != 1 {
					logger.Warningf("application.RegisterService params: %d, results: %d, typeparams: %d", signature.Params().Len(), signature.Results().Len(), signature.TypeParams().Len())
					return nil, nil, ErrBadApplicationPackage
				}

				if !types.Identical(signature.Params().At(0).Type(), types.Universe.Lookup("string").Type()) {
					logger.Warningf("application.RegisterService parameter type: %v", signature.Params().At(0).Type())
					return nil, nil, ErrBadApplicationPackage
				}

				registerEvent = obj
				continue
			}

			// Detect application.NewService
			if len(next) == 0 && obj.Name() == "NewService" && obj.Pkg().Path() == systemPaths.ApplicationPackage {
				fn, ok := obj.(*types.Func)
				if !ok {
					return nil, nil, ErrBadApplicationPackage
				}

				signature := fn.Type().(*types.Signature)
				if signature.Params().Len() != 1 || signature.Results().Len() != 1 || tp.Len() != 1 {
					logger.Warningf("application.NewService params: %d, results: %d, typeparams: %d", signature.Params().Len(), signature.Results().Len(), tp.Len())
					return nil, nil, ErrBadApplicationPackage
				}

				// Schedule unique type param for analysis.
				tgt := target{obj, 0}
				scheduled[tgt] = true
				next = append(next, targetInfo{target: tgt})
				continue
			}
		}
	}

	projected, err := findProjectedServices(pkgs, systemPaths, logger)
	if err != nil {
		return nil, nil, err
	}

	// found tracks service types that have been found so far, for deduplication.
	found := make(map[*types.TypeName]*ServiceBinding)

	return func(yield func(*ServiceBinding) bool) {
		yieldBinding := func(binding *ServiceBinding) bool {
			if previous := found[binding.Type]; previous != nil {
				if previous.Projection != binding.Projection {
					logger.Errorf(
						"service type %s is registered with conflicting frontend binding projections",
						binding.Type,
					)
				}
				return true
			}

			found[binding.Type] = binding
			return yield(binding)
		}

		for _, binding := range projected {
			if !yieldBinding(binding) {
				return
			}
		}

		// Process targets.
		for len(next) > 0 {
			// Pop one target off the next list.
			tgt := next[len(next)-1]
			next = next[:len(next)-1]

			// Prepare indirect binding message.
			indirectMsg := ""
			if tgt.cause.IsValid() {
				indirectMsg = fmt.Sprintf(" (indirectly bound at %s)", tgt.cause)
			}

			for _, instance := range instances[tgt.obj] {
				// Retrieve type argument.
				serviceType := types.Unalias(instance.args.At(tgt.param))

				var named *types.Named

				switch t := serviceType.(type) {
				case *types.Named:
					// Process named type.
					named = t.Origin()

				case *types.TypeParam:
					// Schedule type parameter for analysis.
					newtgt := target{owner[t.Obj()], t.Index()}
					if !scheduled[newtgt] {
						scheduled[newtgt] = true

						// Retrieve position of call to application.NewService
						// that caused this target to be scheduled.
						cause := tgt.cause
						if !tgt.cause.IsValid() {
							// This _is_ a call to application.NewService.
							cause = instance.pos
						}

						// Push on next list.
						next = append(next, targetInfo{newtgt, cause})
					}
					continue

				default:
					logger.Warningf("%s: ignoring anonymous service type %s%s", instance.pos, serviceType, indirectMsg)
					continue
				}

				// Reject interfaces and generic types.
				if types.IsInterface(named.Underlying()) {
					logger.Warningf("%s: ignoring interface service type %s%s", instance.pos, named, indirectMsg)
					continue
				} else if named.TypeParams() != nil {
					logger.Warningf("%s: ignoring generic service type %s%s", instance.pos, named, indirectMsg)
					continue
				}

				// Record and yield type object.
				if !yieldBinding(&ServiceBinding{Type: named.Obj()}) {
					return
				}
			}
		}
	}, registerEvent, nil
}

// findProjectedServices finds direct calls to NewServiceAs and
// NewServiceAsWithOptions. Unlike NewService, the projection type argument is
// an interface, so the concrete service type must be read from the call
// argument's static type rather than from the generic instantiation alone.
func findProjectedServices(pkgs []*packages.Package, systemPaths *config.SystemPaths, logger config.Logger) ([]*ServiceBinding, error) {
	var result []*ServiceBinding
	badApplicationPackage := false

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(node ast.Node) bool {
				if badApplicationPackage {
					return false
				}
				call, ok := node.(*ast.CallExpr)
				if !ok || len(call.Args) == 0 {
					return true
				}

				ident := calledFunctionIdent(call.Fun)
				if ident == nil {
					return true
				}
				obj := pkg.TypesInfo.Uses[ident]
				if obj == nil || obj.Pkg() == nil || obj.Pkg().Path() != systemPaths.ApplicationPackage {
					return true
				}
				if obj.Name() != "NewServiceAs" && obj.Name() != "NewServiceAsWithOptions" {
					return true
				}

				fn, ok := obj.(*types.Func)
				if !ok {
					return true
				}
				signature := fn.Type().(*types.Signature)
				typeParams := signature.TypeParams()
				wantParams := 1
				if obj.Name() == "NewServiceAsWithOptions" {
					wantParams = 2
				}
				if signature.Params().Len() != wantParams || signature.Results().Len() != 1 || typeParams == nil || typeParams.Len() != 1 {
					badApplicationPackage = true
					return false
				}

				instance, ok := pkg.TypesInfo.Instances[ident]
				if !ok || instance.TypeArgs.Len() != 1 {
					return true
				}

				position := pkg.Fset.Position(ident.Pos())
				projection, ok := namedProjection(instance.TypeArgs.At(0))
				if !ok {
					// The implementation of NewServiceAsWithOptions delegates to
					// NewServiceAs with an unresolved type parameter. Its concrete
					// call sites are analysed separately.
					if pkg.PkgPath != systemPaths.ApplicationPackage {
						logger.Warningf("%s: ignoring service binding projection %s: expected a named interface", position, instance.TypeArgs.At(0))
					}
					return true
				}

				interfaceType := projection.Underlying().(*types.Interface).Complete()
				for i := range interfaceType.NumMethods() {
					if method := interfaceType.Method(i); !method.Exported() {
						logger.Errorf("%s: service binding projection %s contains unexported method %s", position, projection, method.Name())
						return true
					}
				}

				argumentType := pkg.TypesInfo.TypeOf(call.Args[0])
				service, ok := namedServicePointer(argumentType)
				if !ok {
					if pkg.PkgPath != systemPaths.ApplicationPackage {
						logger.Warningf(
							"%s: cannot determine concrete service type from projected service argument %s; pass a concrete service pointer directly",
							position,
							argumentType,
						)
					}
					return true
				}

				if !types.AssignableTo(argumentType, projection) {
					logger.Errorf("%s: service type %s does not implement binding projection %s", position, argumentType, projection)
					return true
				}

				result = append(result, &ServiceBinding{Type: service.Obj(), Projection: projection.Obj()})
				return true
			})
		}
	}
	if badApplicationPackage {
		return nil, ErrBadApplicationPackage
	}

	return result, nil
}

func calledFunctionIdent(expr ast.Expr) *ast.Ident {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr
	case *ast.SelectorExpr:
		return expr.Sel
	case *ast.IndexExpr:
		return calledFunctionIdent(expr.X)
	case *ast.IndexListExpr:
		return calledFunctionIdent(expr.X)
	case *ast.ParenExpr:
		return calledFunctionIdent(expr.X)
	default:
		return nil
	}
}

func namedProjection(t types.Type) (*types.Named, bool) {
	if t == nil {
		return nil, false
	}
	named, ok := types.Unalias(t).(*types.Named)
	if !ok {
		return nil, false
	}
	named = named.Origin()
	if _, ok := named.Underlying().(*types.Interface); !ok {
		return nil, false
	}
	if params := named.TypeParams(); params != nil && params.Len() > 0 {
		return nil, false
	}
	return named, true
}

func namedServicePointer(t types.Type) (*types.Named, bool) {
	if t == nil {
		return nil, false
	}
	pointer, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return nil, false
	}
	named, ok := types.Unalias(pointer.Elem()).(*types.Named)
	if !ok {
		return nil, false
	}
	named = named.Origin()
	if params := named.TypeParams(); params != nil && params.Len() > 0 {
		return nil, false
	}
	return named, true
}
