//go:build windows

package edge

import (
	"syscall"
	"testing"
)

var procSetLastErrorForTest = syscall.NewLazyDLL("kernel32.dll").NewProc("SetLastError")

func setLastErrorForTest(e uint32) { procSetLastErrorForTest.Call(uintptr(e)) }

// fakeIUnknown builds AddRef/Release procs that leave a stale last-error on the
// thread (as arbitrary native code may) and count which slot was called.
type fakeIUnknown struct {
	addRefs, releases int
	vtbl              _IUnknownVtbl
}

func newFakeIUnknown(stale uint32) *fakeIUnknown {
	f := &fakeIUnknown{}
	f.vtbl.AddRef = NewComProc(func(this uintptr) uintptr {
		f.addRefs++
		setLastErrorForTest(stale)
		return uintptr(1 + f.addRefs)
	})
	f.vtbl.Release = NewComProc(func(this uintptr) uintptr {
		f.releases++
		setLastErrorForTest(stale)
		return 1
	})
	return f
}

// Release/AddRef return the reference count and are not failures, whatever
// last-error native code left behind. Each wrapper must call its own slot.
func TestReleaseIgnoresStaleLastError(t *testing.T) {
	for _, stale := range []uint32{0, 5, 1400} {
		f := newFakeIUnknown(stale)
		s := &IStream{vtbl: &_IStreamVtbl{_IUnknownVtbl: f.vtbl}}
		if got := s.AddRef(); got != 2 {
			t.Errorf("IStream.AddRef stale=%d: got %d want 2", stale, got)
		}
		if got := s.Release(); got != 1 {
			t.Errorf("IStream.Release stale=%d: got %d want 1", stale, got)
		}
		if f.addRefs != 1 || f.releases != 1 {
			t.Errorf("IStream stale=%d: addRefs=%d releases=%d", stale, f.addRefs, f.releases)
		}

		f = newFakeIUnknown(stale)
		r := &ICoreWebView2WebResourceResponse{vtbl: &_ICoreWebView2WebResourceResponseVtbl{_IUnknownVtbl: f.vtbl}}
		r.AddRef()
		if got := r.Release(); got != 1 || f.addRefs != 1 || f.releases != 1 {
			t.Errorf("Response stale=%d: got %d addRefs=%d releases=%d", stale, got, f.addRefs, f.releases)
		}

		f = newFakeIUnknown(stale)
		d := &ICoreWebView2Deferral{Vtbl: &ICoreWebView2DeferralVtbl{_IUnknownVtbl: f.vtbl}}
		d.Release()
		if f.releases != 1 || f.addRefs != 0 {
			t.Errorf("Deferral.Release stale=%d must call Release only: addRefs=%d releases=%d", stale, f.addRefs, f.releases)
		}

		f = newFakeIUnknown(stale)
		h := &ICoreWebView2HttpRequestHeaders{vtbl: &_ICoreWebView2HttpRequestHeadersVtbl{_IUnknownVtbl: f.vtbl}}
		h.AddRef()
		h.Release()
		if f.addRefs != 1 || f.releases != 1 {
			t.Errorf("HttpRequestHeaders stale=%d: addRefs=%d releases=%d", stale, f.addRefs, f.releases)
		}

		f = newFakeIUnknown(stale)
		o := &ICoreWebView2ObjectCollectionView{vtbl: &_ICoreWebView2ObjectCollectionViewVtbl{_IUnknownVtbl: f.vtbl}}
		o.AddRef()
		o.Release()
		if f.addRefs != 1 || f.releases != 1 {
			t.Errorf("ObjectCollectionView stale=%d: addRefs=%d releases=%d", stale, f.addRefs, f.releases)
		}
	}
}
