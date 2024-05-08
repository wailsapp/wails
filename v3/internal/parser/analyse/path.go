package analyse

import (
	"go/types"
	"slices"

	"github.com/samber/lo"
	"golang.org/x/tools/go/types/typeutil"
)

// Path represents a sequence of access operations
// that when applied to a certain target variable or field
// should yield a binding expression.
//
// The analyser has two principal modes of operation:
//
//   - [Analyser.processReference] starts from
//     an identifier referencing a target variable
//     and walks up its context in the syntax tree;
//     if it finds operations that match a given Path,
//     it pops them off the path and keeps going,
//     otherwise it stops.
//   - [Analyser.processExpression] starts from
//     a complex expression assigned to a target variable or field
//     and breaks it down into simpler parts,
//     recording the operations it finds into a Path
//     until it reaches an identifier referencing a variable.
//
// Both methods access frequently the first or second element
// of the Path they operate upon, hence it is convenient
// to assign them indexes 0 and 1;
// however, both methods frequently prepend elements to the path,
// which may result in inefficient copying and allocation patterns.
//
// To mitigate this problem, elements in a Path are kept
// in reverse order, and many methods are provided
// that translate indices transparently.
type Path []Step

// Step represents a single access operation in a Path.
// Some operations use the adjacent element as an argument.
//
// All special operations defined below are negative values.
// Positive values represent struct field selections.
//
// The underlying type must be large enough to fit an int and a uint32.
// This is due to the fact that step codes must be able to represent
// struct field indices, which have type int, and Go types as encoded
// by a [typeutil.Hasher], which outputs uint32 values.
//
// At the time of writing, the only simple way to satisfy this constraint
// is to choose int64.
type Step int64

const (
	// IndexingStep represents a slice or map indexing operation.
	IndexingStep Step = -1 - iota

	// WeakIndexingStep represents a slice or map indexing operation
	// where we track assignments to the indexed element
	// but not to the map or slice itself.
	WeakIndexingStep

	// ArrayIndexingStep represents an array indexing operation.
	// The distinction is required because arrays do not have reference
	// semantics, while slices do.
	// Moreover, array indexing operations may include
	// an implicit pointer indirection,
	// while slice and map indexing do not operate on pointers.
	ArrayIndexingStep

	// IndirectionStep represents a pointer indirection operation.
	IndirectionStep

	// WeakIndirectionStep represents a pointer indirection operation
	// where we track assignments to the pointee but not to the pointer itself.
	WeakIndirectionStep

	// ChanSendStep signifies that the current target is a channel
	// and values sent into it may be used as bindings,
	// hence we track those values.
	ChanSendStep

	// WeakChanSendStep is equivalent to ChanSendStep,
	// but we do not track assignments to the channel variable.
	WeakChanSendStep

	// ChanRecvStep signifies that the current target is a channel
	// and values sent into it may refer a binding source,
	// hence we track how receives from that channel are used.
	ChanRecvStep

	// WeakChanRecvStep is equivalent to ChanRecvStep,
	// but we do not track assignments to the channel variable.
	WeakChanRecvStep

	// TypeAssertionStep represents a type assertion operation.
	// The adjacent step represents the asserted type
	// as encoded by [PathSet.TypeStep].
	TypeAssertionStep

	// MethodLookupStep represents a method lookup operation.
	// It is used when resolving abstract interface methods to concrete methods.
	// The adjacent step represents the package and name of the method
	// as encoded by [PathSet.MethodStep].
	MethodLookupStep

	// ParamStep signifies that the current target is a function expression
	// that we are looking to resolve to a concrete function.
	// The adjacent step represents the index of a function parameter
	// that should be scheduled as an additional target
	// whenever resolution succeeds.
	ParamStep

	// ResultStep is like ParamStep, but targets function result fields.
	ResultStep

	// InvalidStep is returned by [Path.At]
	// when the given index is out of range.
	// All predicates defined below are false for an InvalidStep.
	InvalidStep
)

// IsSelection returns true if and only if
// the given step code represents a subfield selection.
func (step Step) IsSelection() bool {
	return step >= 0
}

// IsIndexing returns true if and only if
// the given step code is either IndexingStep,
// WeakIndexingStep or ArrayIndexingStep.
func (step Step) IsIndexing() bool {
	return step == IndexingStep || step == WeakIndexingStep || step == ArrayIndexingStep
}

// IsIndirection returns true if and only if
// the given step code is either IndirectionStep or WeakIndirectionStep.
func (step Step) IsIndirection() bool {
	return step == IndirectionStep || step == WeakIndirectionStep
}

// IsChanSend returns true if and only if
// the given step code is either ChanSendStep or WeakChanSendStep.
func (step Step) IsChanSend() bool {
	return step == ChanSendStep || step == WeakChanSendStep
}

// IsChanRecv returns true if and only if
// the given step code is either ChanRecvStep or WeakChanRecvStep.
func (step Step) IsChanRecv() bool {
	return step == ChanRecvStep || step == WeakChanRecvStep
}

// IsFunc returns true if and only if
// the given step code is either ParamStep or ResultStep.
func (step Step) IsFunc() bool {
	return step == ParamStep || step == ResultStep
}

// IsRef returns true if and only if
// the given step code applies to values with reference semantics,
// i.e. copies of that value may be used to modify the original target.
func (step Step) IsRef() bool {
	return step.IsWeakRef() || step.IsStrongRef()

}

// IsStrongRef returns true if and only if
// the given step code applies to values with reference semantics
// and the analyser should track changes to the reference itself.
func (step Step) IsStrongRef() bool {
	return step == IndexingStep ||
		step == IndirectionStep ||
		step == ChanSendStep ||
		step == ChanRecvStep
}

// IsWeakRef returns true if and only if
// the given step code applies to values with reference semantics
// and the analyser should _not_ track changes to the reference itself.
func (step Step) IsWeakRef() bool {
	return step == WeakIndexingStep ||
		step == WeakIndirectionStep ||
		step == WeakChanSendStep ||
		step == WeakChanRecvStep
}

// NewPath is a convenience function that forms a path from the given steps
// by reversing their order. Note that the input slice is reversed in place,
// then returned as is.
func NewPath(steps ...Step) Path {
	slices.Reverse(steps)
	return steps
}

// LeftmostVisibleRef returns the step of the leftmost Ref step
// in the given path that is not preceded by a Func step.
// If no such Ref step exists, LeftmostVisibleRef returns -1.
func (path Path) LeftmostVisibleRef() int {
	ref := -1
	for i := len(path) - 1; i >= 0; i-- {
		if path[i].IsFunc() {
			break
		}

		if path[i].IsRef() {
			ref = i
			break
		}
	}
	return ref
}

// HasRef returns true if and only if
// the given path contains a visible HasRef step
// as defined by [Path.LeftmostVisibleRef].
func (path Path) HasRef() bool {
	return path.LeftmostVisibleRef() >= 0
}

// HasStrongRef returns true if and only if
// the given path contains a visible Ref step
// and the leftmost such is strong.
func (path Path) HasStrongRef() bool {
	ref := path.LeftmostVisibleRef()
	return ref >= 0 && path[ref].IsStrongRef()
}

// HasWeakRef returns true if and only if
// the given path contains a visible Ref step
// and the leftmost such is strong.
func (path Path) HasWeakRef() bool {
	ref := path.LeftmostVisibleRef()
	return ref >= 0 && path[ref].IsWeakRef()
}

// Clone returns a shallow copy of the given path.
func (path Path) Clone() Path {
	return slices.Clone(path)
}

// CloneAndGrow returns a shallow copy of the given path
// whose capacity is sufficient for prepending n more elements
// without further allocations.
func (path Path) CloneAndGrow(n int) Path {
	if n < 0 {
		panic("cannot be negative")
	}

	result := make(Path, len(path), len(path)+n)
	copy(result, path)
	return result
}

// Weaken replaces the leftmost visible
// Ref step in the given path, if any,
// with its weak version,
// then returns the receiver.
func (path Path) Weaken() Path {
	if ref := path.LeftmostVisibleRef(); ref >= 0 {
		switch path[ref] {
		case IndexingStep:
			path[ref] = WeakIndexingStep
		case IndirectionStep:
			path[ref] = WeakIndirectionStep
		case ChanSendStep:
			path[ref] = WeakChanSendStep
		case ChanRecvStep:
			path[ref] = WeakChanRecvStep
		}
	}

	return path
}

// Prepend prepends the given prefix to the receiver
// and returns the resulting path.
// If the receiver has enough capacity,
// no memory is allocated.
//
// After Prepend returns, the input slice is reversed.
func (path Path) Prepend(prefix ...Step) Path {
	slices.Reverse(prefix)
	return path.PrependPath(prefix)
}

// PrependPath prepends the given prefix to the receiver
// and returns the resulting path.
// If the receiver has enough capacity,
// no memory is allocated.
func (path Path) PrependPath(prefix Path) Path {
	return append(path, prefix...)
}

// Consume returns the subpath of the given path
// obtained by removing the first n steps.
func (path Path) Consume(n int) Path {
	return path[:len(path)-n]
}

// At returns the i-th step in the given path.
// If no such step exists, At returns InvalidStep.
func (path Path) At(i int) Step {
	if i >= 0 && i < len(path) {
		return path[len(path)-i-1]
	} else {
		return InvalidStep
	}
}

// Set assigns the given step to the i-th element in the given path.
func (path Path) Set(i int, step Step) Path {
	path[len(path)-i-1] = step
	return path
}

// ConcatPaths is equivalent to
//
//	paths[len(paths)-1].Clone().
//		PrependPath(paths[len(paths)-2]).
//		PrependPath(paths[len(paths)-3])...
//
// i.e. it returns a new path obtained by concatenating the given paths,
// but allocates memory just once.
func ConcatPaths(paths ...Path) Path {
	result := make(Path, lo.SumBy(paths, func(path Path) int { return len(path) }))
	i := 0
	for j := len(paths); j >= 0; j-- {
		copy(result[i:], paths[j])
		i += len(paths[j])
	}
	return result
}

// PathSet is a radix tree that represents
// a set of [Path] values efficiently.
//
// It is used by the static analyser to canonicalise
// access paths and make them usable as map keys.
type PathSet struct {
	root      *pathNode
	hasher    typeutil.Hasher
	methods   []methodId
	methodMap map[methodId]int
}

// PathSetElement is an opaque type that represents an entry in a [PathSet].
type PathSetElement struct {
	*pathNode
}

// pathNode holds data for a node in the [PathSet] tree.
type pathNode struct {
	// path holds the access path represented by this node.
	path Path

	// extensions holds pointers to pathNodes
	// that represent proper extensions of this node.
	// Keys store the first element of the extension suffix.
	extensions map[Step]*pathNode
}

// methodId holds parameters for method lookup.
type methodId struct {
	pkg  *types.Package
	name string
}

// SetHasher sets the type hasher used by PathSet.
//
// The hasher is used by [PathSet.TypeAssertionStep]
// to compute path elements representing type assertions.
//
// A single Hasher created by may be shared among many PathSets.
// See [typeutil.Map.SetHasher] for more information.
func (ps *PathSet) SetHasher(hasher typeutil.Hasher) {
	if ps.root == nil {
		ps.root = &pathNode{
			extensions: make(map[Step]*pathNode),
		}
	}

	ps.hasher = hasher
}

// TypeStep returns the step code that,
// when occurring immediately after a TypeAssertionStep,
// represents an assertion to the given type.
func (ps *PathSet) TypeStep(typ types.Type) Step {
	if ps.root == nil {
		ps.root = &pathNode{
			extensions: make(map[Step]*pathNode),
		}
	}

	return Step(ps.hasher.Hash(typ))
}

// MethodStep returns the step code that,
// when occurring immediately after a MethodLookupStep,
// represents a method lookup with the given id.
func (ps *PathSet) MethodStep(pkg *types.Package, name string) Step {
	id := methodId{pkg, name}
	index, ok := ps.methodMap[id]
	if !ok {
		if ps.methodMap == nil {
			ps.methodMap = make(map[methodId]int)
		}

		index = len(ps.methods)
		ps.methods = append(ps.methods, id)
		ps.methodMap[id] = index
	}

	return Step(index)
}

// MethodStepParams returns the method package and name, if any,
// associated to the given method selection step.
func (ps *PathSet) MethodStepParams(step Step) (*types.Package, string) {
	index := int(step)
	if index >= 0 && index < len(ps.methods) {
		id := ps.methods[index]
		return id.pkg, id.name
	} else {
		return nil, ""
	}
}

// Get returns the unique [PathSetElement]
// associated to the given path in a PathSet.
//
// The following properties are guaranteed to hold:
//   - slices.Equal(path, ps.Get(path).Path());
//   - if slices.Equal(path1, path2),
//     then ps.Get(path1) == ps.Get(path2).
func (ps *PathSet) Get(path Path) PathSetElement {
	if ps.root == nil {
		ps.SetHasher(typeutil.MakeHasher())
		ps.hasher = typeutil.MakeHasher()
	}

	return PathSetElement{ps.root.prepend(path)}
}

// Path returns a shallow copy of the Path
// that the given PathSetElement represents canonically.
func (p PathSetElement) Path() Path {
	return p.path.Clone()
}

// StrongRef is equivalent to p.Path().StrongRef(),
// but does not make a copy of the represented path.
func (p PathSetElement) StrongRef() bool {
	return p.path.HasStrongRef()
}

// prepend returns the pathNode that represents the path obtained
// by prepending the given prefix to the receiver
// in the same [PathSet] n belongs to.
// I.e. if n belongs to set ps, then
//
//	n.prepend(prefix) == ps.Get(n.path.Clone().PrependPath(prefix)).
func (n *pathNode) prepend(prefix Path) *pathNode {
	for len(prefix) > 0 {
		head := prefix[0]
		ext := n.extensions[head]

		if ext == nil {
			// No suitable extension node is present, insert and return a new one.
			result := &pathNode{
				slices.Concat(n.path, prefix),
				make(map[Step]*pathNode),
			}

			n.extensions[head] = result
			return result
		}

		// Extension node found: determine whether to advance or split.
		lcp := LongestCommonPrefix(ext.path[len(n.path):], prefix)
		if lcp < 1 {
			// This should not happen...
			panic("longest common prefix is empty")
		}

		if splitDepth := len(n.path) + lcp; splitDepth == len(ext.path) {
			// Extension node is a prefix of the result path:
			// pass it on to the next iteration.
			prefix = prefix[lcp:]
			n = ext
			continue
		} else if lcp == len(prefix) {
			// Result path is a proper prefix of
			// the extension node: split and replace it.
			result := &pathNode{
				ext.path[:splitDepth],
				map[Step]*pathNode{
					ext.path[splitDepth]: ext,
				},
			}

			n.extensions[head] = result
			return result
		} else {
			// Extension node and result path diverge somewhere:
			// split and replace the extension node, adding a new leaf.

			result := &pathNode{
				slices.Concat(n.path, prefix),
				make(map[Step]*pathNode),
			}

			n.extensions[head] = &pathNode{
				ext.path[:splitDepth],
				map[Step]*pathNode{
					ext.path[splitDepth]: ext,
					prefix[lcp]:          result,
				},
			}

			return result
		}
	}

	return n
}
