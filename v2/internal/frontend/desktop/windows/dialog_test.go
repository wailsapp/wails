//go:build windows

package windows

import (
	"testing"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"golang.org/x/sys/windows"
)

func Test_calculateMessageDialogFlags(t *testing.T) {
	tests := []struct {
		name    string
		options frontend.MessageDialogOptions
		want    uint32
	}{
		{
			name: "Test Info Dialog",
			options: frontend.MessageDialogOptions{
				Type: frontend.InfoDialog,
			},
			want: windows.MB_OK | windows.MB_ICONINFORMATION,
		},
		{
			name: "Test Error Dialog",
			options: frontend.MessageDialogOptions{
				Type: frontend.ErrorDialog,
			},
			want: windows.MB_ICONERROR | windows.MB_OK,
		},
		{
			name: "Test Question Dialog",
			options: frontend.MessageDialogOptions{
				Type: frontend.QuestionDialog,
			},
			want: windows.MB_YESNO,
		},
		{
			name: "Test Question Dialog with default cancel",
			options: frontend.MessageDialogOptions{
				Type:          frontend.QuestionDialog,
				DefaultButton: "No",
			},
			want: windows.MB_YESNO | windows.MB_DEFBUTTON2,
		},
		{
			name: "Test Question Dialog with default cancel (lowercase)",
			options: frontend.MessageDialogOptions{
				Type:          frontend.QuestionDialog,
				DefaultButton: "no",
			},
			want: windows.MB_YESNO | windows.MB_DEFBUTTON2,
		},
		{
			name: "Test Warning Dialog",
			options: frontend.MessageDialogOptions{
				Type: frontend.WarningDialog,
			},
			want: windows.MB_OK | windows.MB_ICONWARNING,
		},
		{
			name: "Test Error Dialog",
			options: frontend.MessageDialogOptions{
				Type: frontend.ErrorDialog,
			},
			want: windows.MB_ICONERROR | windows.MB_OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMessageDialogFlags(tt.options); got != tt.want {
				t.Errorf("calculateMessageDialogFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}
