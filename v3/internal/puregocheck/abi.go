package main

import (
	"go/types"
)

// Stack-argument placement on darwin/arm64 (the stacknarrow rule).
//
// The AAPCS64 register allocation is the same on every arm64 platform: the
// first eight integer-class arguments go in x0..x7, the first eight floating
// point (and HFA) arguments in v0..v7, and whatever does not fit goes on the
// stack. The platforms differ in the stack layout:
//
//   - AAPCS64 (Linux, Android): every stacked argument takes its own 8-byte
//     slot, so a bool, uint16 or float32 on the stack occupies a full slot.
//   - Apple's arm64 ABI ("Writing ARM64 code for Apple platforms"): stacked
//     arguments are packed at their natural size and alignment, so a bool
//     followed by a uint16 share one slot. Once one narrow argument is packed,
//     every later stacked argument also moves, so the first narrow spill is
//     the one to report.
//
// purego v0.11.1 knows about the Apple layout. For a fixed-signature call on
// darwin it bundles the stacked arguments into a struct with C alignment
// (struct_arm64.go: shouldBundleStackArgs, bundleStackArgs), so those calls
// are placed correctly. Its variadic path does not: every Objective-C send
// goes through objc_msgSend(obj ID, cmd SEL, args ...any), and the variadic
// expansion in func.go calls addValue per element, so each stacked element
// gets its own 8-byte slot. That is the gap this rule reports. Linux is not
// affected (AAPCS64 slots), so the rule only runs for darwin.

// argKind is the register class of one argument.
type argKind int

const (
	kindZero  argKind = iota // zero-sized: takes no register or stack space
	kindInt                  // one integer register (scalars, pointers, large structs by reference)
	kindFloat                // one floating point register
	kindHFA                  // homogeneous float aggregate: n floating point registers
	kindSmall                // non-HFA composite of at most 16 bytes: regs integer registers
)

// cArg is one argument as the C callee sees it.
type cArg struct {
	kind   argKind
	n      int   // kindHFA: number of members; kindSmall: integer registers used
	size   int64 // size in bytes of the Go value
	narrow bool  // stacked, it is placed at the wrong offset by purego on darwin/arm64
}

// allocate assigns each argument to registers or the stack in AAPCS64 order.
// preInt is the number of integer registers already used by the caller (2 for
// an Objective-C message: self and _cmd). It returns the index of every
// argument that spills to the stack.
//
// The spill rules follow AAPCS64 C.1-C.17: a scalar that does not fit leaves
// the register counters alone, while a composite that does not fit (an HFA or
// a small aggregate) exhausts its register file, so later arguments of that
// class also go to the stack.
func allocate(args []cArg, preInt int) []int {
	ngrn, nsrn := preInt, 0
	var spilled []int
	for i, a := range args {
		switch a.kind {
		case kindZero:
		case kindInt:
			if ngrn < 8 {
				ngrn++
			} else {
				spilled = append(spilled, i)
			}
		case kindFloat:
			if nsrn < 8 {
				nsrn++
			} else {
				spilled = append(spilled, i)
			}
		case kindHFA:
			if nsrn+a.n <= 8 {
				nsrn += a.n
			} else {
				nsrn = 8
				spilled = append(spilled, i)
			}
		case kindSmall:
			if ngrn+a.n <= 8 {
				ngrn += a.n
			} else {
				ngrn = 8
				spilled = append(spilled, i)
			}
		}
	}
	return spilled
}

// leaf describes the scalar members of a composite, in order.
type leaf struct {
	float bool
	size  int64
}

// leaves flattens a type into its scalar members. ok is false for types the
// C mapping does not cover (interfaces, complex numbers, strings inside
// composites, ...).
func leaves(t types.Type, sizes types.Sizes) (out []leaf, ok bool) {
	switch u := t.Underlying().(type) {
	case *types.Basic:
		switch u.Kind() {
		case types.Bool, types.Int8, types.Int16, types.Int32, types.Int64, types.Int,
			types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uint, types.Uintptr,
			types.UnsafePointer:
			return []leaf{{false, sizes.Sizeof(u)}}, true
		case types.Float32, types.Float64:
			return []leaf{{true, sizes.Sizeof(u)}}, true
		}
		return nil, false
	case *types.Pointer, *types.Map, *types.Chan, *types.Signature, *types.Slice:
		return []leaf{{false, 8}}, true
	case *types.Array:
		elem, ok := leaves(u.Elem(), sizes)
		if !ok {
			return nil, false
		}
		for i := int64(0); i < u.Len(); i++ {
			out = append(out, elem...)
		}
		return out, true
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			f, ok := leaves(u.Field(i).Type(), sizes)
			if !ok {
				return nil, false
			}
			out = append(out, f...)
		}
		return out, true
	}
	return nil, false
}

// classify maps a Go type to its C argument class. ok is false when the type
// cannot be mapped, in which case the caller must not guess its placement.
// String and slice values are passed as their data pointer, which is what the
// purego backend hands the C side.
func classify(t types.Type, sizes types.Sizes) (cArg, bool) {
	if b, ok := t.Underlying().(*types.Basic); ok {
		switch {
		case b.Kind() == types.String || b.Kind() == types.UntypedString:
			return cArg{kind: kindInt, size: 8}, true
		case b.Kind() == types.UntypedNil:
			return cArg{kind: kindInt, size: 8}, true
		case b.Kind() == types.UntypedBool:
			return cArg{kind: kindInt, size: 1, narrow: true}, true
		case b.Kind() == types.UntypedInt || b.Kind() == types.UntypedRune:
			return cArg{kind: kindInt, size: 8}, true
		case b.Kind() == types.UntypedFloat:
			return cArg{kind: kindFloat, size: 8}, true
		}
	}
	switch u := t.Underlying().(type) {
	case *types.Basic:
		switch u.Kind() {
		case types.Float32:
			return cArg{kind: kindFloat, size: 4, narrow: true}, true
		case types.Float64:
			return cArg{kind: kindFloat, size: 8}, true
		}
		if _, ok := leaves(u, sizes); ok {
			size := sizes.Sizeof(u)
			return cArg{kind: kindInt, size: size, narrow: size < 8}, true
		}
		return cArg{}, false
	case *types.Pointer, *types.Map, *types.Chan, *types.Signature, *types.Slice:
		return cArg{kind: kindInt, size: 8}, true
	case *types.Struct, *types.Array:
		ls, ok := leaves(u, sizes)
		if !ok {
			return cArg{}, false
		}
		size := sizes.Sizeof(u)
		if size == 0 {
			return cArg{kind: kindZero}, true
		}
		narrow := size%8 != 0
		for _, l := range ls {
			if l.size < 8 {
				narrow = true
			}
		}
		if hfa, n := isHFA(ls); hfa {
			return cArg{kind: kindHFA, n: n, size: size, narrow: narrow}, true
		}
		if size <= 16 {
			return cArg{kind: kindSmall, n: int((size + 7) / 8), size: size, narrow: narrow}, true
		}
		// Larger composites are passed by reference: a pointer in an integer
		// register, so they are never narrow.
		return cArg{kind: kindInt, size: 8}, true
	}
	return cArg{}, false
}

// isHFA reports whether the leaves are one floating point type with 1 to 4
// members (AAPCS64 homogeneous float aggregate), and how many there are.
func isHFA(ls []leaf) (bool, int) {
	if len(ls) == 0 || len(ls) > 4 || !ls[0].float {
		return false, 0
	}
	for _, l := range ls {
		if !l.float || l.size != ls[0].size {
			return false, 0
		}
	}
	return true, len(ls)
}
