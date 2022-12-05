//go:build windows

package edge

type _ICoreWebView2NavigationCompletedEventArgsVtbl struct {
	_IUnknownVtbl
	GetIsSuccess      ComProc
	GetWebErrorStatus ComProc
	GetNavigationId   ComProc
}

type ICoreWebView2NavigationCompletedEventArgs struct {
	vtbl *_ICoreWebView2NavigationCompletedEventArgsVtbl
}

func (i *ICoreWebView2NavigationCompletedEventArgs) AddRef() uintptr {
	return i.AddRef()
}
