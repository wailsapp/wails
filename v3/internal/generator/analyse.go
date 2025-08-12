package generator

import (
	"fmt"
	"go/token"
	"go/types"
	"iter"

	"github.com/wailsapp/wails/v3/internal/generator/config"
	"golang.org/x/tools/go/packages"
)

// FindServices scans the given packages for invocations
// of the NewService function from the Wails application package.
//
// Whenever one is found and the type of its unique argument
// is a valid service type, the corresponding named type object
// is passed to yield.
//
// Results are deduplicated, i.e. yield is called at most once per object.
//
// If yield returns false, FindBoundTypes returns immediately.
func FindServices(pkgs []*packages.Package, systemPaths *config.SystemPaths, logger config.Logger) (iter.Seq[*types.TypeName], error) {
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

			if len(next) > 0 {
				// application.NewService has been found already.
				continue
			}

			fn, ok := obj.(*types.Func)
			if !ok {
				continue
			}

			// Detect application.NewService
			if fn.Name() == "NewService" && fn.Pkg().Path() == systemPaths.ApplicationPackage {
				// Check signature.
				signature := fn.Type().(*types.Signature)
				if signature.Params().Len() > 2 || signature.Results().Len() != 1 || tp.Len() != 1 || tp.At(0).Obj() == nil {
					logger.Warningf("Param Len: %d, Results Len: %d, tp.Len: %d, tp.At(0).Obj(): %v", signature.Params().Len(), signature.Results().Len(), tp.Len(), tp.At(0).Obj())
					return nil, ErrBadApplicationPackage
				}

				// Schedule unique type param for analysis.
				tgt := target{obj, 0}
				scheduled[tgt] = true
				next = append(next, targetInfo{target: tgt})
			}
		}
	}

	// found tracks service types that have been found so far, for deduplication.
	found := make(map[*types.TypeName]bool)

	return func(yield func(*types.TypeName) bool) {
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
					logger.Warningf("%s: ignoring generic service type %s", instance.pos, named, indirectMsg)
					continue
				}

				// Record and yield type object.
				if !found[named.Obj()] {
					found[named.Obj()] = true
					if !yield(named.Obj()) {
						return
					}
				}
			}
		}
	}, nil
}
