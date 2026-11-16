package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/ebitengine/purego"
)

var (
	objcGetClass            func(name string) uintptr
	objcGetClassList        func(buf *uintptr, n int32) int32
	objectGetClass          func(obj uintptr) uintptr
	classRespondsToSelector func(cls, sel uintptr) bool
	selRegisterName         func(name string) uintptr
)

// resolve checks C symbols, then that every class the backend names exists
// and every selector it sends is implemented by at least one class in the
// runtime (preferring the classes the backend names). A selector nothing
// implements is a guaranteed doesNotRecognizeSelector: crash, which pure Go
// cannot catch.
func resolve(r *report) error {
	loadLibraries(r)
	if data, err := os.ReadFile("/System/Library/CoreServices/SystemVersion.plist"); err == nil {
		if m := regexp.MustCompile(`ProductVersion</key>\s*<string>([^<]+)`).FindSubmatch(data); m != nil {
			r.Loaded = append(r.Loaded, "macOS "+string(m[1]))
		}
	}
	objc, err := purego.Dlopen("/usr/lib/libobjc.A.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("libobjc: %w", err)
	}
	purego.RegisterLibFunc(&objcGetClass, objc, "objc_getClass")
	purego.RegisterLibFunc(&objcGetClassList, objc, "objc_getClassList")
	purego.RegisterLibFunc(&objectGetClass, objc, "object_getClass")
	purego.RegisterLibFunc(&classRespondsToSelector, objc, "class_respondsToSelector")
	purego.RegisterLibFunc(&selRegisterName, objc, "sel_registerName")
	resolveSymbols(r)

	var named []uintptr
	for _, e := range r.Classes {
		if e.Ours {
			continue
		}
		cls := objcGetClass(e.Name)
		if cls == 0 {
			if e.Optional {
				r.Absent = append(r.Absent, "class "+e.Name)
			} else {
				r.Missing = append(r.Missing, "class "+e.Name+" ("+e.Where[0]+")")
			}
			continue
		}
		named = append(named, cls, objectGetClass(cls))
	}

	var all []uintptr
	allClasses := func() []uintptr {
		if all == nil {
			n := objcGetClassList(nil, 0)
			buf := make([]uintptr, n)
			n = objcGetClassList(&buf[0], n)
			for _, c := range buf[:n] {
				all = append(all, c, objectGetClass(c))
			}
		}
		return all
	}
	implemented := func(sel uintptr, classes []uintptr) bool {
		for _, c := range classes {
			if c != 0 && classRespondsToSelector(c, sel) {
				return true
			}
		}
		return false
	}
	for _, e := range r.Selectors {
		if e.Ours {
			continue
		}
		sel := selRegisterName(e.Name)
		if implemented(sel, named) || implemented(sel, allClasses()) {
			continue
		}
		if e.Optional {
			r.Absent = append(r.Absent, "selector "+e.Name)
		} else {
			r.Missing = append(r.Missing, "selector "+e.Name+" ("+e.Where[0]+")")
		}
	}
	return nil
}
