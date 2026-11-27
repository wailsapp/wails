//go:build windows

package w32

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestWindowLongPtrProceduresResolve(testContext *testing.T) {
	for _, procedure := range []*windows.LazyProc{User32GetWindowLongPtrW, User32SetWindowLongPtrW} {
		if err := procedure.Find(); err != nil {
			testContext.Errorf("%s: %v", procedure.Name, err)
		}
	}
}
